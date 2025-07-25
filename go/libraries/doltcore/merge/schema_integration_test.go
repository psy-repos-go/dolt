// Copyright 2020 Dolthub, Inc.
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

package merge_test

import (
	"context"
	"testing"

	"github.com/dolthub/go-mysql-server/sql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dolthub/dolt/go/cmd/dolt/cli"
	"github.com/dolthub/dolt/go/cmd/dolt/commands"
	"github.com/dolthub/dolt/go/libraries/doltcore/doltdb"
	"github.com/dolthub/dolt/go/libraries/doltcore/dtestutils"
	"github.com/dolthub/dolt/go/libraries/doltcore/env"
	"github.com/dolthub/dolt/go/libraries/doltcore/merge"
	"github.com/dolthub/dolt/go/libraries/doltcore/schema"
	"github.com/dolthub/dolt/go/libraries/doltcore/schema/typeinfo"
	"github.com/dolthub/dolt/go/libraries/doltcore/table/editor"
	"github.com/dolthub/dolt/go/store/types"
)

type testCommand struct {
	cmd  cli.Command
	args args
}

func (tc testCommand) exec(t *testing.T, ctx context.Context, dEnv *env.DoltEnv) int {
	cliCtx, err := commands.NewArgFreeCliContext(ctx, dEnv, dEnv.FS)
	require.NoError(t, err)
	return tc.cmd.Exec(ctx, tc.cmd.Name(), tc.args, dEnv, cliCtx)
}

type args []string

// TestMergeSchemas are schema merge integration tests from 2020
func TestMergeSchemas(t *testing.T) {
	for _, test := range mergeSchemaTests {
		t.Run(test.name, func(t *testing.T) {
			testMergeSchemas(t, test)
		})
	}
	for _, test := range mergeSchemaConflictTests {
		t.Run(test.name, func(t *testing.T) {
			testMergeSchemasWithConflicts(t, test)
		})
	}
	for _, test := range mergeForeignKeyTests {
		t.Run(test.name, func(t *testing.T) {
			testMergeForeignKeys(t, test)
		})
	}
}

type mergeSchemaTest struct {
	name  string
	setup []testCommand
	sch   schema.Schema
	skip  bool
}

type mergeSchemaConflictTest struct {
	name        string
	setup       []testCommand
	expConflict merge.SchemaConflict
	expectedErr error
}

type mergeForeignKeyTest struct {
	name          string
	setup         []testCommand
	fkColl        *doltdb.ForeignKeyCollection
	expFKConflict []merge.FKConflict
}

var setupCommon = []testCommand{
	{commands.SqlCmd{}, []string{"-q", "create table test (" +
		"pk int not null primary key," +
		"c1 int not null," +
		"c2 int," +
		"c3 int);"}},
	{commands.SqlCmd{}, []string{"-q", "create index c1_idx on test(c1)"}},
	{commands.AddCmd{}, []string{"."}},
	{commands.CommitCmd{}, []string{"-m", "setup common"}},
	{commands.BranchCmd{}, []string{"other"}},
}

