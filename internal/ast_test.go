package internal_test

import (
	"testing"
	"time"

	"cloud.google.com/go/civil"

	"cloud.google.com/go/spanner"
	"github.com/MakeNowJust/memefish/pkg/ast"
	"github.com/stretchr/testify/assert"

	"github.com/genkami/memeduck/internal"
)

func testAst(t *testing.T, val interface{}, expected ast.Expr) {
	actual, err := internal.ToExpr(val)
	assert.Nil(t, err, "can't convert %#v into Expr", val)
	assert.Equal(t, expected, actual)
}

func TestAstWithNil(t *testing.T) {
	testAst(t, nil, internal.NullLit())
}

func TestAstWithString(t *testing.T) {
	testAst(t, "hoge", internal.StringLit("hoge"))
}

func TestAstWithStringPtr(t *testing.T) {
	var v = "hoge"
	testAst(t, &v, internal.StringLit("hoge"))
	testAst(t, (*string)(nil), internal.NullLit())
}

func TestAstWithNullString(t *testing.T) {
	testAst(t, spanner.NullString{StringVal: "hoge", Valid: true}, internal.StringLit("hoge"))
	testAst(t, spanner.NullString{}, internal.NullLit())
}

func TestAstWithBytes(t *testing.T) {
	testAst(t, []byte{0, 1, 2}, internal.BytesLit([]byte{0, 1, 2}))
}

func TestAstWithInt(t *testing.T) {
	testAst(t, int(123), internal.IntLit(123))
}

func TestAstWithIntPtr(t *testing.T) {
	var v int = 123
	testAst(t, &v, internal.IntLit(123))
	testAst(t, (*int)(nil), internal.NullLit())
}

func TestAstWithInt64(t *testing.T) {
	testAst(t, int64(123), internal.IntLit(123))
}

func TestAstWithInt64Ptr(t *testing.T) {
	var v int64 = 123
	testAst(t, &v, internal.IntLit(123))
	testAst(t, (*int64)(nil), internal.NullLit())
}

func TestAstWithNullInt64(t *testing.T) {
	testAst(t, spanner.NullInt64{Int64: 123, Valid: true}, internal.IntLit(123))
	testAst(t, spanner.NullInt64{}, internal.NullLit())
}

func TestAstWithBool(t *testing.T) {
	testAst(t, true, internal.BoolLit(true))
	testAst(t, false, internal.BoolLit(false))
}

func TestAstWithBoolPtr(t *testing.T) {
	var v bool = true
	testAst(t, &v, internal.BoolLit(true))
	testAst(t, (*bool)(nil), internal.NullLit())
}

func TestAstWithNullBool(t *testing.T) {
	testAst(t, spanner.NullBool{Bool: false, Valid: true}, internal.BoolLit(false))
	testAst(t, spanner.NullBool{}, internal.NullLit())
}

func TestAstWithFloat64(t *testing.T) {
	testAst(t, float64(3.14), internal.FloatLit(3.14))
}

func TestAstWithFloat64Ptr(t *testing.T) {
	var v float64 = 3.14
	testAst(t, &v, internal.FloatLit(3.14))
	testAst(t, (*float64)(nil), internal.NullLit())
}

func TestAstWithNullFloat64(t *testing.T) {
	testAst(t, spanner.NullFloat64{Float64: 1.23, Valid: true}, internal.FloatLit(1.23))
	testAst(t, spanner.NullFloat64{}, internal.NullLit())
}

func TestAstWithTime(t *testing.T) {
	var v = time.Now()
	testAst(t, v, internal.TimeLit(v))
}

func TestAstWithTimePtr(t *testing.T) {
	var v = time.Now()
	testAst(t, &v, internal.TimeLit(v))
	testAst(t, (*time.Time)(nil), internal.NullLit())
}

func TestAstWithNullTime(t *testing.T) {
	var v = time.Now()
	testAst(t, spanner.NullTime{Time: v, Valid: true}, internal.TimeLit(v))
	testAst(t, spanner.NullTime{}, internal.NullLit())
}

func TestAstWithDate(t *testing.T) {
	v, err := civil.ParseDate("2021-05-22")
	assert.Nil(t, err)
	testAst(t, v, internal.DateLit(v))
}

func TestAstWithDatePtr(t *testing.T) {
	v, err := civil.ParseDate("2021-05-22")
	assert.Nil(t, err)
	testAst(t, &v, internal.DateLit(v))
	testAst(t, (*civil.Date)(nil), internal.NullLit())
}

func TestAstWithNullDate(t *testing.T) {
	v, err := civil.ParseDate("2021-05-22")
	assert.Nil(t, err)
	testAst(t, spanner.NullDate{Date: v, Valid: true}, internal.DateLit(v))
	testAst(t, spanner.NullDate{}, internal.NullLit())
}

type customExpr struct{}

func (*customExpr) ToASTExpr() ast.Expr {
	return internal.StringLit("custom expr")
}

func TestAstWithSpannerExpr(t *testing.T) {
	testAst(t, &customExpr{}, internal.StringLit("custom expr"))
}

func TestAstWithSlice(t *testing.T) {
	testAst(t,
		[]interface{}{nil, nil},
		internal.ArrayLit([]ast.Expr{internal.NullLit(), internal.NullLit()}),
	)
	testAst(t,
		[]string{"hoge", "fuga"},
		internal.ArrayLit([]ast.Expr{internal.StringLit("hoge"), internal.StringLit("fuga")}),
	)
	testAst(t,
		[]interface{}{123, "456"},
		internal.ArrayLit([]ast.Expr{internal.IntLit(123), internal.StringLit("456")}),
	)
}
