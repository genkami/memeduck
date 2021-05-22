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

// Delete creates a new DeleteStmt with given table name.
func Delete(table string) *DeleteStmt {
	return &DeleteStmt{
		table: table,
	}
}

// Where appends given conditional expressions to the DELETE statement.
func (s *DeleteStmt) Where(conds ...WhereCond) *DeleteStmt {
	return &DeleteStmt{
		table: s.table,
		conds: append(s.conds, conds...),
	}
}

func (s *DeleteStmt) SQL() (string, error) {
	stmt, err := s.toAST()
	if err != nil {
		return "", err
	}
	return stmt.SQL(), nil
}

func (s *DeleteStmt) toAST() (*ast.Delete, error) {
	if len(s.conds) == 0 {
		return nil, errors.New("no WHERE clause specified")
	}
	cond, err := s.conds[0].ToAstWhere()
	if err != nil {
		return nil, err
	}
	return &ast.Delete{
		TableName: &ast.Ident{Name: s.table},
		Where:     cond,
	}, nil
}

// InsertStmt builds INSERT statements.
type InsertStmt struct {
	table  string
	cols   []string
	values interface{}
}

// Insert creates a new InsertStmt with given table name. and column names.
func Insert(table string, cols []string) *InsertStmt {
	return &InsertStmt{
		table: table,
		cols:  cols,
	}
}

// Values returns an InsertStmt with its values set to given ones.
// It replaces existing values.
func (s *InsertStmt) Values(values interface{}) *InsertStmt {
	return &InsertStmt{
		table:  s.table,
		cols:   s.cols,
		values: values,
	}
}

func (is *InsertStmt) SQL() (string, error) {
	stmt, err := is.toAST()
	if err != nil {
		return "", err
	}
	return stmt.SQL(), nil
}

func (s *InsertStmt) toAST() (*ast.Insert, error) {
	cols := make([]*ast.Ident, 0, len(s.cols))
	for _, name := range s.cols {
		cols = append(cols, &ast.Ident{Name: name})
	}
	if s.values == nil {
		return nil, errors.New("neither VALUES nor SELECT specified")
	}
	input := &ast.ValuesInput{}
	// TODO: support SELECT
	rowsV := reflect.ValueOf(s.values)
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
		TableName: &ast.Ident{Name: s.table},
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
