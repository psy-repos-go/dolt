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

package enginetest

import (
	"fmt"
	"strings"

	"github.com/dolthub/go-mysql-server/enginetest/queries"
	"github.com/dolthub/go-mysql-server/sql"

	"github.com/dolthub/dolt/go/libraries/doltcore/schema"
)

// fillerVarchar pushes the tree into level 3
var fillerVarchar = strings.Repeat("x", 500)

var DoltHistogramTests = []queries.ScriptTest{
	//{
	//	Name: "mcv checking",
	//	SetUpScript: []string{
	//		"CREATE table xy (x bigint primary key, y int, z varchar(500), key(y,z));",
	//		"insert into xy values (0,0,'a'), (1,0,'a'), (2,0,'a'), (3,0,'a'), (4,1,'a'), (5,2,'a')",
	//		"analyze table xy",
	//	},
	//	Assertions: []queries.ScriptTestAssertion{
	//		{
	//			Query: " SELECT mcv_cnt from information_schema.column_statistics join json_table(histogram, '$.statistic.buckets[*]' COLUMNS(mcv_cnt JSON path '$.mcv_counts')) as dt  where table_name = 'xy' and column_name = 'y,z'",
	//			Expected: []sql.Row{
	//				{types.JSONDocument{Val: []interface{}{
	//					float64(4),
	//				}}},
	//			},
	//		},
	//		{
	//			Query: " SELECT mcv from information_schema.column_statistics join json_table(histogram, '$.statistic.buckets[*]' COLUMNS(mcv JSON path '$.mcvs[*]')) as dt  where table_name = 'xy' and column_name = 'y,z'",
	//			Expected: []sql.Row{
	//				{types.JSONDocument{Val: []interface{}{
	//					[]interface{}{float64(0), "a"},
	//				}}},
	//			},
	//		},
	//		{
	//			Query: " SELECT x,z from information_schema.column_statistics join json_table(histogram, '$.statistic.buckets[*]' COLUMNS(x bigint path '$.upper_bound[0]', z text path '$.upper_bound[1]')) as dt  where table_name = 'xy' and column_name = 'y,z'",
	//			Expected: []sql.Row{
	//				{2, "a"},
	//			},
	//		},
	//	},
	//},
	//{
	//	Name: "int pk",
	//	SetUpScript: []string{
	//		"CREATE table xy (x bigint primary key, y varchar(500));",
	//		fmt.Sprintf("insert into xy select x, '%s' from (with recursive inputs(x) as (select 1 union select x+1 from inputs where x < 10000) select * from inputs) dt", fillerVarchar),
	//		fmt.Sprintf("insert into xy select x, '%s'  from (with recursive inputs(x) as (select 10001 union select x+1 from inputs where x < 20000) select * from inputs) dt", fillerVarchar),
	//		fmt.Sprintf("insert into xy select x, '%s'  from (with recursive inputs(x) as (select 20001 union select x+1 from inputs where x < 30000) select * from inputs) dt", fillerVarchar),
	//		"analyze table xy",
	//	},
	//	Assertions: []queries.ScriptTestAssertion{
	//		{
	//			Query:    "SELECT json_length(json_extract(histogram, \"$.statistic.buckets\")) from information_schema.column_statistics where column_name = 'x'",
	//			Expected: []sql.Row{{32}},
	//		},
	//		{
	//			Query:    " SELECT sum(cnt) from information_schema.column_statistics join json_table(histogram, '$.statistic.buckets[*]' COLUMNS(cnt int path '$.row_count')) as dt  where table_name = 'xy' and column_name = 'x'",
	//			Expected: []sql.Row{{float64(30000)}},
	//		},
	//		{
	//			Query:    " SELECT sum(cnt) from information_schema.column_statistics join json_table(histogram, '$.statistic.buckets[*]' COLUMNS(cnt int path '$.null_count')) as dt  where table_name = 'xy' and column_name = 'x'",
	//			Expected: []sql.Row{{float64(0)}},
	//		},
	//		{
	//			Query:    " SELECT sum(cnt) from information_schema.column_statistics join json_table(histogram, '$.statistic.buckets[*]' COLUMNS(cnt int path '$.distinct_count')) as dt  where table_name = 'xy' and column_name = 'x'",
	//			Expected: []sql.Row{{float64(30000)}},
	//		},
	//		{
	//			Query:    " SELECT max(bound_cnt) from information_schema.column_statistics join json_table(histogram, '$.statistic.buckets[*]' COLUMNS(bound_cnt int path '$.bound_count')) as dt  where table_name = 'xy' and column_name = 'x'",
	//			Expected: []sql.Row{{int64(1)}},
	//		},
	//	},
	//},
	//{
	//	Name: "nulls distinct across chunk boundary",
	//	SetUpScript: []string{
	//		"CREATE table xy (x bigint primary key, y varchar(500), z bigint, key(z));",
	//		fmt.Sprintf("insert into xy select x, '%s', x  from (with recursive inputs(x) as (select 1 union select x+1 from inputs where x < 200) select * from inputs) dt", fillerVarchar),
	//		fmt.Sprintf("insert into xy select x, '%s', NULL  from (with recursive inputs(x) as (select 201 union select x+1 from inputs where x < 400) select * from inputs) dt", fillerVarchar),
	//		"analyze table xy",
	//	},
	//	Assertions: []queries.ScriptTestAssertion{
	//		{
	//			Query: "call dolt_stats_wait()",
	//		},
	//		{
	//			Query:    "SELECT json_length(json_extract(histogram, \"$.statistic.buckets\")) from information_schema.column_statistics where column_name = 'z'",
	//			Expected: []sql.Row{{2}},
	//		},
	//		{
	//			// bucket boundary duplication
	//			Query:    "SELECT json_value(histogram, \"$.statistic.distinct_count\", 'signed') from information_schema.column_statistics where column_name = 'z'",
	//			Expected: []sql.Row{{202}},
	//		},
	//		{
	//			Query:    " SELECT sum(cnt) from information_schema.column_statistics join json_table(histogram, '$.statistic.buckets[*]' COLUMNS(cnt int path '$.row_count')) as dt  where table_name = 'xy' and column_name = 'z'",
	//			Expected: []sql.Row{{float64(400)}},
	//		},
	//		{
	//			Query:    " SELECT sum(cnt) from information_schema.column_statistics join json_table(histogram, '$.statistic.buckets[*]' COLUMNS(cnt int path '$.null_count')) as dt  where table_name = 'xy' and column_name = 'z'",
	//			Expected: []sql.Row{{float64(200)}},
	//		},
	//		{
	//			// chunk border double count
	//			Query:    " SELECT sum(cnt) from information_schema.column_statistics join json_table(histogram, '$.statistic.buckets[*]' COLUMNS(cnt int path '$.distinct_count')) as dt  where table_name = 'xy' and column_name = 'z'",
	//			Expected: []sql.Row{{float64(202)}},
	//		},
	//		{
	//			// max bound count is an all nulls chunk
	//			Query:    " SELECT max(bound_cnt) from information_schema.column_statistics join json_table(histogram, '$.statistic.buckets[*]' COLUMNS(bound_cnt int path '$.bound_count')) as dt  where table_name = 'xy' and column_name = 'z'",
	//			Expected: []sql.Row{{int64(183)}},
	//		},
	//	},
	//},
	//{
	//	Name: "int index",
	//	SetUpScript: []string{
	//		"CREATE table xy (x bigint primary key, y varchar(500), z bigint, key(z));",
	//		fmt.Sprintf("insert into xy select x, '%s', x from (with recursive inputs(x) as (select 1 union select x+1 from inputs where x < 10000) select * from inputs) dt", fillerVarchar),
	//		fmt.Sprintf("insert into xy select x, '%s', x  from (with recursive inputs(x) as (select 10001 union select x+1 from inputs where x < 20000) select * from inputs) dt", fillerVarchar),
	//		fmt.Sprintf("insert into xy select x, '%s', NULL  from (with recursive inputs(x) as (select 20001 union select x+1 from inputs where x < 30000) select * from inputs) dt", fillerVarchar),
	//		"analyze table xy",
	//	},
	//	Assertions: []queries.ScriptTestAssertion{
	//		{
	//			Query:    "SELECT json_length(json_extract(histogram, \"$.statistic.buckets\")) from information_schema.column_statistics where column_name = 'z'",
	//			Expected: []sql.Row{{152}},
	//		},
	//		{
	//			Query:    " SELECT sum(cnt) from information_schema.column_statistics join json_table(histogram, '$.statistic.buckets[*]' COLUMNS(cnt int path '$.row_count')) as dt  where table_name = 'xy' and column_name = 'z'",
	//			Expected: []sql.Row{{float64(30000)}},
	//		},
	//		{
	//			Query:    " SELECT sum(cnt) from information_schema.column_statistics join json_table(histogram, '$.statistic.buckets[*]' COLUMNS(cnt int path '$.null_count')) as dt  where table_name = 'xy' and column_name = 'z'",
	//			Expected: []sql.Row{{float64(10000)}},
	//		},
	//		{
	//			// border NULL double count
	//			Query:    " SELECT sum(cnt) from information_schema.column_statistics join json_table(histogram, '$.statistic.buckets[*]' COLUMNS(cnt int path '$.distinct_count')) as dt  where table_name = 'xy' and column_name = 'z'",
	//			Expected: []sql.Row{{float64(20036)}},
	//		},
	//		{
	//			// max bound count is nulls chunk
	//			Query:    " SELECT max(bound_cnt) from information_schema.column_statistics join json_table(histogram, '$.statistic.buckets[*]' COLUMNS(bound_cnt int path '$.bound_count')) as dt  where table_name = 'xy' and column_name = 'z'",
	//			Expected: []sql.Row{{int64(440)}},
	//		},
	//	},
	//},
	{
		Name: "multiint index",
		SetUpScript: []string{
			"CREATE table xy (x bigint primary key, y varchar(500), z bigint, key(x, z));",
			fmt.Sprintf("insert into xy select x, '%s', x+1  from (with recursive inputs(x) as (select 1 union select x+1 from inputs where x < 10000) select * from inputs) dt", fillerVarchar),
			fmt.Sprintf("insert into xy select x, '%s', x+1  from (with recursive inputs(x) as (select 10001 union select x+1 from inputs where x < 20000) select * from inputs) dt", fillerVarchar),
			fmt.Sprintf("insert into xy select x, '%s', NULL from (with recursive inputs(x) as (select 20001 union select x+1 from inputs where x < 30000) select * from inputs) dt", fillerVarchar),
		},
		Assertions: []queries.ScriptTestAssertion{
			{
				Query: "call dolt_stats_wait()",
			},
			{
				Query:    "SELECT json_length(json_extract(histogram, \"$.statistic.buckets\")) from information_schema.column_statistics where column_name = 'x,z'",
				Expected: []sql.Row{{155}},
			},
			{
				Query:    " SELECT sum(cnt) from information_schema.column_statistics join json_table(histogram, '$.statistic.buckets[*]' COLUMNS(cnt int path '$.row_count')) as dt  where table_name = 'xy' and column_name = 'x,z'",
				Expected: []sql.Row{{float64(30000)}},
			},
			{
				Query:    " SELECT sum(cnt) from information_schema.column_statistics join json_table(histogram, '$.statistic.buckets[*]' COLUMNS(cnt int path '$.null_count')) as dt  where table_name = 'xy' and column_name = 'x,z'",
				Expected: []sql.Row{{float64(10000)}},
			},
			{
				Query:    " SELECT sum(cnt) from information_schema.column_statistics join json_table(histogram, '$.statistic.buckets[*]' COLUMNS(cnt int path '$.distinct_count')) as dt  where table_name = 'xy' and column_name = 'x,z'",
				Expected: []sql.Row{{float64(30000)}},
			},
			{
				// max bound count is nulls chunk
				Query:    " SELECT max(bound_cnt) from information_schema.column_statistics join json_table(histogram, '$.statistic.buckets[*]' COLUMNS(bound_cnt int path '$.bound_count')) as dt  where table_name = 'xy' and column_name = 'x,z'",
				Expected: []sql.Row{{int64(1)}},
			},
		},
	},
	{
		Name: "multiint index small",
		SetUpScript: []string{
			"CREATE table xy (x bigint primary key, y varchar(500), z bigint, key(x, z));",
			fmt.Sprintf("insert into xy select x, '%s', x+1  from (with recursive inputs(x) as (select 1 union select x+1 from inputs where x < 2) select * from inputs) dt", fillerVarchar),
			fmt.Sprintf("insert into xy select x, '%s', x+1  from (with recursive inputs(x) as (select 3 union select x+1 from inputs where x < 4) select * from inputs) dt", fillerVarchar),
			fmt.Sprintf("insert into xy select x, '%s', NULL from (with recursive inputs(x) as (select 5 union select x+1 from inputs where x < 6) select * from inputs) dt", fillerVarchar),
		},
		Assertions: []queries.ScriptTestAssertion{
			{
				Query: "call dolt_stats_wait()",
			},
			{
				Query:    "SELECT json_length(json_extract(histogram, \"$.statistic.buckets\")) from information_schema.column_statistics where column_name = 'x,z'",
				Expected: []sql.Row{{1}},
			},
			{
				Query:    " SELECT sum(cnt) from information_schema.column_statistics join json_table(histogram, '$.statistic.buckets[*]' COLUMNS(cnt int path '$.row_count')) as dt  where table_name = 'xy' and column_name = 'x,z'",
				Expected: []sql.Row{{float64(6)}},
			},
			{
				Query:    " SELECT sum(cnt) from information_schema.column_statistics join json_table(histogram, '$.statistic.buckets[*]' COLUMNS(cnt int path '$.null_count')) as dt  where table_name = 'xy' and column_name = 'x,z'",
				Expected: []sql.Row{{float64(2)}},
			},
			{
				Query:    " SELECT sum(cnt) from information_schema.column_statistics join json_table(histogram, '$.statistic.buckets[*]' COLUMNS(cnt int path '$.distinct_count')) as dt  where table_name = 'xy' and column_name = 'x,z'",
				Expected: []sql.Row{{float64(6)}},
			},
			{
				// max bound count is nulls chunk
				Query:    " SELECT max(bound_cnt) from information_schema.column_statistics join json_table(histogram, '$.statistic.buckets[*]' COLUMNS(bound_cnt int path '$.bound_count')) as dt  where table_name = 'xy' and column_name = 'x,z'",
				Expected: []sql.Row{{int64(1)}},
			},
		},
	},
	{
		Name: "several int index",
		SetUpScript: []string{
			"CREATE table xy (x bigint primary key, y varchar(500), z bigint, key(z), key (x,z));",
			fmt.Sprintf("insert into xy select x, '%s', x+1  from (with recursive inputs(x) as (select 1 union select x+1 from inputs where x < 10000) select * from inputs) dt", fillerVarchar),
		},
		Assertions: []queries.ScriptTestAssertion{
			{
				Query: "call dolt_stats_purge()",
			},
			{
				Query:    "SELECT column_name from information_schema.column_statistics",
				Expected: []sql.Row{},
			},
			{
				Query: "analyze table xy",
			},
			{
				Query:    " SELECT column_name from information_schema.column_statistics",
				Expected: []sql.Row{{"x"}, {"z"}, {"x,z"}},
			},
		},
	},
	{
		Name: "varchar pk",
		SetUpScript: []string{
			"CREATE table xy (x varchar(16) primary key, y varchar(500));",
			fmt.Sprintf("insert into xy select cast (x as char), '%s'  from (with recursive inputs(x) as (select 1 union select x+1 from inputs where x < 10000) select * from inputs) dt", fillerVarchar),
			fmt.Sprintf("insert into xy select cast (x as char), '%s'  from (with recursive inputs(x) as (select 10001 union select x+1 from inputs where x < 20000) select * from inputs) dt", fillerVarchar),
			fmt.Sprintf("insert into xy select cast (x as char), '%s' from (with recursive inputs(x) as (select 20001 union select x+1 from inputs where x < 30000) select * from inputs) dt", fillerVarchar),
			"analyze table xy",
		},
		Assertions: []queries.ScriptTestAssertion{
			{
				Query:    "SELECT json_length(json_extract(histogram, \"$.statistic.buckets\")) from information_schema.column_statistics where column_name = 'x'",
				Expected: []sql.Row{{26}},
			},
			{
				Query:    " SELECT sum(cnt) from information_schema.column_statistics join json_table(histogram, '$.statistic.buckets[*]' COLUMNS(cnt int path '$.row_count')) as dt  where table_name = 'xy' and column_name = 'x'",
				Expected: []sql.Row{{float64(30000)}},
			},
			{
				Query:    " SELECT sum(cnt) from information_schema.column_statistics join json_table(histogram, '$.statistic.buckets[*]' COLUMNS(cnt int path '$.null_count')) as dt  where table_name = 'xy' and column_name = 'x'",
				Expected: []sql.Row{{float64(0)}},
			},
			{
				Query:    " SELECT sum(cnt) from information_schema.column_statistics join json_table(histogram, '$.statistic.buckets[*]' COLUMNS(cnt int path '$.distinct_count')) as dt  where table_name = 'xy' and column_name = 'x'",
				Expected: []sql.Row{{float64(30000)}},
			},
			{
				// max bound count is nulls chunk
				Query:    " SELECT max(bound_cnt) from information_schema.column_statistics join json_table(histogram, '$.statistic.buckets[*]' COLUMNS(bound_cnt int path '$.bound_count')) as dt  where table_name = 'xy' and column_name = 'x'",
				Expected: []sql.Row{{int64(1)}},
			},
		},
	},
	{
		Name: "int-varchar inverse ordinal pk",
		SetUpScript: []string{
			"CREATE table xy (x varchar(16), y varchar(500), z bigint, primary key(z,x));",
			fmt.Sprintf("insert into xy select cast (x as char), '%s', x  from (with recursive inputs(x) as (select 1 union select x+1 from inputs where x < 10000) select * from inputs) dt", fillerVarchar),
			fmt.Sprintf("insert into xy select cast (x as char), '%s', x  from (with recursive inputs(x) as (select 10001 union select x+1 from inputs where x < 20000) select * from inputs) dt", fillerVarchar),
			fmt.Sprintf("insert into xy select cast (x as char), '%s', x from (with recursive inputs(x) as (select 20001 union select x+1 from inputs where x < 30000) select * from inputs) dt", fillerVarchar),
			"analyze table xy",
		},
		Assertions: []queries.ScriptTestAssertion{
			{
				Query:    " SELECT column_name from information_schema.column_statistics",
				Expected: []sql.Row{{"z,x"}},
			},
			{
				Query:    "SELECT json_length(json_extract(histogram, \"$.statistic.buckets\")) from information_schema.column_statistics where column_name = 'z,x'",
				Expected: []sql.Row{{42}},
			},
			{
				Query:    " SELECT sum(cnt) from information_schema.column_statistics join json_table(histogram, '$.statistic.buckets[*]' COLUMNS(cnt int path '$.row_count')) as dt  where table_name = 'xy' and column_name = 'z,x'",
				Expected: []sql.Row{{float64(30000)}},
			},
			{
				Query:    " SELECT sum(cnt) from information_schema.column_statistics join json_table(histogram, '$.statistic.buckets[*]' COLUMNS(cnt int path '$.null_count')) as dt  where table_name = 'xy' and column_name = 'z,x'",
				Expected: []sql.Row{{float64(0)}},
			},
			{
				Query:    " SELECT sum(cnt) from information_schema.column_statistics join json_table(histogram, '$.statistic.buckets[*]' COLUMNS(cnt int path '$.distinct_count')) as dt  where table_name = 'xy' and column_name = 'z,x'",
				Expected: []sql.Row{{float64(30000)}},
			},
			{
				// max bound count is nulls chunk
				Query:    " SELECT max(bound_cnt) from information_schema.column_statistics join json_table(histogram, '$.statistic.buckets[*]' COLUMNS(bound_cnt int path '$.bound_count')) as dt  where table_name = 'xy' and column_name = 'z,x'",
				Expected: []sql.Row{{int64(1)}},
			},
		},
	},
}