var mergeSchemaTests = []mergeSchemaTest{
	{
		name:  "no changes",
		setup: []testCommand{},
		sch: schemaFromColsAndIdxs(
			colCollection(
				newColTypeInfo("pk", uint64(3228), typeinfo.Int32Type, true, schema.NotNullConstraint{}),
				newColTypeInfo("c1", uint64(8201), typeinfo.Int32Type, false, schema.NotNullConstraint{}),
				newColTypeInfo("c2", uint64(8539), typeinfo.Int32Type, false),
				newColTypeInfo("c3", uint64(4696), typeinfo.Int32Type, false)),
			schema.NewIndex("c1_idx", []uint64{8201}, []uint64{8201, 3228}, nil, schema.IndexProperties{IsUserDefined: true}),
		),
	},
	{
		name: "add cols, drop cols, merge",
		setup: []testCommand{
			{commands.SqlCmd{}, []string{"-q", "alter table test drop column c2;"}},
			{commands.SqlCmd{}, []string{"-q", "alter table test add column c8 int;"}},
			{commands.AddCmd{}, []string{"."}},
			{commands.CommitCmd{}, []string{"-m", "modified branch main"}},
			{commands.CheckoutCmd{}, []string{"other"}},
			{commands.SqlCmd{}, []string{"-q", "alter table test drop column c3;"}},
			{commands.SqlCmd{}, []string{"-q", "alter table test add column c9 int;"}},
			{commands.AddCmd{}, []string{"."}},
			{commands.CommitCmd{}, []string{"-m", "modified branch other"}},
			{commands.CheckoutCmd{}, []string{env.DefaultInitBranch}},
		},
		sch: schemaFromColsAndIdxs(
			colCollection(
				newColTypeInfo("pk", uint64(3228), typeinfo.Int32Type, true, schema.NotNullConstraint{}),
				newColTypeInfo("c1", uint64(8201), typeinfo.Int32Type, false, schema.NotNullConstraint{}),
				newColTypeInfo("c8", uint64(12393), typeinfo.Int32Type, false),
				newColTypeInfo("c9", uint64(4508), typeinfo.Int32Type, false)),
			schema.NewIndex("c1_idx", []uint64{8201}, []uint64{8201, 3228}, nil, schema.IndexProperties{IsUserDefined: true}),
		),
	},
	{
		name: "add constraint, drop constraint, merge",
		setup: []testCommand{
			{commands.SqlCmd{}, []string{"-q", "alter table test modify c1 int null;"}},
			{commands.AddCmd{}, []string{"."}},
			{commands.CommitCmd{}, []string{"-m", "modified branch main"}},
			{commands.CheckoutCmd{}, []string{"other"}},
			{commands.SqlCmd{}, []string{"-q", "alter table test modify c2 int not null;"}},
			{commands.AddCmd{}, []string{"."}},
			{commands.CommitCmd{}, []string{"-m", "modified branch other"}},
			{commands.CheckoutCmd{}, []string{env.DefaultInitBranch}},
		},
		sch: schemaFromColsAndIdxs(
			colCollection(
				newColTypeInfo("pk", uint64(3228), typeinfo.Int32Type, true, schema.NotNullConstraint{}),
				newColTypeInfo("c1", uint64(8201), typeinfo.Int32Type, false),
				newColTypeInfo("c2", uint64(8539), typeinfo.Int32Type, false, schema.NotNullConstraint{}),
				newColTypeInfo("c3", uint64(4696), typeinfo.Int32Type, false)),
			schema.NewIndex("c1_idx", []uint64{8201}, []uint64{8201, 3228}, nil, schema.IndexProperties{IsUserDefined: true}),
		),
	},
	{
		name: "add index, drop index, merge",
		setup: []testCommand{
			{commands.SqlCmd{}, []string{"-q", "create index c3_idx on test(c3);"}},
			{commands.AddCmd{}, []string{"."}},
			{commands.CommitCmd{}, []string{"-m", "modified branch main"}},
			{commands.CheckoutCmd{}, []string{"other"}},
			{commands.SqlCmd{}, []string{"-q", "alter table test drop index c1_idx;"}},
			{commands.AddCmd{}, []string{"."}},
			{commands.CommitCmd{}, []string{"-m", "modified branch other"}},
			{commands.CheckoutCmd{}, []string{env.DefaultInitBranch}},
		},
		sch: schemaFromColsAndIdxs(
			colCollection(
				newColTypeInfo("pk", uint64(3228), typeinfo.Int32Type, true, schema.NotNullConstraint{}),
				newColTypeInfo("c1", uint64(8201), typeinfo.Int32Type, false, schema.NotNullConstraint{}),
				newColTypeInfo("c2", uint64(8539), typeinfo.Int32Type, false),
				newColTypeInfo("c3", uint64(4696), typeinfo.Int32Type, false)),
			schema.NewIndex("c3_idx", []uint64{4696}, []uint64{4696, 3228}, nil, schema.IndexProperties{IsUserDefined: true}),
		),
	},
	{
		name: "rename columns",
		setup: []testCommand{
			{commands.SqlCmd{}, []string{"-q", "alter table test rename column c3 to c33;"}},
			{commands.AddCmd{}, []string{"."}},
			{commands.CommitCmd{}, []string{"-m", "modified branch main"}},
			{commands.CheckoutCmd{}, []string{"other"}},
			{commands.SqlCmd{}, []string{"-q", "alter table test rename column c2 to c22;"}},
			{commands.AddCmd{}, []string{"."}},
			{commands.CommitCmd{}, []string{"-m", "modified branch other"}},
			{commands.CheckoutCmd{}, []string{env.DefaultInitBranch}},
		},
		// hmmm, we created new columns with a rename?
		sch: schemaFromColsAndIdxs(
			colCollection(
				newColTypeInfo("pk", uint64(3228), typeinfo.Int32Type, true, schema.NotNullConstraint{}),
				newColTypeInfo("c1", uint64(8201), typeinfo.Int32Type, false, schema.NotNullConstraint{}),
				newColTypeInfo("c22", uint64(8539), typeinfo.Int32Type, false),
				newColTypeInfo("c33", uint64(4696), typeinfo.Int32Type, false)),
			schema.NewIndex("c1_idx", []uint64{8201}, []uint64{8201, 3228}, nil, schema.IndexProperties{IsUserDefined: true}),
		),
	},
	{
		name: "rename indexes",
		setup: []testCommand{
			{commands.SqlCmd{}, []string{"-q", "alter table test drop index c1_idx;"}},
			{commands.SqlCmd{}, []string{"-q", "create index c1_index on test(c1);"}},
			{commands.AddCmd{}, []string{"."}},
			{commands.CommitCmd{}, []string{"-m", "modified branch main"}},
		},
		sch: schemaFromColsAndIdxs(
			colCollection(
				newColTypeInfo("pk", uint64(3228), typeinfo.Int32Type, true, schema.NotNullConstraint{}),
				newColTypeInfo("c1", uint64(8201), typeinfo.Int32Type, false, schema.NotNullConstraint{}),
				newColTypeInfo("c2", uint64(8539), typeinfo.Int32Type, false),
				newColTypeInfo("c3", uint64(4696), typeinfo.Int32Type, false)),
			schema.NewIndex("c1_index", []uint64{8201}, []uint64{8201, 3228}, nil, schema.IndexProperties{IsUserDefined: true}),
		),
	},
	{
		name: "add same column on both branches, merge",
		setup: []testCommand{
			{commands.SqlCmd{}, []string{"-q", "alter table test add column c4 int;"}},
			{commands.AddCmd{}, []string{"."}},
			{commands.CommitCmd{}, []string{"-m", "modified branch main"}},
			{commands.CheckoutCmd{}, []string{"other"}},
			{commands.SqlCmd{}, []string{"-q", "alter table test add column c4 int;"}},
			{commands.AddCmd{}, []string{"."}},
			{commands.CommitCmd{}, []string{"-m", "modified branch other"}},
			{commands.CheckoutCmd{}, []string{env.DefaultInitBranch}},
		},
		sch: schemaFromColsAndIdxs(
			colCollection(
				newColTypeInfo("pk", uint64(3228), typeinfo.Int32Type, true, schema.NotNullConstraint{}),
				newColTypeInfo("c1", uint64(8201), typeinfo.Int32Type, false, schema.NotNullConstraint{}),
				newColTypeInfo("c2", uint64(8539), typeinfo.Int32Type, false),
				newColTypeInfo("c3", uint64(4696), typeinfo.Int32Type, false),
				newColTypeInfo("c4", uint64(1716), typeinfo.Int32Type, false)),
			schema.NewIndex("c1_idx", []uint64{8201}, []uint64{8201, 3228}, nil, schema.IndexProperties{IsUserDefined: true}),
		),
	},
	{
		name: "add same index on both branches, merge",
		setup: []testCommand{
			{commands.SqlCmd{}, []string{"-q", "create index c3_idx on test(c3);"}},
			{commands.AddCmd{}, []string{"."}},
			{commands.CommitCmd{}, []string{"-m", "modified branch main"}},
			{commands.CheckoutCmd{}, []string{"other"}},
			{commands.SqlCmd{}, []string{"-q", "create index c3_idx on test(c3);"}},
			{commands.AddCmd{}, []string{"."}},
			{commands.CommitCmd{}, []string{"-m", "modified branch other"}},
			{commands.CheckoutCmd{}, []string{env.DefaultInitBranch}},
		},
		sch: schemaFromColsAndIdxs(
			colCollection(
				newColTypeInfo("pk", uint64(3228), typeinfo.Int32Type, true, schema.NotNullConstraint{}),
				newColTypeInfo("c1", uint64(8201), typeinfo.Int32Type, false, schema.NotNullConstraint{}),
				newColTypeInfo("c2", uint64(8539), typeinfo.Int32Type, false),
				newColTypeInfo("c3", uint64(4696), typeinfo.Int32Type, false)),
			schema.NewIndex("c1_idx", []uint64{8201}, []uint64{8201, 3228}, nil, schema.IndexProperties{IsUserDefined: true}),
			schema.NewIndex("c3_idx", []uint64{4696}, []uint64{4696, 3228}, nil, schema.IndexProperties{IsUserDefined: true}),
		),
	},
}

