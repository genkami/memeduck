package memeduck

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func testInsert(t *testing.T, stmt *InsertStmt, expected string) {
	actual, err := stmt.SQL()
	assert.Nil(t, err, expected)
	assert.Equal(t, expected, actual)
}

func TestInsertWithIntSlice(t *testing.T) {
	testInsert(t,
		Insert("person", []string{"age", "height"}, [][]int{
			{1600, 143},
		}),
		"INSERT INTO person (age, height) VALUES (1600, 143)",
	)
}

func TestInsertWithIntPtrSlice(t *testing.T) {
	var age int = 1600
	var height int = 143
	testInsert(t,
		Insert("person", []string{"age", "height"}, [][]*int{
			{&age, &height},
		}),
		"INSERT INTO person (age, height) VALUES (1600, 143)",
	)
}

func TestInsertWithInt64Slice(t *testing.T) {
	testInsert(t,
		Insert("person", []string{"age", "height"}, [][]int64{
			{1600, 143},
		}),
		"INSERT INTO person (age, height) VALUES (1600, 143)",
	)
}
func TsetInsertWithInt64PtrSlice(t *testing.T) {
	var age int64 = 1600
	var height int64 = 143
	testInsert(t,
		Insert("person", []string{"age", "height"}, [][]*int64{
			{&age, &height},
		}),
		"INSERT INTO person (age, height) VALUES (1600, 143)",
	)
}
