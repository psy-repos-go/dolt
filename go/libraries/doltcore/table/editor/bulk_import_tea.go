// Copyright 2021 Dolthub, Inc.
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

package editor

import (
	"context"

	"github.com/dolthub/dolt/go/store/hash"
	"github.com/dolthub/dolt/go/store/types"
	"github.com/dolthub/dolt/go/store/types/edits"
)

var _ TableEditAccumulator = (*BulkImportTEA)(nil)

type BulkImportTEA struct {
	teaf       TEAFactory
	emptyTuple types.Tuple

	ea      types.EditAccumulator
	rowData types.Map

	// opCount contains the number of edits that would be applied in materializing the edits
	opCount int64
	adds    map[hash.Hash]bool
	deletes map[hash.Hash]bool
}

func (tea *BulkImportTEA) OpCount() int64 {
	return tea.opCount
}

func (tea *BulkImportTEA) Delete(keyHash hash.Hash, key types.Tuple) {
	tea.opCount++
	tea.ea.AddEdit(key, nil)

	tea.deletes[keyHash] = true
	delete(tea.adds, keyHash)
}

func (tea *BulkImportTEA) Insert(keyHash hash.Hash, key types.Tuple, val types.Tuple) {
	tea.opCount++
	tea.ea.AddEdit(key, val)

	tea.adds[keyHash] = true
	delete(tea.deletes, keyHash)
}

func (tea *BulkImportTEA) Get(ctx context.Context, keyHash hash.Hash, key types.Tuple) (*doltKVP, bool, error) {
	if tea.deletes[keyHash] {
		return nil, false, nil
	}

	if tea.adds[keyHash] {
		return &doltKVP{k: key, v: tea.emptyTuple}, true, nil
	}

	v, ok, err := tea.rowData.MaybeGetTuple(ctx, key)

	if err != nil {
		return nil, false, err
	}

	if !ok {
		return nil, false, nil
	}

	return &doltKVP{k: key, v: v}, true, nil
}

func (tea *BulkImportTEA) Commit(ctx context.Context, nbf *types.NomsBinFormat) error {
	panic("Not Supported")
}

func (tea *BulkImportTEA) Rollback(ctx context.Context) error {
	panic("Not Supported")
}

func (tea *BulkImportTEA) MaterializeEdits(ctx context.Context, nbf *types.NomsBinFormat) (m types.Map, err error) {
	ea := tea.ea
	defer ea.Close(ctx)

	itr, err := ea.FinishedEditing()
	if err != nil {
		return types.EmptyMap, err
	}

	currMap := tea.rowData
	for !itr.ReachedEOF() {
		currMap, _, err = types.ApplyNEdits(ctx, itr, currMap, 256*1024)
		if err != nil {
			return types.EmptyMap, err
		}
	}

	*tea = *(tea.teaf.NewTEA(ctx, currMap).(*BulkImportTEA))
	return currMap, nil
}

var _ TEAFactory = (*BulkImportTEAFactory)(nil)

type BulkImportTEAFactory struct {
	nbf       *types.NomsBinFormat
	vrw       types.ValueReadWriter
	directory string
}

func NewBulkImportTEAFactory(nbf *types.NomsBinFormat, vrw types.ValueReadWriter, directory string) *BulkImportTEAFactory {
	return &BulkImportTEAFactory{
		nbf:       nbf,
		vrw:       vrw,
		directory: directory,
	}
}

func (b *BulkImportTEAFactory) NewTEA(ctx context.Context, rowData types.Map) TableEditAccumulator {
	const flushInterval = 256 * 1024

	createMapEA := func() types.EditAccumulator {
		return types.CreateEditAccForMapEdits(b.nbf)
	}

	ea := edits.NewDiskBackedEditAcc(ctx, b.nbf, b.vrw, flushInterval, b.directory, createMapEA)
	return &BulkImportTEA{
		teaf:       b,
		rowData:    rowData,
		ea:         ea,
		adds:       make(map[hash.Hash]bool),
		deletes:    make(map[hash.Hash]bool),
		emptyTuple: types.EmptyTuple(b.nbf),
	}
}