var mergeSchemaConflictTests = []mergeSchemaConflictTest{
	{
		name: "no conflicts",
		expConflict: merge.SchemaConflict{
			TableName: doltdb.TableName{Name: "test"},
		},
	},
	{
		name: "column name collisions",
		setup: []testCommand{
			{commands.SqlCmd{}, []string{"-q", "alter table test rename column c3 to c4;"}},
			{commands.SqlCmd{}, []string{"-q", "alter table test add column C6 int;"}},
			{commands.AddCmd{}, []string{"."}},
			{commands.CommitCmd{}, []string{"-m", "modified branch main"}},
			{commands.CheckoutCmd{}, []string{"other"}},
			{commands.SqlCmd{}, []string{"-q", "alter table test rename column c2 to c4;"}},
			{commands.SqlCmd{}, []string{"-q", "alter table test add column c6 int;"}},
			{commands.AddCmd{}, []string{"."}},
			{commands.CommitCmd{}, []string{"-m", "modified branch other"}},
			{commands.CheckoutCmd{}, []string{env.DefaultInitBranch}},
		},
		expConflict: merge.SchemaConflict{
			TableName: doltdb.TableName{Name: "test"},
			ColConflicts: []merge.ColConflict{
				{
					Kind:   merge.NameCollision,
					Ours:   newColTypeInfo("C6", uint64(13258), typeinfo.Int32Type, false),
					Theirs: newColTypeInfo("c6", uint64(13258), typeinfo.Int32Type, false),
				},
				{
					Kind:   merge.NameCollision,
					Ours:   newColTypeInfo("c4", uint64(4696), typeinfo.Int32Type, false),
					Theirs: newColTypeInfo("c4", uint64(8539), typeinfo.Int32Type, false),
				},
			},
		},
	},
	{
		name: "index name collisions",
		setup: []testCommand{
			{commands.SqlCmd{}, []string{"-q", "create index `both` on test (c1,c2);"}},
			{commands.AddCmd{}, []string{"."}},
			{commands.CommitCmd{}, []string{"-m", "modified branch main"}},
			{commands.CheckoutCmd{}, []string{"other"}},
			{commands.SqlCmd{}, []string{"-q", "create index `both` on test (c2, c3);"}},
			{commands.AddCmd{}, []string{"."}},
			{commands.CommitCmd{}, []string{"-m", "modified branch other"}},
			{commands.CheckoutCmd{}, []string{env.DefaultInitBranch}},
		},
		expConflict: merge.SchemaConflict{
			TableName: doltdb.TableName{Name: "test"},
			IdxConflicts: []merge.IdxConflict{
				{
					Kind:   merge.NameCollision,
					Ours:   schema.NewIndex("both", []uint64{8201, 8539}, []uint64{8201, 8539, 3228}, nil, schema.IndexProperties{IsUserDefined: true}),
					Theirs: schema.NewIndex("both", []uint64{8539, 4696}, []uint64{8539, 4696, 3228}, nil, schema.IndexProperties{IsUserDefined: true}),
				},
			},
		},
	},
	{
		name: "column definition collision",
		setup: []testCommand{
			{commands.SqlCmd{}, []string{"-q", "alter table test add column c40 int;"}},
			{commands.SqlCmd{}, []string{"-q", "alter table test add column c6 bigint;"}},
			{commands.AddCmd{}, []string{"."}},
			{commands.CommitCmd{}, []string{"-m", "modified branch main"}},
			{commands.CheckoutCmd{}, []string{"other"}},
			{commands.SqlCmd{}, []string{"-q", "alter table test add column c40 int;"}},
			{commands.SqlCmd{}, []string{"-q", "alter table test rename column c40 to c44;"}},
			{commands.SqlCmd{}, []string{"-q", "alter table test add column c6 tinyint;"}},
			{commands.AddCmd{}, []string{"."}},
			{commands.CommitCmd{}, []string{"-m", "modified branch other"}},
			{commands.CheckoutCmd{}, []string{env.DefaultInitBranch}},
		},
		expConflict: merge.SchemaConflict{
			TableName: doltdb.TableName{Name: "test"},
			ColConflicts: []merge.ColConflict{
				{
					Kind:   merge.TagCollision,
					Ours:   newColTypeInfo("c40", uint64(679), typeinfo.Int32Type, false),
					Theirs: newColTypeInfo("c44", uint64(679), typeinfo.Int32Type, false),
				},
				{
					Kind:   merge.TagCollision,
					Ours:   newColTypeInfo("c6", uint64(10774), typeinfo.Int64Type, false),
					Theirs: newColTypeInfo("c6", uint64(10774), typeinfo.Int8Type, false),
				},
			},
		},
	},
	{
		name: "index definition collision",
		setup: []testCommand{
			{commands.SqlCmd{}, []string{"-q", "create index c3_idx on test(c3);"}},
			{commands.AddCmd{}, []string{"."}},
			{commands.CommitCmd{}, []string{"-m", "modified branch main"}},
			{commands.CheckoutCmd{}, []string{"other"}},
			{commands.SqlCmd{}, []string{"-q", "create index c3_index on test(c3);"}},
			{commands.AddCmd{}, []string{"."}},
			{commands.CommitCmd{}, []string{"-m", "modified branch other"}},
			{commands.CheckoutCmd{}, []string{env.DefaultInitBranch}},
		},
		expConflict: merge.SchemaConflict{
			TableName: doltdb.TableName{Name: "test"},
			IdxConflicts: []merge.IdxConflict{
				{
					Kind:   merge.TagCollision,
					Ours:   schema.NewIndex("c3_idx", []uint64{4696}, []uint64{4696, 3228}, nil, schema.IndexProperties{IsUserDefined: true}),
					Theirs: schema.NewIndex("c3_index", []uint64{4696}, []uint64{4696, 3228}, nil, schema.IndexProperties{IsUserDefined: true}),
				},
			},
		},
	},
	{
		name: "check definition collision",
		setup: []testCommand{
			{commands.SqlCmd{}, []string{"-q", "alter table test add constraint chk check (c3 > 0);"}},
			{commands.AddCmd{}, []string{"."}},
			{commands.CommitCmd{}, []string{"-m", "modified branch main"}},
			{commands.CheckoutCmd{}, []string{"other"}},
			{commands.SqlCmd{}, []string{"-q", "alter table test add constraint chk check (c3 < 0);"}},
			{commands.AddCmd{}, []string{"."}},
			{commands.CommitCmd{}, []string{"-m", "modified branch other"}},
			{commands.CheckoutCmd{}, []string{env.DefaultInitBranch}},
		},
		expConflict: merge.SchemaConflict{
			TableName: doltdb.TableName{Name: "test"},
			ChkConflicts: []merge.ChkConflict{
				{
					Kind:   merge.TagCollision,
					Ours:   schema.NewCheck("chk", "(`c3` > 0)", true),
					Theirs: schema.NewCheck("chk", "(`c3` < 0)", true),
				},
				{
					Kind:   merge.NameCollision,
					Ours:   schema.NewCheck("chk", "(`c3` > 0)", true),
					Theirs: schema.NewCheck("chk", "(`c3` < 0)", true),
				},
			},
		},
	},
	{
		name: "modified check",
		setup: []testCommand{
			{commands.SqlCmd{}, []string{"-q", "alter table test add constraint chk check (c3 > 0);"}},
			{commands.AddCmd{}, []string{"."}},
			{commands.CommitCmd{}, []string{"-m", "modified branch main"}},

			{commands.CheckoutCmd{}, []string{"other"}},
			{commands.SqlCmd{}, []string{"-q", "alter table test add constraint chk check (c3 > 0);"}},
			{commands.AddCmd{}, []string{"."}},
			{commands.CommitCmd{}, []string{"-m", "modified branch other"}},

			{commands.MergeCmd{}, []string{env.DefaultInitBranch}},
			{commands.CheckoutCmd{}, []string{env.DefaultInitBranch}},
			{commands.MergeCmd{}, []string{"other"}},

			{commands.SqlCmd{}, []string{"-q", "alter table test drop constraint chk;"}},
			{commands.SqlCmd{}, []string{"-q", "alter table test add constraint chk check (c3 > 10);"}},
			{commands.AddCmd{}, []string{"."}},
			{commands.CommitCmd{}, []string{"-m", "modified branch main"}},

			{commands.CheckoutCmd{}, []string{"other"}},
			{commands.SqlCmd{}, []string{"-q", "alter table test drop constraint chk;"}},
			{commands.SqlCmd{}, []string{"-q", "alter table test add constraint chk check (c3 < 10);"}},
			{commands.AddCmd{}, []string{"."}},
			{commands.CommitCmd{}, []string{"-m", "modified branch other"}},
		},
		expConflict: merge.SchemaConflict{
			TableName: doltdb.TableName{Name: "test"},
			ChkConflicts: []merge.ChkConflict{
				{
					Kind:   merge.TagCollision,
					Ours:   schema.NewCheck("chk", "(`c3` > 10)", true),
					Theirs: schema.NewCheck("chk", "(`c3` < 10)", true),
				},
				{
					Kind:   merge.NameCollision,
					Ours:   schema.NewCheck("chk", "(`c3` > 10)", true),
					Theirs: schema.NewCheck("chk", "(`c3` < 10)", true),
				},
			},
		},
	},
	{
		name: "primary key conflicts",
		setup: []testCommand{
			{commands.CheckoutCmd{}, []string{"other"}},
			{commands.SqlCmd{}, []string{"-q", "alter table test drop primary key;"}},
			{commands.AddCmd{}, []string{"."}},
			{commands.CommitCmd{}, []string{"-m", "modified branch other"}},
			{commands.CheckoutCmd{}, []string{env.DefaultInitBranch}},
		},
		expectedErr: merge.ErrMergeWithDifferentPks.New("test"),
	},
}

