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

	test(
		Insert("person", []string{"age", "height"}, [][]int{
			[]int{1600, 143},
		}),
		"INSERT INTO person (age, height) VALUES (1600, 143)",
	)
}
