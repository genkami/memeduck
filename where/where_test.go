package where_test

import (
	"testing"

	"github.com/MakeNowJust/memefish/pkg/ast"
	"github.com/genkami/memeduck/where"
	"github.com/stretchr/testify/assert"
)

type whereCond interface {
	ToAstWhere() (*ast.Where, error)
}

func test(t *testing.T, cond whereCond, expected string) {
	w, err := cond.ToAstWhere()
	assert.Nil(t, err, expected)
	assert.Equal(t, w.Expr.SQL(), expected)
}

func TestBool(t *testing.T) {
	test(t, where.Bool(true), `TRUE`)
	test(t, where.Bool(false), `FALSE`)
}

func TestOp(t *testing.T) {
	test(t, where.Op(1, where.EQ, 1), `1 = 1`)
	test(t, where.Op("hoge", where.NE, "fuga"), `"hoge" != "fuga"`)
	test(t, where.Op(1.23, where.LT, 4.56), `1.23e+00 < 4.56e+00`)
	test(t, where.Op(4.56, where.GT, 1.23), `4.56e+00 > 1.23e+00`)
	test(t, where.Op(1, where.LE, 2), `1 <= 2`)
	test(t, where.Op(2, where.GE, 1), `2 >= 1`)
}
