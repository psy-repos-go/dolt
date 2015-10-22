// This file was generated by nomdl/codegen.

package test

import (
	"github.com/attic-labs/noms/ref"
	"github.com/attic-labs/noms/types"
)

var __testPackageInFile_struct_primitives_CachedRef = __testPackageInFile_struct_primitives_Ref()

// This function builds up a Noms value that describes the type
// package implemented by this file and registers it with the global
// type package definition cache.
func __testPackageInFile_struct_primitives_Ref() ref.Ref {
	p := types.NewPackage([]types.TypeRef{
		types.MakeStructTypeRef("StructPrimitives",
			[]types.Field{
				types.Field{"uint64", types.MakePrimitiveTypeRef(types.UInt64Kind), false},
				types.Field{"uint32", types.MakePrimitiveTypeRef(types.UInt32Kind), false},
				types.Field{"uint16", types.MakePrimitiveTypeRef(types.UInt16Kind), false},
				types.Field{"uint8", types.MakePrimitiveTypeRef(types.UInt8Kind), false},
				types.Field{"int64", types.MakePrimitiveTypeRef(types.Int64Kind), false},
				types.Field{"int32", types.MakePrimitiveTypeRef(types.Int32Kind), false},
				types.Field{"int16", types.MakePrimitiveTypeRef(types.Int16Kind), false},
				types.Field{"int8", types.MakePrimitiveTypeRef(types.Int8Kind), false},
				types.Field{"float64", types.MakePrimitiveTypeRef(types.Float64Kind), false},
				types.Field{"float32", types.MakePrimitiveTypeRef(types.Float32Kind), false},
				types.Field{"bool", types.MakePrimitiveTypeRef(types.BoolKind), false},
				types.Field{"string", types.MakePrimitiveTypeRef(types.StringKind), false},
				types.Field{"blob", types.MakePrimitiveTypeRef(types.BlobKind), false},
				types.Field{"value", types.MakePrimitiveTypeRef(types.ValueKind), false},
			},
			types.Choices{},
		),
	}, []ref.Ref{})
	return types.RegisterPackage(&p)
}

// StructPrimitives

type StructPrimitives struct {
	m   types.Map
	ref *ref.Ref
}

func NewStructPrimitives() StructPrimitives {
	return StructPrimitives{types.NewMap(
		types.NewString("$type"), types.MakeTypeRef(__testPackageInFile_struct_primitives_CachedRef, 0),
		types.NewString("uint64"), types.UInt64(0),
		types.NewString("uint32"), types.UInt32(0),
		types.NewString("uint16"), types.UInt16(0),
		types.NewString("uint8"), types.UInt8(0),
		types.NewString("int64"), types.Int64(0),
		types.NewString("int32"), types.Int32(0),
		types.NewString("int16"), types.Int16(0),
		types.NewString("int8"), types.Int8(0),
		types.NewString("float64"), types.Float64(0),
		types.NewString("float32"), types.Float32(0),
		types.NewString("bool"), types.Bool(false),
		types.NewString("string"), types.NewString(""),
		types.NewString("blob"), types.NewEmptyBlob(),
		types.NewString("value"), types.Bool(false),
	), &ref.Ref{}}
}

type StructPrimitivesDef struct {
	Uint64  uint64
	Uint32  uint32
	Uint16  uint16
	Uint8   uint8
	Int64   int64
	Int32   int32
	Int16   int16
	Int8    int8
	Float64 float64
	Float32 float32
	Bool    bool
	String  string
	Blob    types.Blob
	Value   types.Value
}

