package memeduck_test

import (
	"math"
	"testing"
	"time"

	"cloud.google.com/go/civil"
	"cloud.google.com/go/spanner"
	"github.com/MakeNowJust/memefish/pkg/ast"
	"github.com/stretchr/testify/assert"

	"github.com/genkami/memeduck"
	"github.com/genkami/memeduck/internal"
)

func testInsert(t *testing.T, stmt *memeduck.InsertStmt, expected string) {
	actual, err := stmt.SQL()
	assert.Nil(t, err, expected)
	assert.Equal(t, expected, actual)
}

func TestInsertWithNoValuesClause(t *testing.T) {
	_, err := memeduck.Insert("hoge", []string{"a", "b"}).SQL()
	assert.Error(t, err, "no VALUES")
}

func TestInsertWithEmpty(t *testing.T) {
	_, err := memeduck.Insert("hoge", []string{"a", "b"}).Values([][]int{}).SQL()
	assert.Error(t, err, "empty values")
}

func TestInsertWithNonSliceArgument(t *testing.T) {
	_, err := memeduck.Insert("hoge", []string{"a", "b"}).Values(map[string]string{"hoge": "fuga"}).SQL()
	assert.Error(t, err, "non-slice argument")
}

func TestInsertWithNilInterfaceSlice(t *testing.T) {
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b"}).Values([][]interface{}{
			{nil, nil},
		}),
		`INSERT INTO hoge (a, b) VALUES (NULL, NULL)`,
	)
}

func TestInsertWithStringSlice(t *testing.T) {
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b"}).Values([][]string{
			{"foo", "bar"},
		}),
		`INSERT INTO hoge (a, b) VALUES ("foo", "bar")`,
	)
}

func TestInsertWithStringSliceSlice(t *testing.T) {
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b"}).Values([][][]string{
			{{}, {"a"}, {"b", "c"}},
		}),
		`INSERT INTO hoge (a, b) VALUES (ARRAY[], ARRAY["a"], ARRAY["b", "c"])`,
	)
}

func TestInsertWithStringPtrSlice(t *testing.T) {
	var a = "foo"
	var b = "bar"
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b"}).Values([][]*string{
			{&a, &b},
		}),
		`INSERT INTO hoge (a, b) VALUES ("foo", "bar")`,
	)
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b"}).Values([][]*string{
			{nil, nil},
		}),
		`INSERT INTO hoge (a, b) VALUES (NULL, NULL)`,
	)
}

func TestInsertWithStringPtrSliceSlice(t *testing.T) {
	var a = "foo"
	var b = "bar"
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b", "c"}).Values([][][]*string{
			{{}, {&a}, {&b, nil}},
		}),
		`INSERT INTO hoge (a, b, c) VALUES (ARRAY[], ARRAY["foo"], ARRAY["bar", NULL])`,
	)
}

func TestInsertWithNullStringSlice(t *testing.T) {
	var a = spanner.NullString{StringVal: "foo", Valid: true}
	var b = spanner.NullString{StringVal: "bar", Valid: true}
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b"}).Values([][]spanner.NullString{
			{a, b},
		}),
		`INSERT INTO hoge (a, b) VALUES ("foo", "bar")`,
	)
	var null = spanner.NullString{}
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b"}).Values([][]spanner.NullString{
			{null, null},
		}),
		`INSERT INTO hoge (a, b) VALUES (NULL, NULL)`,
	)
}

func TestInsertWithNullStringSliceSlice(t *testing.T) {
	var a = spanner.NullString{StringVal: "foo", Valid: true}
	var b = spanner.NullString{StringVal: "bar", Valid: true}
	var null = spanner.NullString{}
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b", "c"}).Values([][][]spanner.NullString{
			{{}, {a}, {b, null}},
		}),
		`INSERT INTO hoge (a, b, c) VALUES (ARRAY[], ARRAY["foo"], ARRAY["bar", NULL])`,
	)
}

func TestInsertWithByteSliceSlice(t *testing.T) {
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b"}).Values([][][]byte{
			{{0, 1}, {2, 3, 4}},
		}),
		`INSERT INTO hoge (a, b) VALUES (B"\x00\x01", B"\x02\x03\x04")`,
	)
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b"}).Values([][][]byte{
			{nil, nil},
		}),
		`INSERT INTO hoge (a, b) VALUES (NULL, NULL)`,
	)
}