var DoltStatsStorageTests = []queries.ScriptTest{
	{
		Name: "single-table",
		SetUpScript: []string{
			"CREATE table xy (x bigint primary key, y int, z varchar(500), key(y,z));",
			"insert into xy values (0,0,'a'), (1,0,'a'), (2,0,'a'), (3,0,'a'), (4,1,'a'), (5,2,'a')",
			"analyze table xy",
		},
		Assertions: []queries.ScriptTestAssertion{
			{
				Query: "select database_name, table_name, index_name, columns, types from dolt_statistics",
				Expected: []sql.Row{
					{"mydb", "xy", "primary", "x", "bigint"},
					{"mydb", "xy", "y", "y,z", "int,varchar(500)"},
				},
			},
			{
				Query:    fmt.Sprintf("select %s, %s, %s from dolt_statistics", schema.StatsRowCountColName, schema.StatsDistinctCountColName, schema.StatsNullCountColName),
				Expected: []sql.Row{{uint64(6), uint64(6), uint64(0)}, {uint64(6), uint64(3), uint64(0)}},
			},
			{
				Query: fmt.Sprintf("select %s, %s from dolt_statistics", schema.StatsUpperBoundColName, schema.StatsUpperBoundCntColName),
				Expected: []sql.Row{
					{"5", uint64(1)},
					{"2,a", uint64(1)},
				},
			},
			{
				Query: fmt.Sprintf("select %s, %s, %s, %s, %s from dolt_statistics", schema.StatsMcv1ColName, schema.StatsMcv2ColName, schema.StatsMcv3ColName, schema.StatsMcv4ColName, schema.StatsMcvCountsColName),
				Expected: []sql.Row{
					{"", "", "", "", ""},
					{"0,a", "", "", "", "4"},
				},
			},
		},
	},
	{
		Name: "issue 8964: alternative indexes panic",
		SetUpScript: []string{
			"create table geom_tbl(g geometry not null srid 0)",
			"insert into geom_tbl values (point(0,0)), (linestring(point(1,1), point(2,2)))",
			"alter table geom_tbl add spatial index (g)",
			"CREATE TABLE fullt_tbl (pk BIGINT UNSIGNED PRIMARY KEY, v1 VARCHAR(200), v2 VARCHAR(200), FULLTEXT idx (v1, v2));",
			"INSERT INTO fullt_tbl VALUES (1, 'abc', 'def pqr'), (2, 'ghi', 'jkl'), (3, 'mno', 'mno'), (4, 'stu vwx', 'xyz zyx yzx'), (5, 'ghs', 'mno shg');",
			"create table vector_tbl (id int primary key, v json);",
			`insert into vector_tbl values (1, '[4.0,3.0]'), (2, '[0.0,0.0]'), (3, '[-1.0,1.0]'), (4, '[0.0,-2.0]');`,
			`create vector index v_idx on vector_tbl(v);`,
			"create table gen_tbl (a  int primary key, b int as (a + 1) stored)",
			"insert into gen_tbl (a) values (0), (1), (2)",
			"create index i1 on gen_tbl(b)",
		},
		Assertions: []queries.ScriptTestAssertion{
			{
				Query: "analyze table geom_tbl, fullt_tbl, vector_tbl, gen_tbl",
			},
			{
				Query:    "select table_name, index_name from dolt_statistics",
				Expected: []sql.Row{{"fullt_tbl", "primary"}, {"gen_tbl", "primary"}, {"gen_tbl", "i1"}, {"vector_tbl", "primary"}},
			},
		},
	},
	{
		Name: "comma encoding bug",
		SetUpScript: []string{
			`create table a (a varbinary (32) primary key)`,
			"insert into a values ('hello, world')",
			"analyze table a",
		},
		Assertions: []queries.ScriptTestAssertion{
			{
				Query:    "select count(*) from dolt_statistics",
				Expected: []sql.Row{{1}},
			},
		},
	},
	{
		Name: "comma encoding mcv bug",
		SetUpScript: []string{
			`create table ab (a int primary key, b varbinary (32), t timestamp, index (b,t))`,
			"insert into ab values (1, 'no thank you, world', '2024-03-12 01:18:53'), (2, 'hi, world', '2024-03-12 01:18:53'), (3, 'hello, world', '2024-03-12 01:18:53'), (4, 'hello, world', '2024-03-12 01:18:53'),(5, 'hello, world', '2024-03-12 01:18:53'), (6, 'hello, world', '2024-03-12 01:18:53')",
			"analyze table ab",
		},
		Assertions: []queries.ScriptTestAssertion{
			{
				Query:    "select count(*) from dolt_statistics",
				Expected: []sql.Row{{2}},
			},
		},
	},
	{
		Name: "boundary nils don't panic when trying to convert to the zero type",
		SetUpScript: []string{
			"CREATE table xy (x bigint primary key, y varchar(10), key(y,x));",
			"insert into xy values (0,null),(1,null)",
			"analyze table xy",
		},
		Assertions: []queries.ScriptTestAssertion{
			{
				Query: "select database_name, table_name, index_name, columns, types from dolt_statistics",
				Expected: []sql.Row{
					{"mydb", "xy", "primary", "x", "bigint"},
					{"mydb", "xy", "y", "y,x", "varchar(10),bigint"},
				},
			},
		},
	},
	{
		Name: "binary types round-trip",
		SetUpScript: []string{
			"CREATE table xy (x bigint primary key, y varbinary(10), z binary(14), key(y(9)), key(z));",
			"insert into xy values (0,'row 1', 'row 1'),(1,'row 2', 'row 1')",
			"analyze table xy",
		},
		Assertions: []queries.ScriptTestAssertion{
			{
				Query: "select database_name, table_name, index_name, columns, types from dolt_statistics",
				Expected: []sql.Row{
					{"mydb", "xy", "y", "y", "varbinary(10)"},
					{"mydb", "xy", "primary", "x", "bigint"},
					{"mydb", "xy", "z", "z", "binary(14)"},
				},
			},
			{
				Query:    "select count(*) from dolt_statistics",
				Expected: []sql.Row{{3}},
			},
		},
	},
	{
		Name: "timestamp types round-trip",
		SetUpScript: []string{
			"CREATE table xy (x bigint primary key, y timestamp, key(y));",
			"insert into xy values (0,'2024-03-11 18:52:44'),(1,'2024-03-11 19:22:12')",
			"analyze table xy",
		},
		Assertions: []queries.ScriptTestAssertion{
			{
				Query: "select database_name, table_name, index_name, columns, types from dolt_statistics",
				Expected: []sql.Row{
					{"mydb", "xy", "primary", "x", "bigint"},
					{"mydb", "xy", "y", "y", "timestamp"},
				},
			},
			{
				Query:    "select count(*) from dolt_statistics",
				Expected: []sql.Row{{2}},
			},
		},
	},
	{
		Name: "multi-table",
		SetUpScript: []string{
			"CREATE table xy (x bigint primary key, y int, z varchar(500), key(y,z));",
			"insert into xy values (0,0,'a'), (1,0,'a'), (2,0,'a'), (3,0,'a'), (4,1,'a'), (5,2,'a')",
			"CREATE table ab (a bigint primary key, b int, c int, key(b,c));",
			"insert into ab values (0,0,1), (1,0,1), (2,0,1), (3,0,1), (4,1,1), (5,2,1)",
			"analyze table xy",
			"analyze table ab",
		},
		Assertions: []queries.ScriptTestAssertion{
			{
				Query: "select database_name, table_name, index_name, columns, types  from dolt_statistics where table_name = 'xy'",
				Expected: []sql.Row{
					{"mydb", "xy", "primary", "x", "bigint"},
					{"mydb", "xy", "y", "y,z", "int,varchar(500)"},
				},
			},
			{
				Query:    fmt.Sprintf("select %s, %s, %s from dolt_statistics where table_name = 'xy'", schema.StatsRowCountColName, schema.StatsDistinctCountColName, schema.StatsNullCountColName),
				Expected: []sql.Row{{uint64(6), uint64(6), uint64(0)}, {uint64(6), uint64(3), uint64(0)}},
			},
			{
				Query: "select `table_name`, `index_name` from dolt_statistics",
				Expected: []sql.Row{
					{"ab", "primary"},
					{"ab", "b"},
					{"xy", "primary"},
					{"xy", "y"},
				},
			},
			{
				Query: "select database_name, table_name, index_name, columns, types  from dolt_statistics where table_name = 'ab'",
				Expected: []sql.Row{
					{"mydb", "ab", "primary", "a", "bigint"},
					{"mydb", "ab", "b", "b,c", "int,int"},
				},
			},
			{
				Query:    fmt.Sprintf("select %s, %s, %s from dolt_statistics where table_name = 'ab'", schema.StatsRowCountColName, schema.StatsDistinctCountColName, schema.StatsNullCountColName),
				Expected: []sql.Row{{uint64(6), uint64(6), uint64(0)}, {uint64(6), uint64(3), uint64(0)}},
			},
		},
	},
	{
		// only edited chunks are scanned and re-written
		Name: "incremental stats updates",
		SetUpScript: []string{
			"CREATE table xy (x bigint primary key, y int, z varchar(500), key(y,z));",
			"insert into xy values (0,0,'a'), (2,0,'a'), (4,1,'a'), (6,2,'a')",
			"analyze table xy",
			"insert into xy values (1,0,'a'), (3,0,'a'), (5,2,'a'),  (7,1,'a')",
			"analyze table xy",
		},
		Assertions: []queries.ScriptTestAssertion{
			{
				Query: fmt.Sprintf("select %s, %s, %s from dolt_statistics where table_name = 'xy'", schema.StatsRowCountColName, schema.StatsDistinctCountColName, schema.StatsNullCountColName),
				Expected: []sql.Row{
					{uint64(8), uint64(8), uint64(0)},
					{uint64(8), uint64(3), uint64(0)},
				},
			},
		},
	},
	{
		Name: "incremental stats deletes manual analyze",
		SetUpScript: []string{
			"CREATE table xy (x bigint primary key, y int, z varchar(500), key(y,z));",
			"insert into xy select x, 1, 1 from (with recursive inputs(x) as (select 4 union select x+1 from inputs where x < 1000) select * from inputs) dt;",
			"analyze table xy",
		},
		Assertions: []queries.ScriptTestAssertion{
			{
				Query:    "select count(*) as cnt from dolt_statistics group by table_name, index_name order by cnt",
				Expected: []sql.Row{{6}, {7}},
			},
			{
				Query: "delete from xy where x > 500",
			},
			{
				Query: "analyze table xy",
			},
			{
				Query:    "select count(*) from dolt_statistics group by table_name, index_name",
				Expected: []sql.Row{{4}, {4}},
			},
		},
	},
	{
		Name: "incremental stats deletes auto",
		SetUpScript: []string{
			"CREATE table xy (x bigint primary key, y int, z varchar(500), key(y,z));",
			"insert into xy select x, 1, 1 from (with recursive inputs(x) as (select 4 union select x+1 from inputs where x < 1000) select * from inputs) dt;",
			"analyze table xy",
		},
		Assertions: []queries.ScriptTestAssertion{
			{
				Query:    "select count(*) as cnt from dolt_statistics group by table_name, index_name order by cnt",
				Expected: []sql.Row{{6}, {7}},
			},
			{
				Query: "delete from xy where x > 500",
			},
			{
				Query: "analyze table xy",
			},
			{
				Query:    "select count(*) from dolt_statistics group by table_name, index_name",
				Expected: []sql.Row{{4}, {4}},
			},
		},
	},
	{
		// https://github.com/dolthub/dolt/issues/8504
		Name: "alter index column type",
		SetUpScript: []string{
			"CREATE table xy (x bigint primary key, y varchar(16))",
			"insert into xy values (0,'0'), (1,'1'), (2,'2')",
			"analyze table xy",
		},
		Assertions: []queries.ScriptTestAssertion{
			{
				Query:    "select count(*) from dolt_statistics group by table_name, index_name",
				Expected: []sql.Row{{1}},
			},
			{
				Query: "alter table xy modify column x varchar(16);",
			},
			{
				Query: "insert into xy values ('3', '3')",
			},
			{
				Query: "call dolt_stats_restart()",
			},
			{
				Query: "select sleep(.2)",
			},
			{
				Query:    "select count(*) from dolt_statistics group by table_name, index_name",
				Expected: []sql.Row{{1}},
			},
		},
	},
	{
		Name: "drop primary key",
		SetUpScript: []string{
			"CREATE table xy (x bigint primary key, y varchar(16))",
			"insert into xy values (0,'0'), (1,'1'), (2,'2')",
			"analyze table xy",
		},
		Assertions: []queries.ScriptTestAssertion{
			{
				Query:    "select count(*) from dolt_statistics group by table_name, index_name",
				Expected: []sql.Row{{1}},
			},
			{
				Query: "alter table xy drop primary key",
			},
			{
				Query: "insert into xy values ('3', '3')",
			},
			{
				Query: "analyze table xy",
			},
			{
				Query:    "select count(*) from dolt_statistics group by table_name, index_name",
				Expected: []sql.Row{},
			},
		},
	},
}

