// This file was generated by nomdl/codegen.

package test

import (
	"github.com/attic-labs/noms/ref"
	"github.com/attic-labs/noms/types"
)

var __testPackageInFile_struct_recursive_CachedRef = __testPackageInFile_struct_recursive_Ref()

// This function builds up a Noms value that describes the type
// package implemented by this file and registers it with the global
// type package definition cache.
func __testPackageInFile_struct_recursive_Ref() ref.Ref {
	p := types.NewPackage([]types.TypeRef{
		types.MakeStructTypeRef("Tree",
			[]types.Field{
				types.Field{"children", types.MakeCompoundTypeRef("", types.ListKind, types.MakeTypeRef(ref.Ref{}, 0)), false},
			},
			types.Choices{},
		),
	}, []ref.Ref{})
	return types.RegisterPackage(&p)
}

// Tree

type Tree struct {
	m   types.Map
	ref *ref.Ref
}

func NewTree() Tree {
	return Tree{types.NewMap(
		types.NewString("$type"), types.MakeTypeRef(__testPackageInFile_struct_recursive_CachedRef, 0),
		types.NewString("children"), NewListOfTree(),
	), &ref.Ref{}}
}

type TreeDef struct {
	Children ListOfTreeDef
}

func (def TreeDef) New() Tree {
	return Tree{
		types.NewMap(
			types.NewString("$type"), types.MakeTypeRef(__testPackageInFile_struct_recursive_CachedRef, 0),
			types.NewString("children"), def.Children.New(),
		), &ref.Ref{}}
}

func (s Tree) Def() (d TreeDef) {
	d.Children = s.m.Get(types.NewString("children")).(ListOfTree).Def()
	return
}

var __typeRefForTree = types.MakeTypeRef(__testPackageInFile_struct_recursive_CachedRef, 0)

func (m Tree) TypeRef() types.TypeRef {
	return __typeRefForTree
}

func init() {
	types.RegisterFromValFunction(__typeRefForTree, func(v types.Value) types.Value {
		return TreeFromVal(v)
	})
}

func TreeFromVal(val types.Value) Tree {
	// TODO: Do we still need FromVal?
	if val, ok := val.(Tree); ok {
		return val
	}
	// TODO: Validate here
	return Tree{val.(types.Map), &ref.Ref{}}
}

func (s Tree) NomsValue() types.Value {
	// TODO: Remove this
	return s
}

func (s Tree) InternalImplementation() types.Map {
	return s.m
}

func (s Tree) Equals(other types.Value) bool {
	if other, ok := other.(Tree); ok {
		return s.Ref() == other.Ref()
	}
	return false
}

func (s Tree) Ref() ref.Ref {
	return types.EnsureRef(s.ref, s)
}

func (s Tree) Chunks() (futures []types.Future) {
	futures = append(futures, s.TypeRef().Chunks()...)
	futures = append(futures, s.m.Chunks()...)
	return
}

func (s Tree) Children() ListOfTree {
	return s.m.Get(types.NewString("children")).(ListOfTree)
}

func (s Tree) SetChildren(val ListOfTree) Tree {
	return Tree{s.m.Set(types.NewString("children"), val), &ref.Ref{}}
}

// ListOfTree

type ListOfTree struct {
	l   types.List
	ref *ref.Ref
}

func NewListOfTree() ListOfTree {
	return ListOfTree{types.NewList(), &ref.Ref{}}
}

type ListOfTreeDef []TreeDef

func (def ListOfTreeDef) New() ListOfTree {
	l := make([]types.Value, len(def))
	for i, d := range def {
		l[i] = d.New()
	}
	return ListOfTree{types.NewList(l...), &ref.Ref{}}
}

func (l ListOfTree) Def() ListOfTreeDef {
	d := make([]TreeDef, l.Len())
	for i := uint64(0); i < l.Len(); i++ {
		d[i] = l.l.Get(i).(Tree).Def()
	}
	return d
}