func (def StructPrimitivesDef) New() StructPrimitives {
	return StructPrimitives{
		types.NewMap(
			types.NewString("$type"), types.MakeTypeRef(__testPackageInFile_struct_primitives_CachedRef, 0),
			types.NewString("uint64"), types.UInt64(def.Uint64),
			types.NewString("uint32"), types.UInt32(def.Uint32),
			types.NewString("uint16"), types.UInt16(def.Uint16),
			types.NewString("uint8"), types.UInt8(def.Uint8),
			types.NewString("int64"), types.Int64(def.Int64),
			types.NewString("int32"), types.Int32(def.Int32),
			types.NewString("int16"), types.Int16(def.Int16),
			types.NewString("int8"), types.Int8(def.Int8),
			types.NewString("float64"), types.Float64(def.Float64),
			types.NewString("float32"), types.Float32(def.Float32),
			types.NewString("bool"), types.Bool(def.Bool),
			types.NewString("string"), types.NewString(def.String),
			types.NewString("blob"), def.Blob,
			types.NewString("value"), def.Value,
		), &ref.Ref{}}
}

func (s StructPrimitives) Def() (d StructPrimitivesDef) {
	d.Uint64 = uint64(s.m.Get(types.NewString("uint64")).(types.UInt64))
	d.Uint32 = uint32(s.m.Get(types.NewString("uint32")).(types.UInt32))
	d.Uint16 = uint16(s.m.Get(types.NewString("uint16")).(types.UInt16))
	d.Uint8 = uint8(s.m.Get(types.NewString("uint8")).(types.UInt8))
	d.Int64 = int64(s.m.Get(types.NewString("int64")).(types.Int64))
	d.Int32 = int32(s.m.Get(types.NewString("int32")).(types.Int32))
	d.Int16 = int16(s.m.Get(types.NewString("int16")).(types.Int16))
	d.Int8 = int8(s.m.Get(types.NewString("int8")).(types.Int8))
	d.Float64 = float64(s.m.Get(types.NewString("float64")).(types.Float64))
	d.Float32 = float32(s.m.Get(types.NewString("float32")).(types.Float32))
	d.Bool = bool(s.m.Get(types.NewString("bool")).(types.Bool))
	d.String = s.m.Get(types.NewString("string")).(types.String).String()
	d.Blob = s.m.Get(types.NewString("blob")).(types.Blob)
	d.Value = s.m.Get(types.NewString("value"))
	return
}

var __typeRefForStructPrimitives = types.MakeTypeRef(__testPackageInFile_struct_primitives_CachedRef, 0)

func (m StructPrimitives) TypeRef() types.TypeRef {
	return __typeRefForStructPrimitives
}

func init() {
	types.RegisterFromValFunction(__typeRefForStructPrimitives, func(v types.Value) types.Value {
		return StructPrimitivesFromVal(v)
	})
}

func StructPrimitivesFromVal(val types.Value) StructPrimitives {
	// TODO: Do we still need FromVal?
	if val, ok := val.(StructPrimitives); ok {
		return val
	}
	// TODO: Validate here
	return StructPrimitives{val.(types.Map), &ref.Ref{}}
}

func (s StructPrimitives) NomsValue() types.Value {
	// TODO: Remove this
	return s
}

func (s StructPrimitives) InternalImplementation() types.Map {
	return s.m
}

func (s StructPrimitives) Equals(other types.Value) bool {
	if other, ok := other.(StructPrimitives); ok {
		return s.Ref() == other.Ref()
	}
	return false
}

func (s StructPrimitives) Ref() ref.Ref {
	return types.EnsureRef(s.ref, s)
}

func (s StructPrimitives) Chunks() (futures []types.Future) {
	futures = append(futures, s.TypeRef().Chunks()...)
	futures = append(futures, s.m.Chunks()...)
	return
}

func (s StructPrimitives) Uint64() uint64 {
	return uint64(s.m.Get(types.NewString("uint64")).(types.UInt64))
}

func (s StructPrimitives) SetUint64(val uint64) StructPrimitives {
	return StructPrimitives{s.m.Set(types.NewString("uint64"), types.UInt64(val)), &ref.Ref{}}
}

func (s StructPrimitives) Uint32() uint32 {
	return uint32(s.m.Get(types.NewString("uint32")).(types.UInt32))
}

func (s StructPrimitives) SetUint32(val uint32) StructPrimitives {
	return StructPrimitives{s.m.Set(types.NewString("uint32"), types.UInt32(val)), &ref.Ref{}}
}

