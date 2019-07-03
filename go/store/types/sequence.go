// Copyright 2016 Attic Labs, Inc. All rights reserved.
// Licensed under the Apache License, version 2.0:
// http://www.apache.org/licenses/LICENSE-2.0

package types

import (
	"context"

	"github.com/liquidata-inc/ld/dolt/go/store/d"
	"github.com/liquidata-inc/ld/dolt/go/store/hash"
)

type sequenceItem interface{}

type compareFn func(x int, y int) bool

type sequence interface {
	asValueImpl() valueImpl
	cumulativeNumberOfLeaves(idx int) uint64
	Empty() bool
	Equals(format *Format, other Value) bool
	getChildSequence(ctx context.Context, idx int) sequence
	getCompareFn(f *Format, other sequence) compareFn
	getCompositeChildSequence(ctx context.Context, start uint64, length uint64) sequence
	getItem(idx int, f *Format) sequenceItem
	Hash(*Format) hash.Hash
	isLeaf() bool
	Kind() NomsKind
	Len() uint64
	Less(f *Format, other LesserValuable) bool
	numLeaves() uint64
	seqLen() int
	treeLevel() uint64
	typeOf() *Type
	valueBytes(*Format) []byte
	valueReadWriter() ValueReadWriter
	valuesSlice(f *Format, from, to uint64) []Value
	WalkRefs(f *Format, cb RefCallback)
	writeTo(nomsWriter, *Format)
}

const (
	sequencePartKind   = 0
	sequencePartLevel  = 1
	sequencePartCount  = 2
	sequencePartValues = 3
)

type sequenceImpl struct {
	valueImpl
	len uint64
}

func newSequenceImpl(vrw ValueReadWriter, f *Format, buff []byte, offsets []uint32, len uint64) sequenceImpl {
	return sequenceImpl{valueImpl{vrw, f, buff, offsets}, len}
}

func (seq sequenceImpl) decoderSkipToValues() (valueDecoder, uint64) {
	dec := seq.decoderAtPart(sequencePartCount)
	count := dec.readCount()
	return dec, count
}

func (seq sequenceImpl) decoderAtPart(part uint32) valueDecoder {
	offset := seq.offsets[part] - seq.offsets[sequencePartKind]
	return newValueDecoder(seq.buff[offset:], seq.vrw)
}

func (seq sequenceImpl) Empty() bool {
	return seq.Len() == 0
}

func (seq sequenceImpl) Len() uint64 {
	return seq.len
}

func (seq sequenceImpl) seqLen() int {
	_, count := seq.decoderSkipToValues()
	return int(count)
}

func (seq sequenceImpl) getItemOffset(idx int) int {
	// kind, level, count, elements...
	// 0     1      2      3          n+1
	d.PanicIfTrue(idx+sequencePartValues+1 > len(seq.offsets))
	return int(seq.offsets[idx+sequencePartValues] - seq.offsets[sequencePartKind])
}

func (seq sequenceImpl) decoderSkipToIndex(idx int) valueDecoder {
	offset := seq.getItemOffset(idx)
	return seq.decoderAtOffset(offset)
}
