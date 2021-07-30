package memeduck

import (
	"reflect"
	"strings"

	"github.com/MakeNowJust/memefish/pkg/ast"
	"github.com/pkg/errors"

	"github.com/genkami/memeduck/internal"
)

// SelectStmt builds SELECT statements.
type SelectStmt struct {
	table      string
	cols       []string
	conds      []WhereCond
	ords       []*ordering
	limit      *int
	offset     *int
	asStruct   bool
	subQueries []SubQuery
}

type ordering struct {
	col string
	dir Direction
}

func (o *ordering) toASTOrderByItem() *ast.OrderByItem {
	return &ast.OrderByItem{
		Expr: &ast.Ident{Name: o.col},
		Dir:  ast.Direction(o.dir),
	}
}

// Direction is an ordering direction used by ORDER BY clause.
type Direction ast.Direction

const (
	ASC  Direction = Direction(ast.DirectionAsc)
	DESC Direction = Direction(ast.DirectionDesc)
)

// Select creates a new SelectStmt with given table name and column names.
func Select(table string, cols []string) *SelectStmt {
	return &SelectStmt{
		table: table,
		cols:  cols,
	}
}

func (s *SelectStmt) AsStruct() *SelectStmt {
	var t = *s
	t.asStruct = true
	return &t
}

func (s *SelectStmt) SubQuery(queries ...SubQuery) *SelectStmt {
	var t = *s
	t.subQueries = append(t.subQueries, queries...)
	return &t
}

// Where appends given codintional expressions to the SELECT statement.
func (s *SelectStmt) Where(conds ...WhereCond) *SelectStmt {
	var t = *s
	t.conds = append(t.conds, conds...)
	return &t
}

// OrderBy appends a column to its ORDER BY clause.
func (s *SelectStmt) OrderBy(col string, dir Direction) *SelectStmt {
	var t = *s
	t.ords = append(t.ords, &ordering{
		col: col,
		dir: dir,
	})
	return &t
}

// Limit adds a LIMIT clause to the SELECT statement.
// It replaces existing LIMIT clauses.
func (s *SelectStmt) Limit(limit int) *SelectStmt {
	var t = *s
	t.limit = &limit
	return &t
}

// LimitOffset adds a LIMIT ... OFFSET ... clause to the SELECT statement.
// It replaces existing LIMIT clauses.
func (s *SelectStmt) LimitOffset(limit, offset int) *SelectStmt {
	var t = *s
	t.limit = &limit
	t.offset = &offset
	return &t
}

func (s *SelectStmt) SQL() (string, error) {
	stmt, err := s.toAST()
	if err != nil {
		return "", err
	}
	return stmt.SQL(), nil
}

func (s *SelectStmt) toAST() (*ast.Select, error) {
	var err error
	var where *ast.Where = nil
	if len(s.conds) > 0 {
		where, err = And(s.conds...).ToASTWhere()
		if err != nil {
			return nil, err
		}
	}

	if len(s.cols) <= 0 {
		return nil, errors.New("no columns specified")
	}
	items := make([]ast.SelectItem, 0, len(s.cols))
	for _, col := range s.cols {
		items = append(items, &ast.ExprSelectItem{
			Expr: &ast.Ident{Name: col},
		})
	}
	if len(s.subQueries) > 0 {
		for _, q := range s.subQueries {
			item, err := q.ToAST()
			if err != nil {
				return nil, err
			}
			items = append(items, item)
		}
	}

	var orderBy *ast.OrderBy = nil
	if len(s.ords) > 0 {
		items := make([]*ast.OrderByItem, 0, len(s.ords))
		for _, o := range s.ords {
			items = append(items, o.toASTOrderByItem())
		}
		orderBy = &ast.OrderBy{
			Items: items,
		}
	}

	var limit *ast.Limit = nil
	if s.limit != nil {
		limit = &ast.Limit{
			Count: internal.IntLit(int64(*s.limit)),
		}
		if s.offset != nil {
			limit.Offset = &ast.Offset{
				Value: internal.IntLit(int64(*s.offset)),
			}
		}
	}

	return &ast.Select{
		From: &ast.From{
			Source: &ast.TableName{
				Table: &ast.Ident{Name: s.table},
			},
		},
		AsStruct: s.asStruct,
		Results:  items,
		Where:    where,
		OrderBy:  orderBy,
		Limit:    limit,
	}, nil
}

// UpdateStmt builds UPDATE statements.
type UpdateStmt struct {
	table string
	items []*updateItem
	conds []WhereCond
}

type updateItem struct {
	ident *IdentExpr
	value interface{}
}

func (i *updateItem) toASTUpdateItem() (*ast.UpdateItem, error) {
	if len(i.ident.names) <= 0 {
		return nil, errors.New("empty identifier")
	}
	// NOTE: can't use ast.Path here for any reason.
	path := make([]*ast.Ident, 0, len(i.ident.names))
	for _, name := range i.ident.names {
		path = append(path, &ast.Ident{Name: name})
	}
	expr, err := internal.ToExpr(i.value)
	if err != nil {
		return nil, err
	}
	return &ast.UpdateItem{
		Path: path,
		Expr: expr,
	}, nil
}

