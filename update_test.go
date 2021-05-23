package memeduck_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/genkami/memeduck"
)

func testUpdate(t *testing.T, stmt *memeduck.UpdateStmt, expected string) {
	actual, err := stmt.SQL()
	assert.Nil(t, err, expected)
	assert.Equal(t, expected, actual)
}

func TestUpdate(t *testing.T) {
	testUpdate(t,
		memeduck.Update("hoge").
			Set(memeduck.Ident("a"), 1).
			Where(
				memeduck.Bool(true),
			),
		`UPDATE hoge SET a = 1 WHERE TRUE`,
	)
	testUpdate(t,
		memeduck.Update("hoge").
			Set(memeduck.Ident("a"), 1).
			Set(memeduck.Ident("b"), "foo").
			Where(
				memeduck.Bool(true),
			),
		`UPDATE hoge SET a = 1, b = "foo" WHERE TRUE`,
	)
	testUpdate(t,
		memeduck.Update("hoge").
			Set(memeduck.Ident("a"), 1).
			Set(memeduck.Ident("b"), "foo").
			Where(
				memeduck.Eq(memeduck.Ident("c"), "bar"),
			),
		`UPDATE hoge SET a = 1, b = "foo" WHERE c = "bar"`,
	)
	testUpdate(t,
		memeduck.Update("hoge").
			Set(memeduck.Ident("a"), memeduck.Ident("b")).
			Set(memeduck.Ident("b"), memeduck.Ident("a")).
			Where(
				memeduck.Eq(memeduck.Ident("c"), "bar"),
			),
		`UPDATE hoge SET a = b, b = a WHERE c = "bar"`,
	)
}

func TestUpdateWithNoSet(t *testing.T) {
	_, err := memeduck.Update("hoge").
		Where(
			memeduck.Eq(memeduck.Ident("a"), 1),
		).SQL()
	assert.Error(t, err, "UPDATE without SET clause")
}

func TestUpdateWithNoWhere(t *testing.T) {
	_, err := memeduck.Update("hoge").
		Set(memeduck.Ident("a"), 1).
		SQL()
	assert.Error(t, err, "UPDATE without WHERE clause")
}
