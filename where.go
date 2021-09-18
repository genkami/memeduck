package memeduck

import (
	"github.com/MakeNowJust/memefish/pkg/ast"
	"github.com/pkg/errors"

	"github.com/genkami/memeduck/internal"
)

// WhereCond is a conditional expression that appears in WHERE clauses.
type WhereCond interface {
	ToASTWhere() (*ast.Where, error)
}

// ExprCond is a boolean expression to filter records.
type ExprCond struct {
	expr ast.Expr
}

func (c *ExprCond) ToASTWhere() (*ast.Where, error) {
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
	EQ       BinaryOp = BinaryOp(ast.OpEqual)
	NE       BinaryOp = BinaryOp(ast.OpNotEqual)
	LT       BinaryOp = BinaryOp(ast.OpLess)
	GT       BinaryOp = BinaryOp(ast.OpGreater)
	LE       BinaryOp = BinaryOp(ast.OpLessEqual)
	GE       BinaryOp = BinaryOp(ast.OpGreaterEqual)
	LIKE     BinaryOp = BinaryOp(ast.OpLike)
	NOT_LIKE BinaryOp = BinaryOp(ast.OpNotLike)
)

func (c *OpCond) ToASTWhere() (*ast.Where, error) {
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

// Eq(x, y) is a shorthand for Op(x, EQ, y)
func Eq(lhs, rhs interface{}) *OpCond {
	return Op(lhs, EQ, rhs)
}

// Ne(x, y) is a shorthand for Op(x, NE, y)
func Ne(lhs, rhs interface{}) *OpCond {
	return Op(lhs, NE, rhs)
}

// Lt(x, y) is a shorthand for Op(x, LT, y)
func Lt(lhs, rhs interface{}) *OpCond {
	return Op(lhs, LT, rhs)
}

// Gt(x, y) is a shorthand for Op(x, GT, y)
func Gt(lhs, rhs interface{}) *OpCond {
	return Op(lhs, GT, rhs)
}

// Le(x, y) is a shorthand for Op(x, LE, y)
func Le(lhs, rhs interface{}) *OpCond {
	return Op(lhs, LE, rhs)
}

// Ge(x, y) is a shorthand for Op(x, GE, y)
func Ge(lhs, rhs interface{}) *OpCond {
	return Op(lhs, GE, rhs)
}

// Like(x, y) is a shorthand for Op(x, LIKE, y)
func Like(lhs, rhs interface{}) *OpCond {
	return Op(lhs, LIKE, rhs)
}

// NotLike(x, y) is a shorthand for Op(x, NOT_LIKE, y)
func NotLike(lhs, rhs interface{}) *OpCond {
	return Op(lhs, NOT_LIKE, rhs)
}

// NullCond represents IS NULL or IS NOT NULL predicate.
type NullCond struct {
	not bool
	arg interface{}
}

// IsNull creates `x IS NULL` predicate.
func IsNull(arg interface{}) *NullCond {
	return &NullCond{arg: arg}
}

// IsNotNull creates `x IS NOT NULL` predicate.
func IsNotNull(arg interface{}) *NullCond {
	return &NullCond{arg: arg, not: true}
}

func (c *NullCond) ToASTWhere() (*ast.Where, error) {
	expr, err := internal.ToExpr(c.arg)
	if err != nil {
		return nil, err
	}
	return &ast.Where{
		Expr: &ast.IsNullExpr{
			Not:  c.not,
			Left: expr,
		},
	}, nil
}

// InCond represents IN or NOT IN predicates.
type InCond struct {
	lhs, rhs interface{}
	not      bool
}

// In(x, y) creates `x IN y` predicate.
func In(x, y interface{}) *InCond {
	return &InCond{lhs: x, rhs: y, not: false}
}

func NotIn(x, y interface{}) *InCond {
	return &InCond{lhs: x, rhs: y, not: true}
}

func (c *InCond) ToASTWhere() (*ast.Where, error) {
	lhs, err := internal.ToExpr(c.lhs)
	if err != nil {
		return nil, err
	}
	rhs, err := internal.ToExpr(c.rhs)
	if err != nil {
		return nil, err
	}
	return &ast.Where{
		Expr: &ast.InExpr{
			Not:  c.not,
			Left: lhs,
			Right: &ast.UnnestInCondition{
				Expr: rhs,
			},
		},
	}, nil
}

// BetweenCond represents BETWEEN or NOT BETWEEN predicates.
type BetweenCond struct {
	arg interface{}
	min interface{}
	max interface{}
	not bool
}

// Between(x, min, max) creates `x BETWEEN min AND max` predicate.
func Between(x, min, max interface{}) *BetweenCond {
	return &BetweenCond{arg: x, min: min, max: max}
}

// NotBetween(x, min, max) creates `x NOT BETWEEN min AND max` predicate.
func NotBetween(x, min, max interface{}) *BetweenCond {
	return &BetweenCond{arg: x, min: min, max: max, not: true}
}

func (c *BetweenCond) ToASTWhere() (*ast.Where, error) {
	arg, err := internal.ToExpr(c.arg)
	if err != nil {
		return nil, err
	}
	min, err := internal.ToExpr(c.min)
	if err != nil {
		return nil, err
	}
	max, err := internal.ToExpr(c.max)
	if err != nil {
		return nil, err
	}
	return &ast.Where{
		Expr: &ast.BetweenExpr{
			Not:        c.not,
			Left:       arg,
			RightStart: min,
			RightEnd:   max,
		},
	}, nil
}

// IdentExpr is an identifier.
type IdentExpr struct {
	names []string
}

// Ident creates a new IdentExpr.
// Path expression can be created by passing more than one elements.
func Ident(names ...string) *IdentExpr {
	return &IdentExpr{names: names}
}

func (e *IdentExpr) ToASTExpr() (ast.Expr, error) {
	if len(e.names) <= 0 {
		return nil, errors.New("empty identifier")
	}
	path := &ast.Path{}
	for _, name := range e.names {
		path.Idents = append(path.Idents, &ast.Ident{
			Name: name,
		})
	}
	return path, nil
}

// ParamExpr is a query parameter.
type ParamExpr struct {
	name string
}

// Param createsa new ParamExpr.
func Param(name string) *ParamExpr {
	return &ParamExpr{name: name}
}

func (e *ParamExpr) ToASTExpr() (ast.Expr, error) {
	return &ast.Param{Name: e.name}, nil
}

// LogicalOpCond represents AND/OR operator.
type LogicalOpCond struct {
	op    logicalOp
	conds []WhereCond
}

type logicalOp ast.BinaryOp

const (
	logicalOpAnd logicalOp = logicalOp(ast.OpAnd)
	logicalOpOr  logicalOp = logicalOp(ast.OpOr)
)

// And concatenates more than one WhereConds with AND operator.
func And(conds ...WhereCond) *LogicalOpCond {
	return &LogicalOpCond{
		op:    logicalOpAnd,
		conds: conds,
	}
}

// Or concatenates more than one WhereConds with OR operator.
func Or(conds ...WhereCond) *LogicalOpCond {
	return &LogicalOpCond{
		op:    logicalOpOr,
		conds: conds,
	}
}

func (c *LogicalOpCond) ToASTWhere() (*ast.Where, error) {
	if len(c.conds) <= 0 {
		return nil, errors.New("no conditions")
	}
	where, err := c.conds[0].ToASTWhere()
	if err != nil {
		return nil, err
	}
	acc := where
	for _, cond := range c.conds[1:] {
		where, err = cond.ToASTWhere()
		if err != nil {
			return nil, err
		}
		acc = &ast.Where{
			Expr: &ast.BinaryExpr{
				Op:    ast.BinaryOp(c.op),
				Left:  acc.Expr,
				Right: where.Expr,
			},
		}
	}
	return acc, nil
}
