// Package memeduck provides tools to build Spanner SQL queries.
package memeduck

import (
	"reflect"
	"strconv"

	"cloud.google.com/go/spanner"
	"github.com/MakeNowJust/memefish/pkg/ast"
	"github.com/pkg/errors"
)

// InsertStmt builds INSERT statements.
type InsertStmt struct {
	table string
	cols  []string
	input interface{}
}

// Insert creates a new InsertBuilder with given table name and column names.
func Insert(table string, cols []string, input interface{}) *InsertStmt {
	return &InsertStmt{
		table: table,
		cols:  cols,
		input: input,
	}
}

func (is *InsertStmt) SQL() (string, error) {
	stmt, err := is.toAST()
	if err != nil {
		return "", err
	}
	return stmt.SQL(), nil
}

func (is *InsertStmt) toAST() (*ast.Insert, error) {
	cols := make([]*ast.Ident, 0, len(is.cols))
	for _, name := range is.cols {
		cols = append(cols, &ast.Ident{Name: name})
	}
	input := &ast.ValuesInput{}
	// TODO: check types
	// TODO: support SELECT
	rowsV := reflect.ValueOf(is.input)
	for i := 0; i < rowsV.Len(); i++ {
		rowI := rowsV.Index(i).Interface()
		row, err := toValuesRow(rowI)
		if err != nil {
			return nil, errors.WithMessagef(err, "can't convert %T into SQL row", rowI)
		}
		input.Rows = append(input.Rows, row)
	}
	return &ast.Insert{
		TableName: &ast.Ident{Name: is.table},
		Columns:   cols,
		Input:     input,
	}, nil
}

// TODO:
// []string, []*string, []NullString - STRING ARRAY
// [][]byte - BYTES ARRAY
// []int, []int64, []*int64, []NullInt64 - INT64 ARRAY
// []bool, []*bool, []NullBool - BOOL ARRAY
// float64, *float64, NullFloat64 - FLOAT64
// []float64, []*float64, []NullFloat64 - FLOAT64 ARRAY
// time.Time, *time.Time, NullTime - TIMESTAMP
// []time.Time, []*time.Time, []NullTime - TIMESTAMP ARRAY
// Date, *Date, NullDate - DATE
// []Date, []*Date, []NullDate - DATE ARRAY
// big.Rat, *big.Rat, NullNumeric - NUMERIC
// []big.Rat, []*big.Rat, []NullNumeric - NUMERIC ARRAY
func toValuesRow(val interface{}) (*ast.ValuesRow, error) {
	row := &ast.ValuesRow{}
	valV := reflect.ValueOf(val)
	// TODO: check types
	for i := 0; i < valV.Len(); i++ {
		expr, err := toExpr(valV.Index(i).Interface())
		if err != nil {
			return nil, err
		}
		row.Exprs = append(row.Exprs, &ast.DefaultExpr{Expr: expr})
	}
	return row, nil
}

func toExpr(val interface{}) (ast.Expr, error) {
	switch v := val.(type) {
	case string:
		return stringLit(v), nil
	case *string:
		if v == nil {
			return nullLit(), nil
		}
		return stringLit(*v), nil
	case spanner.NullString:
		if !v.Valid {
			return nullLit(), nil
		}
		return stringLit(v.StringVal), nil
	case []byte:
		if v == nil {
			return nullLit(), nil
		}
		return bytesLit(v), nil
	case int:
		return intLit(int64(v)), nil
	case *int:
		if v == nil {
			return nullLit(), nil
		}
		return intLit(int64(*v)), nil
	case int64:
		return intLit(v), nil
	case *int64:
		if v == nil {
			return nullLit(), nil
		}
		return intLit(*v), nil
	case spanner.NullInt64:
		if !v.Valid {
			return nullLit(), nil
		}
		return intLit(v.Int64), nil
	case bool:
		return boolLit(v), nil
	case *bool:
		if v == nil {
			return nullLit(), nil
		}
		return boolLit(*v), nil
	case spanner.NullBool:
		if !v.Valid {
			return nullLit(), nil
		}
		return boolLit(v.Bool), nil
	case float64:
		return floatLit(v), nil
	case *float64:
		if v == nil {
			return nullLit(), nil
		}
		return floatLit(*v), nil
	case spanner.NullFloat64:
		if !v.Valid {
			return nullLit(), nil
		}
		return floatLit(v.Float64), nil
	default:
		return nil, errors.Errorf("can't convert %T into SQL expr", val)
	}
}

func stringLit(v string) *ast.StringLiteral {
	return &ast.StringLiteral{
		Value: v,
	}
}

func bytesLit(v []byte) *ast.BytesLiteral {
	return &ast.BytesLiteral{
		Value: v,
	}
}

func intLit(v int64) *ast.IntLiteral {
	return &ast.IntLiteral{
		Base:  10,
		Value: strconv.FormatInt(v, 10),
	}
}

func boolLit(v bool) *ast.BoolLiteral {
	return &ast.BoolLiteral{
		Value: v,
	}
}

func floatLit(v float64) *ast.FloatLiteral {
	return &ast.FloatLiteral{
		Value: strconv.FormatFloat(v, 'e', -1, 64),
	}
}

func nullLit() *ast.NullLiteral {
	return &ast.NullLiteral{}
}
