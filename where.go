package memeduck

import (
	"github.com/MakeNowJust/memefish/pkg/ast"
	"github.com/genkami/memeduck/internal"
)

// ExprCond is a boolean expression to filter records.
type ExprCond struct {
	expr ast.Expr
}

func (c *ExprCond) ToAstWhere() (*ast.Where, error) {
	return &ast.Where{
		Expr: c.expr,
	}, nil
}

// Bool creates a new boolean literal.
func Bool(v bool) *ExprCond {
	return &ExprCond{expr: internal.BoolLit(v)}
}

// OpCond is a binary operator expression.
type OpCond struct {
	lhs, rhs interface{}
	op       BinaryOp
}

// Op is a binary operator
type BinaryOp ast.BinaryOp

const (
	EQ BinaryOp = BinaryOp(ast.OpEqual)
	NE BinaryOp = BinaryOp(ast.OpNotEqual)
	LT BinaryOp = BinaryOp(ast.OpLess)
	GT BinaryOp = BinaryOp(ast.OpGreater)
	LE BinaryOp = BinaryOp(ast.OpLessEqual)
	GE BinaryOp = BinaryOp(ast.OpGreaterEqual)
)

func (c *OpCond) ToAstWhere() (*ast.Where, error) {
	lhs, err := internal.ToExpr(c.lhs)
	if err != nil {
		return nil, err
	}
	rhs, err := internal.ToExpr(c.rhs)
	if err != nil {
		return nil, err
	}
	return &ast.Where{
		Expr: &ast.BinaryExpr{
			Op:    ast.BinaryOp(c.op),
			Left:  lhs,
			Right: rhs,
		},
	}, nil
}

// Op creates a new binary operator expression.
func Op(lhs interface{}, op BinaryOp, rhs interface{}) *OpCond {
	return &OpCond{
		lhs: lhs,
		rhs: rhs,
		op:  op,
	}
}

// IdentExpr is an identifier.
type IdentExpr struct {
	name string
}

// Ident creates a new IdentExpr.
func Ident(name string) *IdentExpr {
	return &IdentExpr{name: name}
}

func (ie *IdentExpr) ToASTExpr() ast.Expr {
	return &ast.Ident{Name: ie.name}
}