func TestInsertWithByteSliceSliceSlice(t *testing.T) {
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b", "c"}).Values([][][][]byte{
			{{}, {{0, 1}}, {{2, 3, 4}, nil}},
		}),
		`INSERT INTO hoge (a, b, c) VALUES (ARRAY[], ARRAY[B"\x00\x01"], ARRAY[B"\x02\x03\x04", NULL])`,
	)
}

func TestInsertWithIntSlice(t *testing.T) {
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b"}).Values([][]int{
			{123, 456},
		}),
		`INSERT INTO hoge (a, b) VALUES (123, 456)`,
	)
}

func TestInsertWithIntSliceSlice(t *testing.T) {
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b", "c"}).Values([][][]int{
			{{}, {123}, {456, 789}},
		}),
		`INSERT INTO hoge (a, b, c) VALUES (ARRAY[], ARRAY[123], ARRAY[456, 789])`,
	)
}

func TestInsertWithIntPtrSlice(t *testing.T) {
	var a = int(123)
	var b = int(456)
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b"}).Values([][]*int{
			{&a, &b},
		}),
		`INSERT INTO hoge (a, b) VALUES (123, 456)`,
	)
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b"}).Values([][]*int{
			{nil, nil},
		}),
		`INSERT INTO hoge (a, b) VALUES (NULL, NULL)`,
	)
}

func TestInsertWithIntPtrSliceSlice(t *testing.T) {
	var a = int(123)
	var b = int(456)
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b", "c"}).Values([][][]*int{
			{{}, {&a}, {&b, nil}},
		}),
		`INSERT INTO hoge (a, b, c) VALUES (ARRAY[], ARRAY[123], ARRAY[456, NULL])`,
	)
}

func TestInsertWithInt64Slice(t *testing.T) {
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b"}).Values([][]int64{
			{123, 456},
		}),
		`INSERT INTO hoge (a, b) VALUES (123, 456)`,
	)
}

func TestInsertWithInt64SliceSlice(t *testing.T) {
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b", "c"}).Values([][][]int64{
			{{}, {123}, {456, 789}},
		}),
		`INSERT INTO hoge (a, b, c) VALUES (ARRAY[], ARRAY[123], ARRAY[456, 789])`,
	)
}

func TestInsertWithInt64PtrSlice(t *testing.T) {
	var a = int64(123)
	var b = int64(456)
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b"}).Values([][]*int64{
			{&a, &b},
		}),
		`INSERT INTO hoge (a, b) VALUES (123, 456)`,
	)
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b"}).Values([][]*int64{
			{nil, nil},
		}),
		`INSERT INTO hoge (a, b) VALUES (NULL, NULL)`,
	)
}

func TestInsertWithInt64PtrSliceSlice(t *testing.T) {
	var a = int64(123)
	var b = int64(456)
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b", "c"}).Values([][][]*int64{
			{{}, {&a}, {&b, nil}},
		}),
		`INSERT INTO hoge (a, b, c) VALUES (ARRAY[], ARRAY[123], ARRAY[456, NULL])`,
	)
}

func TestInsertWithNullInt64Slice(t *testing.T) {
	var a = spanner.NullInt64{Int64: 123, Valid: true}
	var b = spanner.NullInt64{Int64: 456, Valid: true}
	var null = spanner.NullInt64{}
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b"}).Values([][]spanner.NullInt64{
			{a, b},
		}),
		`INSERT INTO hoge (a, b) VALUES (123, 456)`,
	)
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b"}).Values([][]spanner.NullInt64{
			{null, null},
		}),
		`INSERT INTO hoge (a, b) VALUES (NULL, NULL)`,
	)
}

func TestInsertWithNullInt64SliceSlice(t *testing.T) {
	var a = spanner.NullInt64{Int64: 123, Valid: true}
	var b = spanner.NullInt64{Int64: 456, Valid: true}
	var null = spanner.NullInt64{}
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b", "c"}).Values([][][]spanner.NullInt64{
			{{}, {a}, {b, null}},
		}),
		`INSERT INTO hoge (a, b, c) VALUES (ARRAY[], ARRAY[123], ARRAY[456, NULL])`,
	)
}

func TestInsertWithBoolSlice(t *testing.T) {
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b"}).Values([][]bool{
			{true, false},
		}),
		`INSERT INTO hoge (a, b) VALUES (TRUE, FALSE)`,
	)
}

