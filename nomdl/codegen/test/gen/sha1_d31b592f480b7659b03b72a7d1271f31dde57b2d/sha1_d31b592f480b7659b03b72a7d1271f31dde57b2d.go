// This file was generated by nomdl/codegen.

package sha1_d31b592f480b7659b03b72a7d1271f31dde57b2d

import (
	"github.com/attic-labs/noms/nomdl/codegen/test/gen/sha1_bbf9c3d7eb6ed891f4b8490b5b81f21f89f7d288"

	"github.com/attic-labs/noms/ref"
	"github.com/attic-labs/noms/types"
)

var __sha1_d31b592f480b7659b03b72a7d1271f31dde57b2dPackageInFile_sha1_d31b592f480b7659b03b72a7d1271f31dde57b2d_CachedRef = __sha1_d31b592f480b7659b03b72a7d1271f31dde57b2dPackageInFile_sha1_d31b592f480b7659b03b72a7d1271f31dde57b2d_Ref()

// This function builds up a Noms value that describes the type
// package implemented by this file and registers it with the global
// type package definition cache.
func __sha1_d31b592f480b7659b03b72a7d1271f31dde57b2dPackageInFile_sha1_d31b592f480b7659b03b72a7d1271f31dde57b2d_Ref() ref.Ref {
	p := types.NewPackage([]types.TypeRef{
		types.MakeStructTypeRef("D",
			[]types.Field{
				types.Field{"structField", types.MakeTypeRef(ref.Parse("sha1-bbf9c3d7eb6ed891f4b8490b5b81f21f89f7d288"), 0), false},
				types.Field{"enumField", types.MakeTypeRef(ref.Parse("sha1-bbf9c3d7eb6ed891f4b8490b5b81f21f89f7d288"), 1), false},
			},
			types.Choices{},
		),
		types.MakeStructTypeRef("DUser",
			[]types.Field{
				types.Field{"Dfield", types.MakeTypeRef(ref.Ref{}, 0), false},
			},
			types.Choices{},
		),
	}, []ref.Ref{
		ref.Parse("sha1-bbf9c3d7eb6ed891f4b8490b5b81f21f89f7d288"),
	})
	return types.RegisterPackage(&p)
}

// D

type D struct {
	m   types.Map
	ref *ref.Ref
}

func NewD() D {
	return D{types.NewMap(
		types.NewString("$type"), types.MakeTypeRef(__sha1_d31b592f480b7659b03b72a7d1271f31dde57b2dPackageInFile_sha1_d31b592f480b7659b03b72a7d1271f31dde57b2d_CachedRef, 0),
		types.NewString("structField"), sha1_bbf9c3d7eb6ed891f4b8490b5b81f21f89f7d288.NewS(),
		types.NewString("enumField"), types.UInt32(0),
	), &ref.Ref{}}
}

type DDef struct {
	StructField sha1_bbf9c3d7eb6ed891f4b8490b5b81f21f89f7d288.SDef
	EnumField   sha1_bbf9c3d7eb6ed891f4b8490b5b81f21f89f7d288.E
}

func (def DDef) New() D {
	return D{
		types.NewMap(
			types.NewString("$type"), types.MakeTypeRef(__sha1_d31b592f480b7659b03b72a7d1271f31dde57b2dPackageInFile_sha1_d31b592f480b7659b03b72a7d1271f31dde57b2d_CachedRef, 0),
			types.NewString("structField"), def.StructField.New(),
			types.NewString("enumField"), types.UInt32(def.EnumField),
		), &ref.Ref{}}
}

func (s D) Def() (d DDef) {
	d.StructField = s.m.Get(types.NewString("structField")).(sha1_bbf9c3d7eb6ed891f4b8490b5b81f21f89f7d288.S).Def()
	d.EnumField = sha1_bbf9c3d7eb6ed891f4b8490b5b81f21f89f7d288.E(s.m.Get(types.NewString("enumField")).(types.UInt32))
	return
}

var __typeRefForD = types.MakeTypeRef(__sha1_d31b592f480b7659b03b72a7d1271f31dde57b2dPackageInFile_sha1_d31b592f480b7659b03b72a7d1271f31dde57b2d_CachedRef, 0)

func (m D) TypeRef() types.TypeRef {
	return __typeRefForD
}

func init() {
	types.RegisterFromValFunction(__typeRefForD, func(v types.Value) types.Value {
		return DFromVal(v)
	})
}