var setupForeignKeyTests = []testCommand{
	{commands.SqlCmd{}, []string{"-q", "create table test (" +
		"pk int not null primary key," +
		"t1 int not null," +
		"t2 int," +
		"t3 int);"}},
	{commands.SqlCmd{}, []string{"-q", "alter table test add index t1_idx (t1);"}},
	{commands.SqlCmd{}, []string{"-q", "create table quiz (" +
		"pk int not null primary key," +
		"q1 int not null," +
		"q2 int not null," +
		"index q2_idx (q2)," +
		"constraint q1_fk foreign key (q1) references test(t1));"}},
	{commands.AddCmd{}, []string{"."}},
	{commands.CommitCmd{}, []string{"-m", "setup common"}},
	{commands.BranchCmd{}, []string{"other"}},
}

var mergeForeignKeyTests = []mergeForeignKeyTest{
	{
		name:  "no changes",
		setup: []testCommand{},
		fkColl: fkCollection(doltdb.ForeignKey{
			Name:                   "q1_fk",
			TableName:              doltdb.TableName{Name: "quiz"},
			TableIndex:             "q1_fk",
			TableColumns:           []uint64{13001},
			ReferencedTableName:    doltdb.TableName{Name: "test"},
			ReferencedTableIndex:   "t1_idx",
			ReferencedTableColumns: []uint64{12111},
			UnresolvedFKDetails: doltdb.UnresolvedFKDetails{
				TableColumns:           []string{"q1"},
				ReferencedTableColumns: []string{"t1"},
			},
		}),
		expFKConflict: []merge.FKConflict{},
	},
	// {
	//	name: "add foreign key, drop foreign key, merge",
	//	setup: []testCommand{
	//		{commands.SqlCmd{}, []string{"-q", "alter table quiz add constraint q2_fk foreign key (q2) references test(t2);"}},
	//		{commands.AddCmd{}, []string{"."}},
	//		{commands.CommitCmd{}, []string{"-m", "modified branch main"}},
	//		{commands.CheckoutCmd{}, []string{"other"}},
	//		{commands.SqlCmd{}, []string{"-q", "alter table quiz drop constraint q1_fk;"}},
	//		{commands.AddCmd{}, []string{"."}},
	//		{commands.CommitCmd{}, []string{"-m", "modified branch other"}},
	//		{commands.CheckoutCmd{}, []string{"main"}},
	//	},
	//	fkColl: fkCollection(
	//		&doltdb.ForeignKey{
	//			Name:                   "q2_fk",
	//			TableName:              "quiz",
	//			TableIndex:             "dolt_fk_2",
	//			TableColumns:           []uint64{12},
	//			ReferencedTableName:    "test",
	//			ReferencedTableIndex:   "dolt_fk_2",
	//			ReferencedTableColumns: []uint64{2}}),
	//	expFKConflict: []merge.FKConflict{},
	// },
}