func TestInsertWithBoolSliceSlice(t *testing.T) {
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b", "c"}).Values([][][]bool{
			{{}, {true}, {false, true}},
		}),
		`INSERT INTO hoge (a, b, c) VALUES (ARRAY[], ARRAY[TRUE], ARRAY[FALSE, TRUE])`,
	)
}

func TestInsertWithBoolPtrSlice(t *testing.T) {
	var a = true
	var b = false
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b"}).Values([][]*bool{
			{&a, &b},
		}),
		`INSERT INTO hoge (a, b) VALUES (TRUE, FALSE)`,
	)
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b"}).Values([][]*bool{
			{nil, nil},
		}),
		`INSERT INTO hoge (a, b) VALUES (NULL, NULL)`,
	)
}

func TestInsertWithBoolPtrSliceSlice(t *testing.T) {
	var a = true
	var b = false
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b", "c"}).Values([][][]*bool{
			{{}, {&a}, {&b, nil}},
		}),
		`INSERT INTO hoge (a, b, c) VALUES (ARRAY[], ARRAY[TRUE], ARRAY[FALSE, NULL])`,
	)
}

func TestInsertWithNullBoolSlice(t *testing.T) {
	var a = spanner.NullBool{Bool: true, Valid: true}
	var b = spanner.NullBool{Bool: false, Valid: true}
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b"}).Values([][]spanner.NullBool{
			{a, b},
		}),
		`INSERT INTO hoge (a, b) VALUES (TRUE, FALSE)`,
	)
	var null = spanner.NullBool{}
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b"}).Values([][]spanner.NullBool{
			{null, null},
		}),
		`INSERT INTO hoge (a, b) VALUES (NULL, NULL)`,
	)
}

func TestInsertWithNullBoolSliceSlice(t *testing.T) {
	var a = spanner.NullBool{Bool: true, Valid: true}
	var b = spanner.NullBool{Bool: false, Valid: true}
	var null = spanner.NullBool{}
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b", "c"}).Values([][][]spanner.NullBool{
			{{}, {a}, {b, null}},
		}),
		`INSERT INTO hoge (a, b, c) VALUES (ARRAY[], ARRAY[TRUE], ARRAY[FALSE, NULL])`,
	)
}

func TestInsertWithFloat64Slice(t *testing.T) {
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b", "c", "d", "e", "f"}).Values([][]float64{
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

func TestInsertWithFloat64SliceSlice(t *testing.T) {
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b", "c"}).Values([][][]float64{
			{{}, {0}, {31.5, math.Inf(1)}},
		}),
		`INSERT INTO hoge (a, b, c) VALUES (ARRAY[], ARRAY[0e+00], ARRAY[3.15e+01, +Inf])`,
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
		memeduck.Insert("hoge", []string{"a", "b", "c", "d", "e", "f"}).Values([][]*float64{
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

func TestInsertWithFloat64PtrSliceSlice(t *testing.T) {
	var a float64 = 0
	var b float64 = 31.5
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b", "c"}).Values([][][]*float64{
			{{}, {&a}, {&b, nil}},
		}),
		`INSERT INTO hoge (a, b, c) VALUES (ARRAY[], ARRAY[0e+00], ARRAY[3.15e+01, NULL])`,
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
		memeduck.Insert("hoge", []string{"a", "b", "c", "d", "e", "f"}).Values([][]spanner.NullFloat64{
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

func TestInsertWithNullFloat64SliceSlice(t *testing.T) {
	var a = spanner.NullFloat64{Float64: 0, Valid: true}
	var b = spanner.NullFloat64{Float64: 31.5, Valid: true}
	var null = spanner.NullFloat64{}
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b", "c"}).Values([][][]spanner.NullFloat64{
			{{}, {a}, {b, null}},
		}),
		`INSERT INTO hoge (a, b, c) VALUES (ARRAY[], ARRAY[0e+00], ARRAY[3.15e+01, NULL])`,
	)
}

