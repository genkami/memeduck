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

func testAST(t *testing.T, val interface{}, expected ast.Expr) {
	actual, err := internal.ToExpr(val)
	assert.Nil(t, err, "can't convert %#v into Expr", val)
	assert.Equal(t, expected, actual)
}

func TestASTWithNil(t *testing.T) {
	testAST(t, nil, internal.NullLit())
}

func TestASTWithString(t *testing.T) {
	testAST(t, "hoge", internal.StringLit("hoge"))
}

func TestASTWithStringPtr(t *testing.T) {
	var v = "hoge"
	testAST(t, &v, internal.StringLit("hoge"))
	testAST(t, (*string)(nil), internal.NullLit())
}

func TestASTWithNullString(t *testing.T) {
	testAST(t, spanner.NullString{StringVal: "hoge", Valid: true}, internal.StringLit("hoge"))
	testAST(t, spanner.NullString{}, internal.NullLit())
}

func TestASTWithBytes(t *testing.T) {
	testAST(t, []byte{0, 1, 2}, internal.BytesLit([]byte{0, 1, 2}))
}

func TestASTWithInt(t *testing.T) {
	testAST(t, int(123), internal.IntLit(123))
}

func TestASTWithIntPtr(t *testing.T) {
	var v int = 123
	testAST(t, &v, internal.IntLit(123))
	testAST(t, (*int)(nil), internal.NullLit())
}

func TestASTWithInt64(t *testing.T) {
	testAST(t, int64(123), internal.IntLit(123))
}

func TestASTWithInt64Ptr(t *testing.T) {
	var v int64 = 123
	testAST(t, &v, internal.IntLit(123))
	testAST(t, (*int64)(nil), internal.NullLit())
}

func TestASTWithNullInt64(t *testing.T) {
	testAST(t, spanner.NullInt64{Int64: 123, Valid: true}, internal.IntLit(123))
	testAST(t, spanner.NullInt64{}, internal.NullLit())
}

func TestASTWithBool(t *testing.T) {
	testAST(t, true, internal.BoolLit(true))
	testAST(t, false, internal.BoolLit(false))
}

func TestASTWithBoolPtr(t *testing.T) {
	var v bool = true
	testAST(t, &v, internal.BoolLit(true))
	testAST(t, (*bool)(nil), internal.NullLit())
}

func TestASTWithNullBool(t *testing.T) {
	testAST(t, spanner.NullBool{Bool: false, Valid: true}, internal.BoolLit(false))
	testAST(t, spanner.NullBool{}, internal.NullLit())
}

func TestASTWithFloat64(t *testing.T) {
	testAST(t, float64(3.14), internal.FloatLit(3.14))
}

func TestASTWithFloat64Ptr(t *testing.T) {
	var v float64 = 3.14
	testAST(t, &v, internal.FloatLit(3.14))
	testAST(t, (*float64)(nil), internal.NullLit())
}

func TestASTWithNullFloat64(t *testing.T) {
	testAST(t, spanner.NullFloat64{Float64: 1.23, Valid: true}, internal.FloatLit(1.23))
	testAST(t, spanner.NullFloat64{}, internal.NullLit())
}

func TestASTWithTime(t *testing.T) {
	var v = time.Now()
	testAST(t, v, internal.TimeLit(v))
}

func TestASTWithTimePtr(t *testing.T) {
	var v = time.Now()
	testAST(t, &v, internal.TimeLit(v))
	testAST(t, (*time.Time)(nil), internal.NullLit())
}

func TestASTWithNullTime(t *testing.T) {
	var v = time.Now()
	testAST(t, spanner.NullTime{Time: v, Valid: true}, internal.TimeLit(v))
	testAST(t, spanner.NullTime{}, internal.NullLit())
}

func TestASTWithDate(t *testing.T) {
	v, err := civil.ParseDate("2021-05-22")
	assert.Nil(t, err)
	testAST(t, v, internal.DateLit(v))
}

func TestASTWithDatePtr(t *testing.T) {
	v, err := civil.ParseDate("2021-05-22")
	assert.Nil(t, err)
	testAST(t, &v, internal.DateLit(v))
	testAST(t, (*civil.Date)(nil), internal.NullLit())
}

func TestASTWithNullDate(t *testing.T) {
	v, err := civil.ParseDate("2021-05-22")
	assert.Nil(t, err)
	testAST(t, spanner.NullDate{Date: v, Valid: true}, internal.DateLit(v))
	testAST(t, spanner.NullDate{}, internal.NullLit())
}

type customExpr struct{}

func (*customExpr) ToASTExpr() ast.Expr {
	return internal.StringLit("custom expr")
}

func TestASTWithSpannerExpr(t *testing.T) {
	testAST(t, &customExpr{}, internal.StringLit("custom expr"))
}

func TestASTWithSlice(t *testing.T) {
	testAST(t,
		[]interface{}{nil, nil},
		internal.ArrayLit([]ast.Expr{internal.NullLit(), internal.NullLit()}),
	)
	testAST(t,
		[]string{"hoge", "fuga"},
		internal.ArrayLit([]ast.Expr{internal.StringLit("hoge"), internal.StringLit("fuga")}),
	)
	testAST(t,
		[]interface{}{123, "456"},
		internal.ArrayLit([]ast.Expr{internal.IntLit(123), internal.StringLit("456")}),
	)
}