func colCollection(cols ...schema.Column) *schema.ColCollection {
	return schema.NewColCollection(cols...)
}

// SchemaFromColsAndIdxs creates a Schema from a ColCollection and an IndexCollection.
func schemaFromColsAndIdxs(allCols *schema.ColCollection, indexes ...schema.Index) schema.Schema {
	sch := schema.MustSchemaFromCols(allCols)
	sch.Indexes().AddIndex(indexes...)
	return sch
}

func newColTypeInfo(name string, tag uint64, typeInfo typeinfo.TypeInfo, partOfPK bool, constraints ...schema.ColConstraint) schema.Column {
	c, err := schema.NewColumnWithTypeInfo(name, tag, typeInfo, partOfPK, "", false, "", constraints...)
	if err != nil {
		panic("could not create column")
	}
	return c
}

func fkCollection(fks ...doltdb.ForeignKey) *doltdb.ForeignKeyCollection {
	fkc, err := doltdb.NewForeignKeyCollection(fks...)
	if err != nil {
		panic(err)
	}
	return fkc
}

func testMergeSchemas(t *testing.T, test mergeSchemaTest) {
	if test.skip {
		t.Skip()
		return
	}

	ctx := context.Background()
	dEnv := dtestutils.CreateTestEnv()
	defer dEnv.DoltDB(ctx).Close()

	cliCtx, _ := commands.NewArgFreeCliContext(ctx, dEnv, dEnv.FS)

	for _, c := range setupCommon {
		exit := c.exec(t, ctx, dEnv)
		require.Equal(t, 0, exit)
	}
	for _, c := range test.setup {
		exit := c.exec(t, ctx, dEnv)
		require.Equal(t, 0, exit)
	}

	// assert that we're on main
	exitCode := commands.CheckoutCmd{}.Exec(ctx, "checkout", []string{env.DefaultInitBranch}, dEnv, cliCtx)
	require.Equal(t, 0, exitCode)

	// merge branches
	exitCode = commands.MergeCmd{}.Exec(ctx, "merge", []string{"other"}, dEnv, cliCtx)
	assert.Equal(t, 0, exitCode)

	wr, err := dEnv.WorkingRoot(ctx)
	assert.NoError(t, err)
	tbl, ok, err := wr.GetTable(ctx, doltdb.TableName{Name: "test"})
	assert.True(t, ok)
	require.NoError(t, err)
	sch, err := tbl.GetSchema(ctx)
	require.NoError(t, err)

	assert.Equal(t, test.sch.GetAllCols(), sch.GetAllCols())
	assert.Equal(t, test.sch.Indexes(), sch.Indexes())
}