func TestInsertWithTimeSlice(t *testing.T) {
	var a = parseTime(t, "2020-06-06T12:34:56.123456Z")
	var b = parseTime(t, "2021-08-10T00:01:23.456789+09:00")
	var c = parseTime(t, "2022-12-08T14:22:51.837583-04:30")
	var d = parseTime(t, "2023-10-10T08:43:17.536829+00:00")
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b", "c", "d"}).Values([][]time.Time{
			{a, b, c, d},
		}),
		`INSERT INTO hoge (a, b, c, d) VALUES (`+
			`TIMESTAMP "2020-06-06T12:34:56.123456Z", `+
			`TIMESTAMP "2021-08-10T00:01:23.456789+09:00", `+
			`TIMESTAMP "2022-12-08T14:22:51.837583-04:30", `+
			`TIMESTAMP "2023-10-10T08:43:17.536829Z")`,
	)
}

func TestInsertWithTimeSliceSlice(t *testing.T) {
	var a = parseTime(t, "2020-06-06T12:34:56.123456Z")
	var b = parseTime(t, "2021-08-10T00:01:23.456789+09:00")
	var c = parseTime(t, "2022-12-08T14:22:51.837583-04:30")
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b", "c"}).Values([][][]time.Time{
			{{}, {a}, {b, c}},
		}),
		`INSERT INTO hoge (a, b, c) VALUES (`+
			`ARRAY[], `+
			`ARRAY[TIMESTAMP "2020-06-06T12:34:56.123456Z"], `+
			`ARRAY[TIMESTAMP "2021-08-10T00:01:23.456789+09:00", `+
			`TIMESTAMP "2022-12-08T14:22:51.837583-04:30"])`,
	)
}

func TestInsertWithTimePtrSlice(t *testing.T) {
	var a = parseTime(t, "2020-06-06T12:34:56.123456Z")
	var b = parseTime(t, "2021-08-10T00:01:23.456789+09:00")
	var c = parseTime(t, "2022-12-08T14:22:51.837583-04:30")
	var d = parseTime(t, "2023-10-10T08:43:17.536829+00:00")
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b", "c", "d"}).Values([][]*time.Time{
			{&a, &b, &c, &d, nil},
		}),
		`INSERT INTO hoge (a, b, c, d) VALUES (`+
			`TIMESTAMP "2020-06-06T12:34:56.123456Z", `+
			`TIMESTAMP "2021-08-10T00:01:23.456789+09:00", `+
			`TIMESTAMP "2022-12-08T14:22:51.837583-04:30", `+
			`TIMESTAMP "2023-10-10T08:43:17.536829Z", `+
			`NULL)`,
	)
}

func TestInsertWithTimePtrSliceSlice(t *testing.T) {
	var a = parseTime(t, "2020-06-06T12:34:56.123456Z")
	var b = parseTime(t, "2021-08-10T00:01:23.456789+09:00")
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b", "c"}).Values([][][]*time.Time{
			{{}, {&a}, {&b, nil}},
		}),
		`INSERT INTO hoge (a, b, c) VALUES (`+
			`ARRAY[], `+
			`ARRAY[TIMESTAMP "2020-06-06T12:34:56.123456Z"], `+
			`ARRAY[TIMESTAMP "2021-08-10T00:01:23.456789+09:00", NULL])`,
	)
}

func TestInsertWithNullTimeSlice(t *testing.T) {
	var a = spanner.NullTime{Time: parseTime(t, "2020-06-06T12:34:56.123456Z"), Valid: true}
	var b = spanner.NullTime{Time: parseTime(t, "2021-08-10T00:01:23.456789+09:00"), Valid: true}
	var c = spanner.NullTime{Time: parseTime(t, "2022-12-08T14:22:51.837583-04:30"), Valid: true}
	var d = spanner.NullTime{Time: parseTime(t, "2023-10-10T08:43:17.536829+00:00"), Valid: true}
	var e = spanner.NullTime{}
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b", "c", "d"}).Values([][]spanner.NullTime{
			{a, b, c, d, e},
		}),
		`INSERT INTO hoge (a, b, c, d) VALUES (`+
			`TIMESTAMP "2020-06-06T12:34:56.123456Z", `+
			`TIMESTAMP "2021-08-10T00:01:23.456789+09:00", `+
			`TIMESTAMP "2022-12-08T14:22:51.837583-04:30", `+
			`TIMESTAMP "2023-10-10T08:43:17.536829Z", `+
			`NULL)`,
	)
}

