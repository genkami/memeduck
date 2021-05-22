// Package memeduck provides tools to build Spanner SQL queries.
package memeduck

import (
	"reflect"

	"github.com/MakeNowJust/memefish/pkg/ast"
	"github.com/pkg/errors"

	"github.com/genkami/memeduck/internal"
)

// WhereCond is a conditional expression that appears in WHERE clauses.
type WhereCond interface {
	ToAstWhere() (*ast.Where, error)
}

// DeleteStmt build DELETE statements.
type DeleteStmt struct {
	table string
	conds []WhereCond
}

// Delete creates a new DeleteStmt with given table name and where clause.
func Delete(table string, conds ...WhereCond) *DeleteStmt {
	return &DeleteStmt{
		table: table,
		conds: conds,
	}
}

func (ds *DeleteStmt) SQL() (string, error) {
	stmt, err := ds.toAST()
	if err != nil {
		return "", err
	}
	return stmt.SQL(), nil
}

func (ds *DeleteStmt) toAST() (*ast.Delete, error) {
	cond, err := ds.conds[0].ToAstWhere()
	if err != nil {
		return nil, err
	}
	return &ast.Delete{
		TableName: &ast.Ident{Name: ds.table},
		Where:     cond,
	}, nil
}

// InsertStmt builds INSERT statements.
type InsertStmt struct {
	table string
	cols  []string
	input interface{}
}

// Insert creates a new InsertStmt with given table name and column names.
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
	// TODO: support SELECT
	rowsV := reflect.ValueOf(is.input)
	if rowsV.Type().Kind() != reflect.Slice {
		return nil, errors.New("values it not a slice")
	}
	if rowsV.Len() <= 0 {
		return nil, errors.New("empty values")
	}
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

func toValuesRow(val interface{}) (*ast.ValuesRow, error) {
	row := &ast.ValuesRow{}
	valV := reflect.ValueOf(val)
	// TODO: check types
	for i := 0; i < valV.Len(); i++ {
		expr, err := internal.ToExpr(valV.Index(i).Interface())
		if err != nil {
			return nil, err
		}
		row.Exprs = append(row.Exprs, &ast.DefaultExpr{Expr: expr})
	}
	return row, nil
}
