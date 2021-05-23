package memeduck_test

import (
	"testing"

	"github.com/MakeNowJust/memefish/pkg/ast"
	"github.com/genkami/memeduck"
	"github.com/stretchr/testify/assert"
)

type whereCond interface {
	ToASTWhere() (*ast.Where, error)
}

func testWhere(t *testing.T, cond whereCond, expected string) {
	w, err := cond.ToASTWhere()
	assert.Nil(t, err, expected)
	assert.Equal(t, expected, w.Expr.SQL())
}

type astExpr interface {
	ToASTExpr() (ast.Expr, error)
}

func testExpr(t *testing.T, expr astExpr, expected string) {
	e, err := expr.ToASTExpr()
	assert.Nil(t, err, expected)
	assert.Equal(t, expected, e.SQL())
}

func TestBool(t *testing.T) {
	testWhere(t, memeduck.Bool(true), `TRUE`)
	testWhere(t, memeduck.Bool(false), `FALSE`)
}

func TestIdent(t *testing.T) {
	_, err := memeduck.Ident().ToASTExpr()
	assert.Error(t, err, "empty ident")
	testExpr(t, memeduck.Ident("a"), `a`)
	testExpr(t, memeduck.Ident("abc"), `abc`)
	testExpr(t, memeduck.Ident("TRUE"), "`TRUE`")
	testExpr(t, memeduck.Ident("a", "b"), `a.b`)
	testExpr(t, memeduck.Ident("TRUE", "FALSE"), "`TRUE`.`FALSE`")
}

func TestParam(t *testing.T) {
	testExpr(t, memeduck.Param("a"), `@a`)
	testExpr(t, memeduck.Param("abc"), `@abc`)
}

func TestOp(t *testing.T) {
	testWhere(t, memeduck.Op(1, memeduck.EQ, 1), `1 = 1`)
	testWhere(t, memeduck.Op("hoge", memeduck.NE, "fuga"), `"hoge" != "fuga"`)
	testWhere(t, memeduck.Op(1.23, memeduck.LT, 4.56), `1.23e+00 < 4.56e+00`)
	testWhere(t, memeduck.Op(4.56, memeduck.GT, 1.23), `4.56e+00 > 1.23e+00`)
	testWhere(t, memeduck.Op(1, memeduck.LE, 2), `1 <= 2`)
	testWhere(t, memeduck.Op(2, memeduck.GE, 1), `2 >= 1`)
	testWhere(t, memeduck.Op("hoge", memeduck.LIKE, "ho%"), `"hoge" LIKE "ho%"`)
	testWhere(t, memeduck.Op("hoge", memeduck.NOT_LIKE, "ho%"), `"hoge" NOT LIKE "ho%"`)

	testWhere(t, memeduck.Eq(1, 1), `1 = 1`)
	testWhere(t, memeduck.Ne("hoge", "fuga"), `"hoge" != "fuga"`)
	testWhere(t, memeduck.Lt(1.23, 4.56), `1.23e+00 < 4.56e+00`)
	testWhere(t, memeduck.Gt(4.56, 1.23), `4.56e+00 > 1.23e+00`)
	testWhere(t, memeduck.Le(1, 2), `1 <= 2`)
	testWhere(t, memeduck.Ge(2, 1), `2 >= 1`)
	testWhere(t, memeduck.Like("hoge", "ho%"), `"hoge" LIKE "ho%"`)
	testWhere(t, memeduck.NotLike("hoge", "ho%"), `"hoge" NOT LIKE "ho%"`)
}

func TestIsNullAndIsNotNull(t *testing.T) {
	testWhere(t, memeduck.IsNull(memeduck.Ident("hoge")), `hoge IS NULL`)
	testWhere(t, memeduck.IsNotNull(memeduck.Ident("fuga")), `fuga IS NOT NULL`)
}

func TestAnd(t *testing.T) {
	_, err := memeduck.And().ToASTWhere()
	assert.Error(t, err, "empty AND")
	testWhere(t,
		memeduck.And(
			memeduck.Op(1, memeduck.EQ, 1),
		),
		`1 = 1`,
	)
	testWhere(t,
		memeduck.And(
			memeduck.Op(1, memeduck.EQ, 1),
			memeduck.Op("hoge", memeduck.EQ, "hoge"),
		),
		`1 = 1 AND "hoge" = "hoge"`,
	)
	testWhere(t,
		memeduck.And(
			memeduck.Op(1, memeduck.EQ, 1),
			memeduck.Op("hoge", memeduck.EQ, "hoge"),
			memeduck.Op(true, memeduck.EQ, true),
		),
		`1 = 1 AND "hoge" = "hoge" AND TRUE = TRUE`,
	)
}

func TestOr(t *testing.T) {
	_, err := memeduck.Or().ToASTWhere()
	assert.Error(t, err, "empty Or")
	testWhere(t,
		memeduck.Or(
			memeduck.Op(1, memeduck.EQ, 1),
		),
		`1 = 1`,
	)
	testWhere(t,
		memeduck.Or(
			memeduck.Op(1, memeduck.EQ, 1),
			memeduck.Op("hoge", memeduck.EQ, "hoge"),
		),
		`1 = 1 OR "hoge" = "hoge"`,
	)
	testWhere(t,
		memeduck.Or(
			memeduck.Op(1, memeduck.EQ, 1),
			memeduck.Op("hoge", memeduck.EQ, "hoge"),
			memeduck.Op(true, memeduck.EQ, true),
		),
		`1 = 1 OR "hoge" = "hoge" OR TRUE = TRUE`,
	)

	// TODO: this shoud pass
	// testWhere(t,
	// 	memeduck.And(
	// 		memeduck.Eq(1, 1),
	// 		memeduck.Or(
	// 			memeduck.Eq(2, 2),
	// 			memeduck.Eq(3, 3),
	// 		),
	// 	),
	// 	`1 = 1 AND (2 = 2 OR 3 = 3)`,
	// )
}