func TestInsertWithNullTimeSliceSlice(t *testing.T) {
	var a = spanner.NullTime{Time: parseTime(t, "2020-06-06T12:34:56.123456Z"), Valid: true}
	var b = spanner.NullTime{Time: parseTime(t, "2021-08-10T00:01:23.456789+09:00"), Valid: true}
	var null = spanner.NullTime{}
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b", "c"}).Values([][][]spanner.NullTime{
			{{}, {a}, {b, null}},
		}),
		`INSERT INTO hoge (a, b, c) VALUES (`+
			`ARRAY[], `+
			`ARRAY[TIMESTAMP "2020-06-06T12:34:56.123456Z"], `+
			`ARRAY[TIMESTAMP "2021-08-10T00:01:23.456789+09:00", NULL])`,
	)
}

func TestInsertWithDateSlice(t *testing.T) {
	var a = parseDate(t, "2024-03-02")
	var b = parseDate(t, "2025-06-20")
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b"}).Values([][]civil.Date{
			{a, b},
		}),
		`INSERT INTO hoge (a, b) VALUES (`+
			`DATE "2024-03-02", `+
			`DATE "2025-06-20")`,
	)
}

func TestInsertWithDateSliceSlice(t *testing.T) {
	var a = parseDate(t, "2024-03-02")
	var b = parseDate(t, "2025-06-20")
	var c = parseDate(t, "2026-03-05")
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b", "c"}).Values([][][]civil.Date{
			{{}, {a}, {b, c}},
		}),
		`INSERT INTO hoge (a, b, c) VALUES (`+
			`ARRAY[], `+
			`ARRAY[DATE "2024-03-02"], `+
			`ARRAY[DATE "2025-06-20", DATE "2026-03-05"])`,
	)
}

func TestInsertWithDatePtrSlice(t *testing.T) {
	var a = parseDate(t, "2024-03-02")
	var b = parseDate(t, "2025-06-20")
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b", "c"}).Values([][]*civil.Date{
			{&a, &b, nil},
		}),
		`INSERT INTO hoge (a, b, c) VALUES (`+
			`DATE "2024-03-02", `+
			`DATE "2025-06-20", `+
			`NULL)`,
	)
}

func TestInsertWithDatePtrSliceSlice(t *testing.T) {
	var a = parseDate(t, "2024-03-02")
	var b = parseDate(t, "2025-06-20")
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b", "c"}).Values([][][]*civil.Date{
			{{}, {&a}, {&b, nil}},
		}),
		`INSERT INTO hoge (a, b, c) VALUES (`+
			`ARRAY[], `+
			`ARRAY[DATE "2024-03-02"], `+
			`ARRAY[DATE "2025-06-20", NULL])`,
	)
}
func TestInsertWithNullDateSlice(t *testing.T) {
	var a = spanner.NullDate{Date: parseDate(t, "2024-03-02"), Valid: true}
	var b = spanner.NullDate{Date: parseDate(t, "2025-06-20"), Valid: true}
	var c = spanner.NullDate{}
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b", "c"}).Values([][]spanner.NullDate{
			{a, b, c},
		}),
		`INSERT INTO hoge (a, b, c) VALUES (`+
			`DATE "2024-03-02", `+
			`DATE "2025-06-20", `+
			`NULL)`,
	)
}

func TestInsertWithNullDateSliceSlice(t *testing.T) {
	var a = spanner.NullDate{Date: parseDate(t, "2024-03-02"), Valid: true}
	var b = spanner.NullDate{Date: parseDate(t, "2025-06-20"), Valid: true}
	var null = spanner.NullDate{}
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b", "c"}).Values([][][]spanner.NullDate{
			{{}, {a}, {b, null}},
		}),
		`INSERT INTO hoge (a, b, c) VALUES (`+
			`ARRAY[], `+
			`ARRAY[DATE "2024-03-02"], `+
			`ARRAY[DATE "2025-06-20", NULL])`,
	)
}

func TestInsertWithParamSlice(t *testing.T) {
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b"}).Values([][]*memeduck.ParamExpr{
			{memeduck.Param("a"), memeduck.Param("b")},
		}),
		`INSERT INTO hoge (a, b) VALUES (@a, @b)`,
	)
}

type testInsertCustomExpr struct {
	a, b string
}

func (e *testInsertCustomExpr) ToASTExpr() (ast.Expr, error) {
	return internal.ArrayLit([]ast.Expr{
		internal.StringLit(e.a),
		internal.StringLit(e.b),
	}), nil
}