func testMergeSchemasWithConflicts(t *testing.T, test mergeSchemaConflictTest) {
	ctx := sql.NewEmptyContext()
	getSchema := func(t *testing.T, dEnv *env.DoltEnv) schema.Schema {
		wr, err := dEnv.WorkingRoot(ctx)
		assert.NoError(t, err)
		tbl, ok, err := wr.GetTable(ctx, doltdb.TableName{Name: "test"})
		assert.True(t, ok)
		require.NoError(t, err)
		sch, err := tbl.GetSchema(ctx)
		require.NoError(t, err)
		return sch
	}

	dEnv := dtestutils.CreateTestEnv()
	defer dEnv.DoltDB(ctx).Close()
	for _, c := range setupCommon {
		exit := c.exec(t, ctx, dEnv)
		require.Equal(t, 0, exit)
	}

	ancSch := getSchema(t, dEnv)

	for _, c := range test.setup {
		exit := c.exec(t, ctx, dEnv)
		require.Equal(t, 0, exit)
	}

	cliCtx, _ := commands.NewArgFreeCliContext(ctx, dEnv, dEnv.FS)

	// assert that we're on main
	exitCode := commands.CheckoutCmd{}.Exec(ctx, "checkout", []string{env.DefaultInitBranch}, dEnv, cliCtx)
	require.Equal(t, 0, exitCode)

	mainSch := getSchema(t, dEnv)

	exitCode = commands.CheckoutCmd{}.Exec(ctx, "checkout", []string{"other"}, dEnv, cliCtx)
	require.Equal(t, 0, exitCode)

	otherSch := getSchema(t, dEnv)

	_, actConflicts, mergeInfo, _, err := merge.SchemaMerge(ctx, types.Format_Default, mainSch, otherSch, ancSch, doltdb.TableName{Name: "test"})
	assert.False(t, mergeInfo.InvalidateSecondaryIndexes)
	if test.expectedErr != nil {
		// We don't use errors.Is here because errors generated by `Kind.New` compare stack traces in their `Is` implementation.
		assert.Equal(t, err.Error(), test.expectedErr.Error(), "Expected error '%s', instead got '%s'", test.expectedErr.Error(), err.Error())
		return
	}

	require.NoError(t, err)
	assert.Equal(t, actConflicts.TableName.Name, "test")

	assert.Equal(t, test.expConflict.Count(), actConflicts.Count())

	require.Equal(t, len(test.expConflict.IdxConflicts), len(actConflicts.IdxConflicts))
	for i, acc := range actConflicts.IdxConflicts {
		assert.True(t, test.expConflict.IdxConflicts[i].Ours.Equals(acc.Ours))
		assert.True(t, test.expConflict.IdxConflicts[i].Theirs.Equals(acc.Theirs))
	}

	require.Equal(t, len(test.expConflict.ColConflicts), len(actConflicts.ColConflicts))
	for i, icc := range actConflicts.ColConflicts {
		assert.True(t, test.expConflict.ColConflicts[i].Ours.Equals(icc.Ours))
		assert.True(t, test.expConflict.ColConflicts[i].Theirs.Equals(icc.Theirs))
	}

	require.Equal(t, len(test.expConflict.ChkConflicts), len(actConflicts.ChkConflicts))
	for i, icc := range actConflicts.ChkConflicts {
		assert.True(t, test.expConflict.ChkConflicts[i].Ours == icc.Ours)
		assert.True(t, test.expConflict.ChkConflicts[i].Theirs == icc.Theirs)
	}
}

