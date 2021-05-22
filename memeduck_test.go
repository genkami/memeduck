package memeduck

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInsert(t *testing.T) {
	test := func(is *InsertStmt, expected string) {
		actual, err := is.SQL()
		assert.Nil(t, err, expected)
		assert.Equal(t, expected, actual)
	}

	func() {
		test(
			Insert("person", []string{"age", "height"}, [][]int{
				{1600, 143},
			}),
			"INSERT INTO person (age, height) VALUES (1600, 143)",
		)
	}()
	func() {
		var age int = 1600
		var height int = 143
		test(
			Insert("person", []string{"age", "height"}, [][]*int{
				{&age, &height},
			}),
			"INSERT INTO person (age, height) VALUES (1600, 143)",
		)
	}()
	func() {
		test(
			Insert("person", []string{"age", "height"}, [][]int64{
				{1600, 143},
			}),
			"INSERT INTO person (age, height) VALUES (1600, 143)",
		)
	}()
	func() {
		var age int64 = 1600
		var height int64 = 143
		test(
			Insert("person", []string{"age", "height"}, [][]*int64{
				{&age, &height},
			}),
			"INSERT INTO person (age, height) VALUES (1600, 143)",
		)
	}()
}