func TestInsertWithCustomExpr(t *testing.T) {
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b"}).Values([][]*testInsertCustomExpr{
			{&testInsertCustomExpr{"aaa", "bbb"}, &testInsertCustomExpr{"ccc", "ddd"}},
		}),
		`INSERT INTO hoge (a, b) VALUES (ARRAY["aaa", "bbb"], ARRAY["ccc", "ddd"])`,
	)
}

type testInsertGoStruct struct {
	A, B, C string
}

func TestInsertWithGoStruct(t *testing.T) {
	testInsert(t,
		memeduck.Insert("hoge", []string{"A", "B", "C"}).Values([]testInsertGoStruct{
			testInsertGoStruct{A: "AAA", B: "BBB", C: "CCC"},
		}),
		`INSERT INTO hoge (A, B, C) VALUES ("AAA", "BBB", "CCC")`,
	)
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b", "c"}).Values([]testInsertGoStruct{
			testInsertGoStruct{A: "AAA", B: "BBB", C: "CCC"},
		}),
		`INSERT INTO hoge (a, b, c) VALUES ("AAA", "BBB", "CCC")`,
	)
	testInsert(t,
		memeduck.Insert("hoge", []string{"A", "B"}).Values([]testInsertGoStruct{
			testInsertGoStruct{A: "AAA", B: "BBB", C: "CCC"},
		}),
		`INSERT INTO hoge (A, B) VALUES ("AAA", "BBB")`,
	)
	testInsert(t,
		memeduck.Insert("hoge", []string{"A", "B", "C"}).Values([]*testInsertGoStruct{
			&testInsertGoStruct{A: "AAA", B: "BBB", C: "CCC"},
		}),
		`INSERT INTO hoge (A, B, C) VALUES ("AAA", "BBB", "CCC")`,
	)
}

type testInsertGoStructWithTags struct {
	A string `spanner:"ColumnA"`
	B string `spanner:"ColumnB"`
	C string
}

func TestInsertWithGoStructWithTags(t *testing.T) {
	testInsert(t,
		memeduck.Insert("hoge", []string{"ColumnA", "ColumnB", "C"}).Values([]testInsertGoStructWithTags{
			testInsertGoStructWithTags{A: "AAA", B: "BBB", C: "CCC"},
		}),
		`INSERT INTO hoge (ColumnA, ColumnB, C) VALUES ("AAA", "BBB", "CCC")`,
	)
	testInsert(t,
		memeduck.Insert("hoge", []string{"ColumnA", "C"}).Values([]testInsertGoStructWithTags{
			testInsertGoStructWithTags{A: "AAA", B: "BBB", C: "CCC"},
		}),
		`INSERT INTO hoge (ColumnA, C) VALUES ("AAA", "CCC")`,
	)
	testInsert(t,
		memeduck.Insert("hoge", []string{"ColumnA", "ColumnB", "C"}).Values([]*testInsertGoStructWithTags{
			&testInsertGoStructWithTags{A: "AAA", B: "BBB", C: "CCC"},
		}),
		`INSERT INTO hoge (ColumnA, ColumnB, C) VALUES ("AAA", "BBB", "CCC")`,
	)
}

func TestInsertWithHeteroSlice(t *testing.T) {
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b", "c", "d"}).Values([][]interface{}{
			{int64(123), "45", []byte{6}, nil},
		}),
		`INSERT INTO hoge (a, b, c, d) VALUES (123, "45", B"\x06", NULL)`,
	)
	testInsert(t,
		memeduck.Insert("hoge", []string{"a", "b", "c", "d"}).Values([]interface{}{
			[]interface{}{int64(123), "45", []byte{6}, nil},
			[]int{1, 2, 3, 4},
		}),
		`INSERT INTO hoge (a, b, c, d) VALUES (123, "45", B"\x06", NULL), (1, 2, 3, 4)`,
	)
}

func parseTime(t *testing.T, s string) time.Time {
	ts, err := time.Parse(time.RFC3339Nano, s)
	assert.Nil(t, err, "failed to parse %s", s)
	return ts
}

func parseDate(t *testing.T, s string) civil.Date {
	d, err := civil.ParseDate(s)
	assert.Nil(t, err, "failed to parse %s", s)
	return d
}
