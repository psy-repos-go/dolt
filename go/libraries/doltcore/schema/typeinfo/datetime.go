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

package typeinfo

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/dolthub/go-mysql-server/sql"
	gmstypes "github.com/dolthub/go-mysql-server/sql/types"
	"github.com/dolthub/vitess/go/sqltypes"

	"github.com/dolthub/dolt/go/store/types"
)

const (
	datetimeTypeParam_SQL           = "sql"
	datetimeTypeParam_SQL_Date      = "date"
	datetimeTypeParam_SQL_Datetime  = "datetime"
	datetimeTypeParam_SQL_Timestamp = "timestamp"
	datetimeTypeParam_Precision     = "precision"
)

type datetimeType struct {
	sqlDatetimeType sql.DatetimeType
}

var _ TypeInfo = (*datetimeType)(nil)
var (
	DateType      = &datetimeType{gmstypes.Date}
	DatetimeType  = &datetimeType{gmstypes.DatetimeMaxPrecision}
	TimestampType = &datetimeType{gmstypes.TimestampMaxPrecision}
)

func CreateDatetimeTypeFromSqlType(typ sql.DatetimeType) *datetimeType {
	return &datetimeType{typ}
}

func CreateDatetimeTypeFromParams(params map[string]string) (TypeInfo, error) {
	if sqlType, ok := params[datetimeTypeParam_SQL]; ok {
		precision := 6
		if precisionParam, ok := params[datetimeTypeParam_Precision]; ok {
			var err error
			precision, err = strconv.Atoi(precisionParam)
			if err != nil {
				return nil, err
			}
		}
		switch sqlType {
		case datetimeTypeParam_SQL_Date:
			return DateType, nil
		case datetimeTypeParam_SQL_Datetime:
			gmsType, err := gmstypes.CreateDatetimeType(sqltypes.Datetime, precision)
			if err != nil {
				return nil, err
			}
			return CreateDatetimeTypeFromSqlType(gmsType), nil
		case datetimeTypeParam_SQL_Timestamp:
			gmsType, err := gmstypes.CreateDatetimeType(sqltypes.Timestamp, precision)
			if err != nil {
				return nil, err
			}
			return CreateDatetimeTypeFromSqlType(gmsType), nil
		default:
			return nil, fmt.Errorf(`create datetime type info has invalid param "%v"`, sqlType)
		}
	} else {
		return nil, fmt.Errorf(`create datetime type info is missing param "%v"`, datetimeTypeParam_SQL)
	}
}

// ConvertNomsValueToValue implements TypeInfo interface.
func (ti *datetimeType) ConvertNomsValueToValue(v types.Value) (interface{}, error) {
	if val, ok := v.(types.Timestamp); ok {
		if ti.Equals(DateType) {
			return time.Time(val).Truncate(24 * time.Hour).UTC(), nil
		}
		return time.Time(val).UTC(), nil
	}
	if _, ok := v.(types.Null); ok || v == nil {
		return nil, nil
	}
	return nil, fmt.Errorf(`"%v" cannot convert NomsKind "%v" to a value`, ti.String(), v.Kind())
}

// ReadFrom reads a go value from a noms types.CodecReader directly
func (ti *datetimeType) ReadFrom(_ *types.NomsBinFormat, reader types.CodecReader) (interface{}, error) {
	k := reader.ReadKind()
	switch k {
	case types.TimestampKind:
		t, err := reader.ReadTimestamp()

		if err != nil {
			return nil, err
		}

		if ti.Equals(DateType) {
			return t.Truncate(24 * time.Hour).UTC(), nil
		}
		return t.UTC(), nil
	case types.NullKind:
		return nil, nil
	}

	return nil, fmt.Errorf(`"%v" cannot convert NomsKind "%v" to a value`, ti.String(), k)
}

// ConvertValueToNomsValue implements TypeInfo interface.
func (ti *datetimeType) ConvertValueToNomsValue(ctx context.Context, vrw types.ValueReadWriter, v interface{}) (types.Value, error) {
	//TODO: handle the zero value as a special case that is valid for all ranges
	if v == nil {
		return types.NullValue, nil
	}
	timeVal, _, err := ti.sqlDatetimeType.Convert(ctx, v)
	if err != nil {
		return nil, err
	}
	val, ok := timeVal.(time.Time)
	if ok {
		return types.Timestamp(val), nil
	}
	return nil, fmt.Errorf(`"%v" has unexpectedly encountered a value of type "%T" from embedded type`, ti.String(), v)
}

// Equals implements TypeInfo interface.
func (ti *datetimeType) Equals(other TypeInfo) bool {
	if other == nil {
		return false
	}
	if ti2, ok := other.(*datetimeType); ok {
		return ti.sqlDatetimeType.Type() == ti2.sqlDatetimeType.Type()
	}
	return false
}

// FormatValue implements TypeInfo interface.
func (ti *datetimeType) FormatValue(v types.Value) (*string, error) {
	if _, ok := v.(types.Null); ok || v == nil {
		return nil, nil
	}
	timeVal, err := ti.ConvertNomsValueToValue(v)
	if err != nil {
		return nil, err
	}
	val, ok := timeVal.(time.Time)
	if !ok {
		return nil, fmt.Errorf(`"%v" has unexpectedly encountered a value of type "%T" from embedded type`, ti.String(), v)
	}
	if ti.sqlDatetimeType.Type() == sqltypes.Date {
		res := val.Format(sql.DateLayout)
		return &res, nil
	} else {
		res := val.Format(sql.TimestampDatetimeLayout)
		return &res, nil
	}
}

