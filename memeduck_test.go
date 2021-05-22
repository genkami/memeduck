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

func TestInsertWithIntSlice(t *testing.T) {
	testInsert(t,
		Insert("hoge", []string{"a", "b"}, [][]int{
			{123, 456},
		}),
		"INSERT INTO hoge (a, b) VALUES (123, 456)",
	)
}

func TestInsertWithIntPtrSlice(t *testing.T) {
	var a = int(123)
	var b = int(456)
	testInsert(t,
		Insert("hoge", []string{"a", "b"}, [][]*int{
			{&a, &b},
		}),
		"INSERT INTO hoge (a, b) VALUES (123, 456)",
	)
	testInsert(t,
		Insert("hoge", []string{"a", "b"}, [][]*int{
			{nil, nil},
		}),
		"INSERT INTO hoge (a, b) VALUES (NULL, NULL)",
	)
}

func TestInsertWithInt64Slice(t *testing.T) {
	testInsert(t,
		Insert("hoge", []string{"a", "b"}, [][]int64{
			{123, 456},
		}),
		"INSERT INTO hoge (a, b) VALUES (123, 456)",
	)
}

func TsetInsertWithInt64PtrSlice(t *testing.T) {
	var a = int64(123)
	var b = int64(456)
	testInsert(t,
		Insert("hoge", []string{"a", "b"}, [][]*int64{
			{&a, &b},
		}),
		"INSERT INTO hoge (a, b) VALUES (123, 456)",
	)
	testInsert(t,
		Insert("hoge", []string{"a", "b"}, [][]*int64{
			{nil, nil},
		}),
		"INSERT INTO hoge (a, b) VALUES (NULL, NULLtes)",
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
		"INSERT INTO hoge (a, b) VALUES (123, 456)",
	)
	testInsert(t,
		Insert("hoge", []string{"a", "b"}, [][]spanner.NullInt64{
			{null, null},
		}),
		"INSERT INTO hoge (a, b) VALUES (NULL, NULL)",
	)
}
