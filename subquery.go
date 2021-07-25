package memeduck

import "github.com/MakeNowJust/memefish/pkg/ast"

type SubQuery interface {
	ToAST() (ast.SelectItem, error)
}

type ScalarSubQueryStmt struct {
	as    string
	query *SelectStmt
}

func ScalarSubQuery(stmt *SelectStmt) *ScalarSubQueryStmt {
	return &ScalarSubQueryStmt{
		query: stmt,
	}
}

func (s *ScalarSubQueryStmt) As(as string) *ScalarSubQueryStmt {
	var t = *s
	t.as = as
	return &t
}

func (s *ScalarSubQueryStmt) ToAST() (ast.SelectItem, error) {
	stmt, err := s.query.toAST()
	if err != nil {
		return nil, err
	}
	expr := &ast.ScalarSubQuery{
		Query: stmt,
	}
	if s.as == "" {
		return &ast.ExprSelectItem{
			Expr: expr,
		}, nil
	}
	return &ast.Alias{
		Expr: expr,
		As: &ast.AsAlias{
			Alias: &ast.Ident{
				Name: s.as,
			},
		},
	}, nil
}

type ArraySubQueryStmt struct {
	as    string
	query *SelectStmt
}

func ArraySubQuery(stmt *SelectStmt) *ArraySubQueryStmt {
	return &ArraySubQueryStmt{
		query: stmt,
	}
}

func (s *ArraySubQueryStmt) As(as string) *ArraySubQueryStmt {
	var t = *s
	t.as = as
	return &t
}

func (s *ArraySubQueryStmt) ToAST() (ast.SelectItem, error) {
	stmt, err := s.query.toAST()
	if err != nil {
		return nil, err
	}
	expr := &ast.ArraySubQuery{
		Query: stmt,
	}
	if s.as == "" {
		return &ast.ExprSelectItem{
			Expr: expr,
		}, nil
	}
	return &ast.Alias{
		Expr: expr,
		As: &ast.AsAlias{
			Alias: &ast.Ident{
				Name: s.as,
			},
		},
	}, nil
}
