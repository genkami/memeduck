package memeduck_test

import (
	"math"
	"testing"

	"cloud.google.com/go/spanner"
	"github.com/stretchr/testify/assert"

	"github.com/genkami/memeduck"
)

func testInsert(t *testing.T, stmt *memeduck.InsertStmt, expected string) {
	actual, err := stmt.SQL()
	assert.Nil(t, err, expected)
	assert.Equal(t, expected, actual)
}

func TestInsertWithStringSlice(t *testing.T) {
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b"}, [][]string{
			{"foo", "bar"},
		}),
		`INSERT INTO hoge (a, b) VALUES ("foo", "bar")`,
	)
}

func TestInsertWithStringPtrSlice(t *testing.T) {
	var a = "foo"
	var b = "bar"
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b"}, [][]*string{
			{&a, &b},
		}),
		`INSERT INTO hoge (a, b) VALUES ("foo", "bar")`,
	)
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b"}, [][]*string{
			{nil, nil},
		}),
		`INSERT INTO hoge (a, b) VALUES (NULL, NULL)`,
	)
}

func TestInsertWithNullStringSlice(t *testing.T) {
	var a = spanner.NullString{StringVal: "foo", Valid: true}
	var b = spanner.NullString{StringVal: "bar", Valid: true}
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b"}, [][]spanner.NullString{
			{a, b},
		}),
		`INSERT INTO hoge (a, b) VALUES ("foo", "bar")`,
	)
	var null = spanner.NullString{}
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b"}, [][]spanner.NullString{
			{null, null},
		}),
		`INSERT INTO hoge (a, b) VALUES (NULL, NULL)`,
	)
}

func TestInsertWithByteSliceSlice(t *testing.T) {
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b"}, [][][]byte{
			{{0, 1}, {2, 3, 4}},
		}),
		`INSERT INTO hoge (a, b) VALUES (B"\x00\x01", B"\x02\x03\x04")`,
	)
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b"}, [][][]byte{
			{nil, nil},
		}),
		`INSERT INTO hoge (a, b) VALUES (NULL, NULL)`,
	)
}

func TestInsertWithIntSlice(t *testing.T) {
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b"}, [][]int{
			{123, 456},
		}),
		`INSERT INTO hoge (a, b) VALUES (123, 456)`,
	)
}

func TestInsertWithIntPtrSlice(t *testing.T) {
	var a = int(123)
	var b = int(456)
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b"}, [][]*int{
			{&a, &b},
		}),
		`INSERT INTO hoge (a, b) VALUES (123, 456)`,
	)
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b"}, [][]*int{
			{nil, nil},
		}),
		`INSERT INTO hoge (a, b) VALUES (NULL, NULL)`,
	)
}

func TestInsertWithInt64Slice(t *testing.T) {
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b"}, [][]int64{
			{123, 456},
		}),
		`INSERT INTO hoge (a, b) VALUES (123, 456)`,
	)
}

func TestInsertWithInt64PtrSlice(t *testing.T) {
	var a = int64(123)
	var b = int64(456)
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b"}, [][]*int64{
			{&a, &b},
		}),
		`INSERT INTO hoge (a, b) VALUES (123, 456)`,
	)
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b"}, [][]*int64{
			{nil, nil},
		}),
		`INSERT INTO hoge (a, b) VALUES (NULL, NULL)`,
	)
}

func TestInsertWithNullInt64Slice(t *testing.T) {
	var a = spanner.NullInt64{Int64: 123, Valid: true}
	var b = spanner.NullInt64{Int64: 456, Valid: true}
	var null = spanner.NullInt64{}
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b"}, [][]spanner.NullInt64{
			{a, b},
		}),
		`INSERT INTO hoge (a, b) VALUES (123, 456)`,
	)
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b"}, [][]spanner.NullInt64{
			{null, null},
		}),
		`INSERT INTO hoge (a, b) VALUES (NULL, NULL)`,
	)
}

func TestInsertWithBoolSlice(t *testing.T) {
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b"}, [][]bool{
			{true, false},
		}),
		`INSERT INTO hoge (a, b) VALUES (TRUE, FALSE)`,
	)
}

func TestInsertWithBoolPtrSlice(t *testing.T) {
	var a = true
	var b = false
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b"}, [][]*bool{
			{&a, &b},
		}),
		`INSERT INTO hoge (a, b) VALUES (TRUE, FALSE)`,
	)
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b"}, [][]*bool{
			{nil, nil},
		}),
		`INSERT INTO hoge (a, b) VALUES (NULL, NULL)`,
	)
}

func TestInsertWithNullBoolSlice(t *testing.T) {
	var a = spanner.NullBool{Bool: true, Valid: true}
	var b = spanner.NullBool{Bool: false, Valid: true}
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b"}, [][]spanner.NullBool{
			{a, b},
		}),
		`INSERT INTO hoge (a, b) VALUES (TRUE, FALSE)`,
	)
	var null = spanner.NullBool{}
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b"}, [][]spanner.NullBool{
			{null, null},
		}),
		`INSERT INTO hoge (a, b) VALUES (NULL, NULL)`,
	)
}

func TestInsertWithFloat64Slice(t *testing.T) {
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b", "c", "d", "e", "f"}, [][]float64{
			{
				1.0,
				0,
				3.1415926535,
				math.NaN(),
				math.Inf(1),
				math.Inf(-1),
			},
		}),
		`INSERT INTO hoge (a, b, c, d, e, f) VALUES (`+
			`1e+00, `+
			`0e+00, `+
			`3.1415926535e+00, `+
			`NaN, `+
			`+Inf, `+
			`-Inf)`,
	)
}

func TestInsertWithFloat64PtrSlice(t *testing.T) {
	var a float64 = 1.0
	var b float64 = 0
	var c float64 = 3.1415926535
	var d float64 = math.NaN()
	var e float64 = math.Inf(1)
	var f float64 = math.Inf(-1)
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b", "c", "d", "e", "f"}, [][]*float64{
			{&a, &b, &c, &d, &e, &f, nil},
		}),
		`INSERT INTO hoge (a, b, c, d, e, f) VALUES (`+
			`1e+00, `+
			`0e+00, `+
			`3.1415926535e+00, `+
			`NaN, `+
			`+Inf, `+
			`-Inf, `+
			`NULL)`,
	)
}

func TestInsertWithNullFloat64Slice(t *testing.T) {
	var a = spanner.NullFloat64{Float64: 1.0, Valid: true}
	var b = spanner.NullFloat64{Float64: 0, Valid: true}
	var c = spanner.NullFloat64{Float64: 3.1415926535, Valid: true}
	var d = spanner.NullFloat64{Float64: math.NaN(), Valid: true}
	var e = spanner.NullFloat64{Float64: math.Inf(1), Valid: true}
	var f = spanner.NullFloat64{Float64: math.Inf(-1), Valid: true}
	var g = spanner.NullFloat64{}
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b", "c", "d", "e", "f"}, [][]spanner.NullFloat64{
			{a, b, c, d, e, f, g},
		}),
		`INSERT INTO hoge (a, b, c, d, e, f) VALUES (`+
			`1e+00, `+
			`0e+00, `+
			`3.1415926535e+00, `+
			`NaN, `+
			`+Inf, `+
			`-Inf, `+
			`NULL)`,
	)
}