// GetTypeIdentifier implements TypeInfo interface.
func (ti *datetimeType) GetTypeIdentifier() Identifier {
	return DatetimeTypeIdentifier
}

// GetTypeParams implements TypeInfo interface.
func (ti *datetimeType) GetTypeParams() map[string]string {
	params := map[string]string{}
	switch ti.sqlDatetimeType.Type() {
	case sqltypes.Date:
		params[datetimeTypeParam_SQL] = datetimeTypeParam_SQL_Date
	case sqltypes.Datetime:
		params[datetimeTypeParam_SQL] = datetimeTypeParam_SQL_Datetime
		params[datetimeTypeParam_Precision] = strconv.Itoa(ti.sqlDatetimeType.Precision())
	case sqltypes.Timestamp:
		params[datetimeTypeParam_SQL] = datetimeTypeParam_SQL_Timestamp
		params[datetimeTypeParam_Precision] = strconv.Itoa(ti.sqlDatetimeType.Precision())
	default:
		panic(fmt.Errorf(`unknown datetime type info sql type "%v"`, ti.sqlDatetimeType.Type().String()))
	}
	return params
}

// IsValid implements TypeInfo interface.
func (ti *datetimeType) IsValid(v types.Value) bool {
	if val, ok := v.(types.Timestamp); ok {
		_, _, err := ti.sqlDatetimeType.Convert(context.Background(), time.Time(val))
		if err != nil {
			return false
		}
		return true
	}
	if _, ok := v.(types.Null); ok || v == nil {
		return true
	}
	return false
}

// NomsKind implements TypeInfo interface.
func (ti *datetimeType) NomsKind() types.NomsKind {
	return types.TimestampKind
}

// Promote implements TypeInfo interface.
func (ti *datetimeType) Promote() TypeInfo {
	return &datetimeType{ti.sqlDatetimeType.Promote().(sql.DatetimeType)}
}

// String implements TypeInfo interface.
func (ti *datetimeType) String() string {
	return fmt.Sprintf(`Datetime(SQL: "%v")`, ti.sqlDatetimeType.String())
}

// ToSqlType implements TypeInfo interface.
func (ti *datetimeType) ToSqlType() sql.Type {
	return ti.sqlDatetimeType
}

// datetimeTypeConverter is an internal function for GetTypeConverter that handles the specific type as the source TypeInfo.
func datetimeTypeConverter(ctx context.Context, src *datetimeType, destTi TypeInfo) (tc TypeConverter, needsConversion bool, err error) {
	switch dest := destTi.(type) {
	case *bitType:
		return wrapConvertValueToNomsValue(dest.ConvertValueToNomsValue)
	case *blobStringType:
		return wrapConvertValueToNomsValue(dest.ConvertValueToNomsValue)
	case *boolType:
		return wrapConvertValueToNomsValue(dest.ConvertValueToNomsValue)
	case *datetimeType:
		return wrapIsValid(dest.IsValid, src, dest)
	case *decimalType:
		return wrapConvertValueToNomsValue(dest.ConvertValueToNomsValue)
	case *enumType:
		return wrapConvertValueToNomsValue(dest.ConvertValueToNomsValue)
	case *floatType:
		return wrapConvertValueToNomsValue(dest.ConvertValueToNomsValue)
	case *geomcollType:
		return wrapConvertValueToNomsValue(dest.ConvertValueToNomsValue)
	case *geometryType:
		return wrapConvertValueToNomsValue(dest.ConvertValueToNomsValue)
	case *inlineBlobType:
		return wrapConvertValueToNomsValue(dest.ConvertValueToNomsValue)
	case *intType:
		return wrapConvertValueToNomsValue(dest.ConvertValueToNomsValue)
	case *jsonType:
		return wrapConvertValueToNomsValue(dest.ConvertValueToNomsValue)
	case *linestringType:
		return wrapConvertValueToNomsValue(dest.ConvertValueToNomsValue)
	case *multilinestringType:
		return wrapConvertValueToNomsValue(dest.ConvertValueToNomsValue)
	case *multipointType:
		return wrapConvertValueToNomsValue(dest.ConvertValueToNomsValue)
	case *multipolygonType:
		return wrapConvertValueToNomsValue(dest.ConvertValueToNomsValue)
	case *pointType:
		return wrapConvertValueToNomsValue(dest.ConvertValueToNomsValue)
	case *polygonType:
		return wrapConvertValueToNomsValue(dest.ConvertValueToNomsValue)
	case *setType:
		return wrapConvertValueToNomsValue(dest.ConvertValueToNomsValue)
	case *timeType:
		return wrapConvertValueToNomsValue(dest.ConvertValueToNomsValue)
	case *uintType:
		return wrapConvertValueToNomsValue(dest.ConvertValueToNomsValue)
	case *uuidType:
		return wrapConvertValueToNomsValue(dest.ConvertValueToNomsValue)
	case *varBinaryType:
		return wrapConvertValueToNomsValue(dest.ConvertValueToNomsValue)
	case *varStringType:
		return wrapConvertValueToNomsValue(dest.ConvertValueToNomsValue)
	case *yearType:
		return wrapConvertValueToNomsValue(dest.ConvertValueToNomsValue)
	default:
		return nil, false, UnhandledTypeConversion.New(src.String(), destTi.String())
	}
}