func ListOfTreeFromVal(val types.Value) ListOfTree {
	// TODO: Do we still need FromVal?
	if val, ok := val.(ListOfTree); ok {
		return val
	}
	// TODO: Validate here
	return ListOfTree{val.(types.List), &ref.Ref{}}
}

func (l ListOfTree) NomsValue() types.Value {
	// TODO: Remove this
	return l
}

func (l ListOfTree) InternalImplementation() types.List {
	return l.l
}

func (l ListOfTree) Equals(other types.Value) bool {
	if other, ok := other.(ListOfTree); ok {
		return l.Ref() == other.Ref()
	}
	return false
}

func (l ListOfTree) Ref() ref.Ref {
	return types.EnsureRef(l.ref, l)
}

func (l ListOfTree) Chunks() (futures []types.Future) {
	futures = append(futures, l.TypeRef().Chunks()...)
	futures = append(futures, l.l.Chunks()...)
	return
}

// A Noms Value that describes ListOfTree.
var __typeRefForListOfTree types.TypeRef

func (m ListOfTree) TypeRef() types.TypeRef {
	return __typeRefForListOfTree
}

func init() {
	__typeRefForListOfTree = types.MakeCompoundTypeRef("", types.ListKind, types.MakeTypeRef(__testPackageInFile_struct_recursive_CachedRef, 0))
	types.RegisterFromValFunction(__typeRefForListOfTree, func(v types.Value) types.Value {
		return ListOfTreeFromVal(v)
	})
}

func (l ListOfTree) Len() uint64 {
	return l.l.Len()
}

func (l ListOfTree) Empty() bool {
	return l.Len() == uint64(0)
}

func (l ListOfTree) Get(i uint64) Tree {
	return l.l.Get(i).(Tree)
}

func (l ListOfTree) Slice(idx uint64, end uint64) ListOfTree {
	return ListOfTree{l.l.Slice(idx, end), &ref.Ref{}}
}

func (l ListOfTree) Set(i uint64, val Tree) ListOfTree {
	return ListOfTree{l.l.Set(i, val), &ref.Ref{}}
}

func (l ListOfTree) Append(v ...Tree) ListOfTree {
	return ListOfTree{l.l.Append(l.fromElemSlice(v)...), &ref.Ref{}}
}

func (l ListOfTree) Insert(idx uint64, v ...Tree) ListOfTree {
	return ListOfTree{l.l.Insert(idx, l.fromElemSlice(v)...), &ref.Ref{}}
}

func (l ListOfTree) Remove(idx uint64, end uint64) ListOfTree {
	return ListOfTree{l.l.Remove(idx, end), &ref.Ref{}}
}

func (l ListOfTree) RemoveAt(idx uint64) ListOfTree {
	return ListOfTree{(l.l.RemoveAt(idx)), &ref.Ref{}}
}

func (l ListOfTree) fromElemSlice(p []Tree) []types.Value {
	r := make([]types.Value, len(p))
	for i, v := range p {
		r[i] = v
	}
	return r
}

type ListOfTreeIterCallback func(v Tree, i uint64) (stop bool)

func (l ListOfTree) Iter(cb ListOfTreeIterCallback) {
	l.l.Iter(func(v types.Value, i uint64) bool {
		return cb(v.(Tree), i)
	})
}

type ListOfTreeIterAllCallback func(v Tree, i uint64)

func (l ListOfTree) IterAll(cb ListOfTreeIterAllCallback) {
	l.l.IterAll(func(v types.Value, i uint64) {
		cb(v.(Tree), i)
	})
}

type ListOfTreeFilterCallback func(v Tree, i uint64) (keep bool)

func (l ListOfTree) Filter(cb ListOfTreeFilterCallback) ListOfTree {
	nl := NewListOfTree()
	l.IterAll(func(v Tree, i uint64) {
		if cb(v, i) {
			nl = nl.Append(v)
		}
	})
	return nl
}