func DFromVal(val types.Value) D {
	// TODO: Do we still need FromVal?
	if val, ok := val.(D); ok {
		return val
	}
	// TODO: Validate here
	return D{val.(types.Map), &ref.Ref{}}
}

func (s D) NomsValue() types.Value {
	// TODO: Remove this
	return s
}

func (s D) InternalImplementation() types.Map {
	return s.m
}

func (s D) Equals(other types.Value) bool {
	if other, ok := other.(D); ok {
		return s.Ref() == other.Ref()
	}
	return false
}

func (s D) Ref() ref.Ref {
	return types.EnsureRef(s.ref, s)
}

func (s D) Chunks() (futures []types.Future) {
	futures = append(futures, s.TypeRef().Chunks()...)
	futures = append(futures, s.m.Chunks()...)
	return
}

func (s D) StructField() sha1_bbf9c3d7eb6ed891f4b8490b5b81f21f89f7d288.S {
	return s.m.Get(types.NewString("structField")).(sha1_bbf9c3d7eb6ed891f4b8490b5b81f21f89f7d288.S)
}

func (s D) SetStructField(val sha1_bbf9c3d7eb6ed891f4b8490b5b81f21f89f7d288.S) D {
	return D{s.m.Set(types.NewString("structField"), val), &ref.Ref{}}
}

func (s D) EnumField() sha1_bbf9c3d7eb6ed891f4b8490b5b81f21f89f7d288.E {
	return sha1_bbf9c3d7eb6ed891f4b8490b5b81f21f89f7d288.E(s.m.Get(types.NewString("enumField")).(types.UInt32))
}

func (s D) SetEnumField(val sha1_bbf9c3d7eb6ed891f4b8490b5b81f21f89f7d288.E) D {
	return D{s.m.Set(types.NewString("enumField"), types.UInt32(val)), &ref.Ref{}}
}

// DUser

type DUser struct {
	m   types.Map
	ref *ref.Ref
}

func NewDUser() DUser {
	return DUser{types.NewMap(
		types.NewString("$type"), types.MakeTypeRef(__sha1_d31b592f480b7659b03b72a7d1271f31dde57b2dPackageInFile_sha1_d31b592f480b7659b03b72a7d1271f31dde57b2d_CachedRef, 1),
		types.NewString("Dfield"), NewD(),
	), &ref.Ref{}}
}

type DUserDef struct {
	Dfield DDef
}

func (def DUserDef) New() DUser {
	return DUser{
		types.NewMap(
			types.NewString("$type"), types.MakeTypeRef(__sha1_d31b592f480b7659b03b72a7d1271f31dde57b2dPackageInFile_sha1_d31b592f480b7659b03b72a7d1271f31dde57b2d_CachedRef, 1),
			types.NewString("Dfield"), def.Dfield.New(),
		), &ref.Ref{}}
}

func (s DUser) Def() (d DUserDef) {
	d.Dfield = s.m.Get(types.NewString("Dfield")).(D).Def()
	return
}

var __typeRefForDUser = types.MakeTypeRef(__sha1_d31b592f480b7659b03b72a7d1271f31dde57b2dPackageInFile_sha1_d31b592f480b7659b03b72a7d1271f31dde57b2d_CachedRef, 1)

func (m DUser) TypeRef() types.TypeRef {
	return __typeRefForDUser
}

func init() {
	types.RegisterFromValFunction(__typeRefForDUser, func(v types.Value) types.Value {
		return DUserFromVal(v)
	})
}

func DUserFromVal(val types.Value) DUser {
	// TODO: Do we still need FromVal?
	if val, ok := val.(DUser); ok {
		return val
	}
	// TODO: Validate here
	return DUser{val.(types.Map), &ref.Ref{}}
}

func (s DUser) NomsValue() types.Value {
	// TODO: Remove this
	return s
}

func (s DUser) InternalImplementation() types.Map {
	return s.m
}

func (s DUser) Equals(other types.Value) bool {
	if other, ok := other.(DUser); ok {
		return s.Ref() == other.Ref()
	}
	return false
}

func (s DUser) Ref() ref.Ref {
	return types.EnsureRef(s.ref, s)
}

func (s DUser) Chunks() (futures []types.Future) {
	futures = append(futures, s.TypeRef().Chunks()...)
	futures = append(futures, s.m.Chunks()...)
	return
}

func (s DUser) Dfield() D {
	return s.m.Get(types.NewString("Dfield")).(D)
}

func (s DUser) SetDfield(val D) DUser {
	return DUser{s.m.Set(types.NewString("Dfield"), val), &ref.Ref{}}
}