var StatBranchTests = []queries.ScriptTest{
	{
		Name: "multi branch stats",
		SetUpScript: []string{
			"CREATE table xy (x bigint primary key, y int, z varchar(500), key(y,z));",
			"insert into xy values (0,0,'a'), (1,0,'a'), (2,0,'a'), (3,0,'a'), (4,1,'a'), (5,2,'a')",
			"call dolt_commit('-Am', 'xy')",
			"call dolt_checkout('-b','feat')",
			"CREATE table ab (a bigint primary key, b int, c int, key(b,c));",
			"insert into ab values (0,0,1), (1,0,1), (2,0,1), (3,0,1), (4,1,1), (5,2,1)",
			"call dolt_commit('-Am', 'ab')",
			"call dolt_checkout('main')",
		},
		Assertions: []queries.ScriptTestAssertion{
			{
				Query: "call dolt_stats_wait()",
			},
			{
				Query: "select table_name, index_name, row_count from dolt_statistics",
				Expected: []sql.Row{
					{"xy", "primary", uint64(6)},
					{"xy", "y", uint64(6)},
				},
			},
			{
				Query: "select table_name, index_name, row_count from dolt_statistics as of 'feat'",
				Expected: []sql.Row{
					{"ab", "primary", uint64(6)},
					{"ab", "b", uint64(6)},
					{"xy", "primary", uint64(6)},
					{"xy", "y", uint64(6)},
				},
			},
			{
				Query: "select table_name, index_name, row_count from dolt_statistics as of 'main'",
				Expected: []sql.Row{
					{"xy", "primary", uint64(6)},
					{"xy", "y", uint64(6)},
				},
			},
			{
				Query: "call dolt_checkout('feat')",
			},
			{
				Query: "insert into xy values ('6',3,'a')",
			},
			{
				Query: "call dolt_commit('-am', 'cm')",
			},
			{
				Query: "call dolt_stats_wait()",
			},
			{
				Query: "select table_name, index_name, row_count from dolt_statistics as of 'feat'",
				Expected: []sql.Row{
					{"ab", "primary", uint64(6)},
					{"ab", "b", uint64(6)},
					{"xy", "primary", uint64(7)},
					{"xy", "y", uint64(7)},
				},
			},
			{
				Query: "select table_name, index_name, row_count from dolt_statistics as of 'main'",
				Expected: []sql.Row{
					{"xy", "primary", uint64(6)},
					{"xy", "y", uint64(6)},
				},
			},
		},
	},
	{
		Name: "issue #7710: branch connection string errors",
		SetUpScript: []string{
			"CREATE table xy (x bigint primary key, y int, z varchar(500), key(y,z));",
			"insert into xy values (0,0,'a'), (1,0,'a'), (2,0,'a'), (3,0,'a'), (4,1,'a'), (5,2,'a')",
			"use `mydb/main`",
		},
		Assertions: []queries.ScriptTestAssertion{
			{
				Query: "analyze table xy",
				Expected: []sql.Row{
					{"xy", "analyze", "status", "OK"},
				},
			},
		},
	},
}