func testMergeForeignKeys(t *testing.T, test mergeForeignKeyTest) {
	ctx := context.Background()
	dEnv := dtestutils.CreateTestEnv()
	defer dEnv.DoltDB(ctx).Close()
	for _, c := range setupForeignKeyTests {
		exit := c.exec(t, ctx, dEnv)
		require.Equal(t, 0, exit)
	}

	ancRoot, err := dEnv.WorkingRoot(ctx)
	require.NoError(t, err)

	for _, c := range test.setup {
		exit := c.exec(t, ctx, dEnv)
		require.Equal(t, 0, exit)
	}

	cliCtx, _ := commands.NewArgFreeCliContext(ctx, dEnv, dEnv.FS)

	// assert that we're on main
	exitCode := commands.CheckoutCmd{}.Exec(ctx, "checkout", []string{env.DefaultInitBranch}, dEnv, cliCtx)
	require.Equal(t, 0, exitCode)

	mainWS, err := dEnv.WorkingSet(ctx)
	require.NoError(t, err)
	mainRoot := mainWS.WorkingRoot()

	exitCode = commands.CheckoutCmd{}.Exec(ctx, "checkout", []string{"other"}, dEnv, cliCtx)
	require.Equal(t, 0, exitCode)

	otherWS, err := dEnv.WorkingSet(ctx)
	require.NoError(t, err)
	otherRoot := otherWS.WorkingRoot()

	opts := editor.TestEditorOptions(dEnv.DoltDB(ctx).ValueReadWriter())
	mo := merge.MergeOpts{IsCherryPick: false}
	result, err := merge.MergeRoots(sql.NewContext(ctx), mainRoot, otherRoot, ancRoot, mainWS, otherWS, opts, mo)
	assert.NoError(t, err)

	fkc, err := result.Root.GetForeignKeyCollection(ctx)
	assert.NoError(t, err)
	assert.Equal(t, test.fkColl.Count(), fkc.Count())

	err = test.fkColl.Iter(func(expFK doltdb.ForeignKey) (stop bool, err error) {
		actFK, ok := fkc.GetByTags(expFK.TableColumns, expFK.ReferencedTableColumns)
		assert.True(t, ok)
		assert.Equal(t, expFK, actFK)
		return false, nil
	})
	assert.NoError(t, err)
}
