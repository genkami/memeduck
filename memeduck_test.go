package memeduck

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInsert(t *testing.T) {
	test := func(ivb *InsertIntoValuesBuilder, expected string) {
		actual, err := ivb.SQL()
		assert.Nil(t, err, expected)
		assert.Equal(t, expected, actual)
	}

	test(
		InsertInto("person", []string{"age", "height"}).Values([][]int{
			[]int{1600, 143},
		}),
		"INSERT INTO person (age, height) VALUES (1600, 143)",
	)
}
