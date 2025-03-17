// Copyright 2025 Dolthub, Inc.
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

package val

import (
	"context"
	"github.com/dolthub/dolt/go/store/hash"
	"github.com/dolthub/go-mysql-server/sql"
	"github.com/dolthub/go-mysql-server/sql/types"
)

const BytePeekLength = 128

// ValueStore is an interface for a key-value store that can store byte sequences, keyed by a content hash.
// The only implementation is tree.NodeStore, but ValueStore can be used without depending on the tree package.
// This is useful for type handlers.
type ValueStore interface {
	ReadBytes(ctx context.Context, h hash.Hash) ([]byte, error)
	WriteBytes(ctx context.Context, val []byte) (hash.Hash, error)
}

// ImmutableValue represents a content-addressed value stored in a ValueStore.
// The contents are loaded lazily and stored in |Buf|
type ImmutableValue struct {
	Addr hash.Hash
	Buf  []byte
	vs   ValueStore
}

func NewImmutableValue(addr hash.Hash, vs ValueStore) ImmutableValue {
	return ImmutableValue{Addr: addr, vs: vs}
}

func (t *ImmutableValue) GetBytes(ctx context.Context) ([]byte, error) {
	if t.Buf == nil {
		if t.Addr.IsEmpty() {
			t.Buf = []byte{}
			return t.Buf, nil
		}
		buf, err := t.vs.ReadBytes(ctx, t.Addr)
		if err != nil {
			return nil, err
		}
		t.Buf = buf
	}
	return t.Buf, nil
}

type TextStorage struct {
	ImmutableValue
	maxByteLength int64
}

func (t TextStorage) IsExactLength() bool {
	return false
}

func (t TextStorage) MaxByteLength() int64 {
	return t.maxByteLength
}

func (t TextStorage) Unwrap(ctx context.Context) (string, error) {
	buf, err := t.GetBytes(ctx)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

func (t TextStorage) UnwrapAny(ctx context.Context) (interface{}, error) {
	return t.Unwrap(ctx)
}

func (t TextStorage) WithMaxByteLength(maxByteLength int64) *TextStorage {
	return &TextStorage{
		ImmutableValue: NewImmutableValue(t.Addr, t.vs),
		maxByteLength:  maxByteLength,
	}
}

func (t TextStorage) Compare(ctx context.Context, other interface{}) (cmp int, comparable bool, err error) {
	otherTextStorage, ok := other.(TextStorage)
	if !ok {
		return 0, false, nil
	}
	if otherTextStorage.Addr == t.Addr {
		return 0, true, nil
	}
	return 0, false, nil
}

func NewTextStorage(addr hash.Hash, vs ValueStore) *TextStorage {
	return &TextStorage{
		ImmutableValue: NewImmutableValue(addr, vs),
		maxByteLength:  types.LongText.MaxByteLength(),
	}
}

var _ sql.StringWrapper = &TextStorage{}

type ByteArray struct {
	ImmutableValue
	maxByteLength int64
}

func (b *ByteArray) IsExactLength() bool {
	return false
}

func (b *ByteArray) MaxByteLength() int64 {
	return b.maxByteLength
}

func (b *ByteArray) Compare(ctx context.Context, other interface{}) (cmp int, comparable bool, err error) {
	otherByteArray, ok := other.(ByteArray)
	if !ok {
		return 0, false, nil
	}
	if otherByteArray.Addr == b.Addr {
		return 0, true, nil
	}
	return 0, false, nil
}

func NewByteArray(addr hash.Hash, vs ValueStore) *ByteArray {
	return &ByteArray{
		ImmutableValue: NewImmutableValue(addr, vs),
		maxByteLength:  types.LongBlob.MaxByteLength(),
	}
}

func (b *ByteArray) ToBytes(ctx context.Context) ([]byte, error) {
	return b.GetBytes(ctx)
}

func (b *ByteArray) UnwrapAny(ctx context.Context) (interface{}, error) {
	return b.ToBytes(ctx)
}

func (b ByteArray) Unwrap(ctx context.Context) ([]byte, error) {
	return b.GetBytes(ctx)
}

func (b *ByteArray) ToString(ctx context.Context) (string, error) {
	buf, err := b.ToBytes(ctx)
	if err != nil {
		return "", err
	}
	toShow := BytePeekLength
	if len(buf) < toShow {
		toShow = len(buf)
	}
	return string(buf[:toShow]), nil
}

var _ sql.BytesWrapper = &ByteArray{}

func (t ByteArray) WithMaxByteLength(maxByteLength int64) *ByteArray {
	return &ByteArray{
		ImmutableValue: t.ImmutableValue,
		maxByteLength:  maxByteLength,
	}
}
