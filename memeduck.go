// Package memeduck provides tools to build Spanner SQL queries.
package memeduck

import (
	"reflect"
	"strconv"

	"github.com/MakeNowJust/memefish/pkg/ast"
	"github.com/pkg/errors"
)

// InsertIntoBuilder builds INSERT statements.
type InsertIntoBuilder struct {
	table string
	cols  []string
}

// Insert creates a new InsertBuilder with given table name and column names.
func InsertInto(table string, cols []string) *InsertIntoBuilder {
	return &InsertIntoBuilder{
		table: table,
		cols:  cols,
	}
}

// Values adds a VALUES clause to the insert statement.
func (ib *InsertIntoBuilder) Values(rows interface{}) *InsertIntoValuesBuilder {
	return &InsertIntoValuesBuilder{
		ib:   ib,
		rows: rows,
	}
}

// InsertIntoValuesBuilder builds INSERT statements with VALUES clauses.
type InsertIntoValuesBuilder struct {
	ib   *InsertIntoBuilder
	rows interface{}
}

func (ivb *InsertIntoValuesBuilder) SQL() (string, error) {
	stmt, err := ivb.toStmt()
	if err != nil {
		return "", err
	}
	return stmt.SQL(), nil
}

func (ivb *InsertIntoValuesBuilder) toStmt() (*ast.Insert, error) {
	cols := make([]*ast.Ident, 0, len(ivb.ib.cols))
	for _, name := range ivb.ib.cols {
		cols = append(cols, &ast.Ident{Name: name})
	}
	input := &ast.ValuesInput{}
	// TODO: check types
	rowsV := reflect.ValueOf(ivb.rows)
	for i := 0; i < rowsV.Len(); i++ {
		rowI := rowsV.Index(i).Interface()
		row, err := toValuesRow(rowI)
		if err != nil {
			return nil, errors.WithMessagef(err, "can't convert %T into SQL row", rowI)
		}
		input.Rows = append(input.Rows, row)
	}
	return &ast.Insert{
		TableName: &ast.Ident{Name: ivb.ib.table},
		Columns:   cols,
		Input:     input,
	}, nil
}

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
	case int:
		return toIntLit(int64(v)), nil
	default:
		return nil, errors.Errorf("can't convert %T into SQL expr", val)
	}
}

func toIntLit(v int64) *ast.IntLiteral {
	return &ast.IntLiteral{
		Base:  10,
		Value: strconv.FormatInt(v, 10),
	}
}
