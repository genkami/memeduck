package memeduck

import (
	"testing"

	"cloud.google.com/go/spanner"
	"github.com/stretchr/testify/assert"
)

func testInsert(t *testing.T, stmt *InsertStmt, expected string) {
	actual, err := stmt.SQL()
	assert.Nil(t, err, expected)
	assert.Equal(t, expected, actual)
}

func TestInsertWithStringSlice(t *testing.T) {
	testInsert(t,
		Insert("hoge", []string{"a", "b"}, [][]string{
			{"foo", "bar"},
		}),
		`INSERT INTO hoge (a, b) VALUES ("foo", "bar")`,
	)
}

func TestInsertWithStringPtrSlice(t *testing.T) {
	var a = "foo"
	var b = "bar"
	testInsert(t,
		Insert("hoge", []string{"a", "b"}, [][]*string{
			{&a, &b},
		}),
		`INSERT INTO hoge (a, b) VALUES ("foo", "bar")`,
	)
	testInsert(t,
		Insert("hoge", []string{"a", "b"}, [][]*string{
			{nil, nil},
		}),
		`INSERT INTO hoge (a, b) VALUES (NULL, NULL)`,
	)
}

func TestInsertWithNullStringSlice(t *testing.T) {
	var a = spanner.NullString{StringVal: "foo", Valid: true}
	var b = spanner.NullString{StringVal: "bar", Valid: true}
	testInsert(t,
		Insert("hoge", []string{"a", "b"}, [][]spanner.NullString{
			{a, b},
		}),
		`INSERT INTO hoge (a, b) VALUES ("foo", "bar")`,
	)
	var null = spanner.NullString{}
	testInsert(t,
		Insert("hoge", []string{"a", "b"}, [][]spanner.NullString{
			{null, null},
		}),
		`INSERT INTO hoge (a, b) VALUES (NULL, NULL)`,
	)
}

func TestInsertWithByteSliceSlice(t *testing.T) {
	testInsert(t,
		Insert("hoge", []string{"a", "b"}, [][][]byte{
			{{0, 1}, {2, 3, 4}},
		}),
		`INSERT INTO hoge (a, b) VALUES (B"\x00\x01", B"\x02\x03\x04")`,
	)
	testInsert(t,
		Insert("hoge", []string{"a", "b"}, [][][]byte{
			{nil, nil},
		}),
		`INSERT INTO hoge (a, b) VALUES (NULL, NULL)`,
	)
}

func TestInsertWithIntSlice(t *testing.T) {
	testInsert(t,
		Insert("hoge", []string{"a", "b"}, [][]int{
			{123, 456},
		}),
		`INSERT INTO hoge (a, b) VALUES (123, 456)`,
	)
}

func TestInsertWithIntPtrSlice(t *testing.T) {
	var a = int(123)
	var b = int(456)
	testInsert(t,
		Insert("hoge", []string{"a", "b"}, [][]*int{
			{&a, &b},
		}),
		`INSERT INTO hoge (a, b) VALUES (123, 456)`,
	)
	testInsert(t,
		Insert("hoge", []string{"a", "b"}, [][]*int{
			{nil, nil},
		}),
		`INSERT INTO hoge (a, b) VALUES (NULL, NULL)`,
	)
}

func TestInsertWithInt64Slice(t *testing.T) {
	testInsert(t,
		Insert("hoge", []string{"a", "b"}, [][]int64{
			{123, 456},
		}),
		`INSERT INTO hoge (a, b) VALUES (123, 456)`,
	)
}

func TsetInsertWithInt64PtrSlice(t *testing.T) {
	var a = int64(123)
	var b = int64(456)
	testInsert(t,
		Insert("hoge", []string{"a", "b"}, [][]*int64{
			{&a, &b},
		}),
		`INSERT INTO hoge (a, b) VALUES (123, 456)`,
	)
	testInsert(t,
		Insert("hoge", []string{"a", "b"}, [][]*int64{
			{nil, nil},
		}),
		`INSERT INTO hoge (a, b) VALUES (NULL, NULLtes)`,
	)
}

func TestInsertWithNullInt64Slice(t *testing.T) {
	var a = spanner.NullInt64{Int64: 123, Valid: true}
	var b = spanner.NullInt64{Int64: 456, Valid: true}
	var null = spanner.NullInt64{}
	testInsert(t,
		Insert("hoge", []string{"a", "b"}, [][]spanner.NullInt64{
			{a, b},
		}),
		`INSERT INTO hoge (a, b) VALUES (123, 456)`,
	)
	testInsert(t,
		Insert("hoge", []string{"a", "b"}, [][]spanner.NullInt64{
			{null, null},
		}),
		`INSERT INTO hoge (a, b) VALUES (NULL, NULL)`,
	)
}

func TestInsertWithBoolSlice(t *testing.T) {
	testInsert(t,
		Insert("hoge", []string{"a", "b"}, [][]bool{
			{true, false},
		}),
		`INSERT INTO hoge (a, b) VALUES (TRUE, FALSE)`,
	)
}

func TestInsertWithBoolPtrSlice(t *testing.T) {
	var a = true
	var b = false
	testInsert(t,
		Insert("hoge", []string{"a", "b"}, [][]*bool{
			{&a, &b},
		}),
		`INSERT INTO hoge (a, b) VALUES (TRUE, FALSE)`,
	)
	testInsert(t,
		Insert("hoge", []string{"a", "b"}, [][]*bool{
			{nil, nil},
		}),
		`INSERT INTO hoge (a, b) VALUES (NULL, NULL)`,
	)
}

func TestInsertWithNullBoolSlice(t *testing.T) {
	var a = spanner.NullBool{Bool: true, Valid: true}
	var b = spanner.NullBool{Bool: false, Valid: true}
	testInsert(t,
		Insert("hoge", []string{"a", "b"}, [][]spanner.NullBool{
			{a, b},
		}),
		`INSERT INTO hoge (a, b) VALUES (TRUE, FALSE)`,
	)
	var null = spanner.NullBool{}
	testInsert(t,
		Insert("hoge", []string{"a", "b"}, [][]spanner.NullBool{
			{null, null},
		}),
		`INSERT INTO hoge (a, b) VALUES (NULL, NULL)`,
	)
}
