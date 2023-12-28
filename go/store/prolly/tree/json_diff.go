// Copyright 2023 Dolthub, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tree

import (
	"bytes"
	"fmt"
	"github.com/dolthub/go-mysql-server/sql/types"
	"io"
	"reflect"
	"strings"
)

type JsonDiff struct {
	Key      string
	From, To *types.JSONDocument
	Type     DiffType
}

type jsonKeyPair struct {
	key   string
	value interface{}
}

// Differ computes the diff between two json objects.
// Once we migrate JSON storage to a format that can be read without being fully deserialized, we should be able to
// use Differ instead.
type JsonDiffer struct {
	root                           string
	currentFromPair, currentToPair *jsonKeyPair
	from, to                       types.JSONIter
	subDiffer                      *JsonDiffer
}

func NewJsonDiffer(root string, from, to types.JsonObject) JsonDiffer {
	fromIter := types.NewJSONIter(from)
	toIter := types.NewJSONIter(to)
	return JsonDiffer{
		root: root,
		from: fromIter,
		to:   toIter,
	}
}

func (differ *JsonDiffer) appendKey(key string) string {
	escapedKey := strings.Replace(key, "\"", "\\\"", -1)
	return fmt.Sprintf("%s.\"%s\"", differ.root, escapedKey)
}

func (differ *JsonDiffer) Next() (diff JsonDiff, err error) {
	for {
		if differ.subDiffer != nil {
			diff, err := differ.subDiffer.Next()
			if err == io.EOF {
				differ.subDiffer = nil
				differ.currentFromPair = nil
				differ.currentToPair = nil
				continue
			} else if err != nil {
				return JsonDiff{}, err
			}
			return diff, nil
		}
		if differ.currentFromPair == nil && differ.from.HasNext() {
			key, value, err := differ.from.Next()
			if err != nil {
				return JsonDiff{}, err
			}
			differ.currentFromPair = &jsonKeyPair{key, value}
		}

		if differ.currentToPair == nil && differ.to.HasNext() {
			key, value, err := differ.to.Next()
			if err != nil {
				return JsonDiff{}, err
			}
			differ.currentToPair = &jsonKeyPair{key, value}
		}

		if differ.currentFromPair == nil && differ.currentToPair == nil {
			return JsonDiff{}, io.EOF
		}

		var diffType DiffType

		if differ.currentFromPair == nil && differ.currentToPair != nil {
			diffType = AddedDiff
		} else if differ.currentFromPair != nil && differ.currentToPair == nil {
			diffType = RemovedDiff
		} else { // differ.currentFromPair != nil && differ.currentToPair != nil
			keyCmp := bytes.Compare([]byte(differ.currentFromPair.key), []byte(differ.currentToPair.key))
			if keyCmp > 0 {
				// `to` key comes before `from` key. Right key must have been inserted.
				diffType = AddedDiff
			} else if keyCmp < 0 {
				// `to` key comes after `from` key. Right key must have been deleted.
				diffType = RemovedDiff
			} else {
				key := differ.currentFromPair.key
				fromValue := differ.currentFromPair.value
				toValue := differ.currentToPair.value
				if reflect.TypeOf(fromValue) != reflect.TypeOf(toValue) {
					diffType = ModifiedDiff
				} else {
					switch from := fromValue.(type) {
					case types.JsonObject:
						// Recursively compare the objects to generate diffs.
						subDiffer := NewJsonDiffer(differ.appendKey(key), from, toValue.(types.JsonObject))
						differ.subDiffer = &subDiffer
						continue
					case types.JsonArray:
						if reflect.DeepEqual(fromValue, toValue) {
							differ.currentFromPair = nil
							differ.currentToPair = nil
							continue
						} else {
							diffType = ModifiedDiff
						}
					default:
						if fromValue == toValue {
							differ.currentFromPair = nil
							differ.currentToPair = nil
							continue
						}
						diffType = ModifiedDiff
					}
				}
			}
		}

		switch diffType {
		case AddedDiff:
			key, value := differ.currentToPair.key, differ.currentToPair.value
			result := JsonDiff{
				Key:  differ.appendKey(key),
				From: nil,
				To:   &types.JSONDocument{Val: value},
				Type: AddedDiff,
			}
			differ.currentToPair = nil
			return result, nil
		case RemovedDiff:
			key, value := differ.currentFromPair.key, differ.currentFromPair.value
			result := JsonDiff{
				Key:  differ.appendKey(key),
				From: &types.JSONDocument{Val: value},
				To:   nil,
				Type: RemovedDiff,
			}
			differ.currentFromPair = nil
			return result, nil
		case ModifiedDiff:
			key := differ.currentFromPair.key
			from := differ.currentFromPair.value
			to := differ.currentToPair.value
			result := JsonDiff{
				Key:  differ.appendKey(key),
				From: &types.JSONDocument{Val: from},
				To:   &types.JSONDocument{Val: to},
				Type: ModifiedDiff,
			}
			differ.currentFromPair = nil
			differ.currentToPair = nil
			return result, nil
		default:
			panic("unreachable")
		}
	}
}