func (s StructPrimitives) Uint16() uint16 {
	return uint16(s.m.Get(types.NewString("uint16")).(types.UInt16))
}

func (s StructPrimitives) SetUint16(val uint16) StructPrimitives {
	return StructPrimitives{s.m.Set(types.NewString("uint16"), types.UInt16(val)), &ref.Ref{}}
}

func (s StructPrimitives) Uint8() uint8 {
	return uint8(s.m.Get(types.NewString("uint8")).(types.UInt8))
}

func (s StructPrimitives) SetUint8(val uint8) StructPrimitives {
	return StructPrimitives{s.m.Set(types.NewString("uint8"), types.UInt8(val)), &ref.Ref{}}
}

func (s StructPrimitives) Int64() int64 {
	return int64(s.m.Get(types.NewString("int64")).(types.Int64))
}

func (s StructPrimitives) SetInt64(val int64) StructPrimitives {
	return StructPrimitives{s.m.Set(types.NewString("int64"), types.Int64(val)), &ref.Ref{}}
}

func (s StructPrimitives) Int32() int32 {
	return int32(s.m.Get(types.NewString("int32")).(types.Int32))
}

func (s StructPrimitives) SetInt32(val int32) StructPrimitives {
	return StructPrimitives{s.m.Set(types.NewString("int32"), types.Int32(val)), &ref.Ref{}}
}

func (s StructPrimitives) Int16() int16 {
	return int16(s.m.Get(types.NewString("int16")).(types.Int16))
}

func (s StructPrimitives) SetInt16(val int16) StructPrimitives {
	return StructPrimitives{s.m.Set(types.NewString("int16"), types.Int16(val)), &ref.Ref{}}
}

func (s StructPrimitives) Int8() int8 {
	return int8(s.m.Get(types.NewString("int8")).(types.Int8))
}

func (s StructPrimitives) SetInt8(val int8) StructPrimitives {
	return StructPrimitives{s.m.Set(types.NewString("int8"), types.Int8(val)), &ref.Ref{}}
}

func (s StructPrimitives) Float64() float64 {
	return float64(s.m.Get(types.NewString("float64")).(types.Float64))
}

func (s StructPrimitives) SetFloat64(val float64) StructPrimitives {
	return StructPrimitives{s.m.Set(types.NewString("float64"), types.Float64(val)), &ref.Ref{}}
}

func (s StructPrimitives) Float32() float32 {
	return float32(s.m.Get(types.NewString("float32")).(types.Float32))
}

func (s StructPrimitives) SetFloat32(val float32) StructPrimitives {
	return StructPrimitives{s.m.Set(types.NewString("float32"), types.Float32(val)), &ref.Ref{}}
}

func (s StructPrimitives) Bool() bool {
	return bool(s.m.Get(types.NewString("bool")).(types.Bool))
}

func (s StructPrimitives) SetBool(val bool) StructPrimitives {
	return StructPrimitives{s.m.Set(types.NewString("bool"), types.Bool(val)), &ref.Ref{}}
}

func (s StructPrimitives) String() string {
	return s.m.Get(types.NewString("string")).(types.String).String()
}

func (s StructPrimitives) SetString(val string) StructPrimitives {
	return StructPrimitives{s.m.Set(types.NewString("string"), types.NewString(val)), &ref.Ref{}}
}

func (s StructPrimitives) Blob() types.Blob {
	return s.m.Get(types.NewString("blob")).(types.Blob)
}

func (s StructPrimitives) SetBlob(val types.Blob) StructPrimitives {
	return StructPrimitives{s.m.Set(types.NewString("blob"), val), &ref.Ref{}}
}

func (s StructPrimitives) Value() types.Value {
	return s.m.Get(types.NewString("value"))
}

func (s StructPrimitives) SetValue(val types.Value) StructPrimitives {
	return StructPrimitives{s.m.Set(types.NewString("value"), val), &ref.Ref{}}
}