// Update creates a new UpdateStmt with given table name.
func Update(table string) *UpdateStmt {
	return &UpdateStmt{
		table: table,
	}
}

// Set adds a assignment clause to the UPDATE statement.
func (s *UpdateStmt) Set(id *IdentExpr, value interface{}) *UpdateStmt {
	var t = *s
	t.items = append(t.items, &updateItem{
		ident: id,
		value: value,
	})
	return &t
}

// Where adds a WHERE clause to the UPDATE statement.
func (s *UpdateStmt) Where(conds ...WhereCond) *UpdateStmt {
	var t = *s
	t.conds = append(t.conds, conds...)
	return &t
}

func (s *UpdateStmt) SQL() (string, error) {
	stmt, err := s.toAST()
	if err != nil {
		return "", err
	}
	return stmt.SQL(), nil
}

func (s *UpdateStmt) toAST() (*ast.Update, error) {
	if len(s.items) <= 0 {
		return nil, errors.New("no SET clause is specified")
	}
	items := make([]*ast.UpdateItem, 0, len(s.items))
	for _, item := range s.items {
		astItem, err := item.toASTUpdateItem()
		if err != nil {
			return nil, err
		}
		items = append(items, astItem)
	}

	cond, err := And(s.conds...).ToASTWhere()
	if err != nil {
		return nil, err
	}
	return &ast.Update{
		TableName: &ast.Ident{Name: s.table},
		Updates:   items,
		Where:     cond,
	}, nil
}

// DeleteStmt builds DELETE statements.
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
	var t = *s
	t.conds = append(t.conds, conds...)
	return &t
}

func (s *DeleteStmt) SQL() (string, error) {
	stmt, err := s.toAST()
	if err != nil {
		return "", err
	}
	return stmt.SQL(), nil
}

func (s *DeleteStmt) toAST() (*ast.Delete, error) {
	cond, err := And(s.conds...).ToASTWhere()
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
	// TODO: support SELECT
	var input ast.InsertInput
	var err error
	rowsV := reflect.ValueOf(s.values)
	if rowsV.Type().Kind() == reflect.Slice {
		input, err = s.sliceToInsertInput(rowsV)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.Errorf("can't create InsertInput")
	}
	return &ast.Insert{
		TableName: &ast.Ident{Name: s.table},
		Columns:   cols,
		Input:     input,
	}, nil
}

func (s *InsertStmt) sliceToInsertInput(rowsV reflect.Value) (ast.InsertInput, error) {
	input := &ast.ValuesInput{}
	if rowsV.Len() <= 0 {
		return nil, errors.New("empty values")
	}
	for i := 0; i < rowsV.Len(); i++ {
		rowI := rowsV.Index(i).Interface()
		row, err := s.toValuesRow(rowI)
		if err != nil {
			return nil, errors.WithMessagef(err, "can't convert %T into SQL row", rowI)
		}
		input.Rows = append(input.Rows, row)
	}
	return input, nil
}

func (s *InsertStmt) toValuesRow(val interface{}) (*ast.ValuesRow, error) {
	valV := reflect.ValueOf(val)
	switch valV.Type().Kind() {
	case reflect.Slice:
		return s.sliceToValuesRow(valV)
	case reflect.Struct:
		return s.structToValuesRow(valV)
	case reflect.Ptr:
		if valV.Type().Elem().Kind() == reflect.Struct {
			return s.structToValuesRow(valV.Elem())
		}
		return nil, errors.Errorf("%s is neither struct nor slice", valV.Type().String())
	default:
		return nil, errors.Errorf("%s is neither struct nor slice", valV.Type().String())
	}
}

// The type of valV is guaranteed to be slice here.
func (s *InsertStmt) sliceToValuesRow(valV reflect.Value) (*ast.ValuesRow, error) {
	row := &ast.ValuesRow{}
	for i := 0; i < valV.Len(); i++ {
		expr, err := internal.ToExpr(valV.Index(i).Interface())
		if err != nil {
			return nil, err
		}
		row.Exprs = append(row.Exprs, &ast.DefaultExpr{Expr: expr})
	}
	return row, nil
}

// The type of valV is guaranteed to be struct here.
func (s *InsertStmt) structToValuesRow(valV reflect.Value) (*ast.ValuesRow, error) {
	row := &ast.ValuesRow{}
	valT := valV.Type()
	numField := valT.NumField()
	for _, colName := range s.cols {
		colFound := false
		for i := 0; i < numField; i++ {
			ft := valT.Field(i)
			if !columnNameMatches(&ft, colName) {
				continue
			}
			colFound = true
			expr, err := internal.ToExpr(valV.Field(i).Interface())
			if err != nil {
				return nil, err
			}
			row.Exprs = append(row.Exprs, &ast.DefaultExpr{Expr: expr})
		}
		if !colFound {
			return nil, errors.Errorf("type %s does not have column %s", valT.String(), colName)
		}
	}
	return row, nil
}

func columnNameMatches(field *reflect.StructField, colName string) bool {
	tag := field.Tag.Get("spanner")
	if tag == "" {
		return strings.EqualFold(field.Name, colName)
	} else if tag == "-" {
		return false
	}
	return tag == colName
}
