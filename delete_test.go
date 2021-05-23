package memeduck_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/genkami/memeduck"
)

func testDelete(t *testing.T, stmt *memeduck.DeleteStmt, expected string) {
	actual, err := stmt.SQL()
	assert.Nil(t, err, expected)
	assert.Equal(t, expected, actual)
}

func TestDelete(t *testing.T) {
	testDelete(t,
		memeduck.Delete("hoge").Where(
			memeduck.Bool(true),
		),
		`DELETE FROM hoge WHERE TRUE`,
	)
	testDelete(t,
		memeduck.Delete("hoge").Where(
			memeduck.Eq(1, 2),
		),
		`DELETE FROM hoge WHERE 1 = 2`,
	)
	testDelete(t,
		memeduck.Delete("hoge").Where(
			memeduck.Ne(memeduck.Ident("a"), "foo"),
		),
		`DELETE FROM hoge WHERE a != "foo"`,
	)
	testDelete(t,
		memeduck.Delete("hoge").Where(
			memeduck.Eq(memeduck.Ident("a"), 1),
			memeduck.Eq(memeduck.Ident("b"), "2"),
		),
		`DELETE FROM hoge WHERE a = 1 AND b = "2"`,
	)
	testDelete(t,
		memeduck.Delete("hoge").Where(
			memeduck.Eq(memeduck.Ident("a"), 1),
			memeduck.Eq(memeduck.Ident("b"), "2"),
			memeduck.Ne(memeduck.Ident("c"), []byte{3}),
		),
		`DELETE FROM hoge WHERE a = 1 AND b = "2" AND c != B"\x03"`,
	)
}

func TestDeleteWithMultipleWhereClause(t *testing.T) {
	testDelete(t,
		memeduck.Delete("hoge").Where(
			memeduck.Eq(memeduck.Ident("a"), 1),
		).Where(
			memeduck.Eq(memeduck.Ident("b"), "2"),
			memeduck.Ne(memeduck.Ident("c"), []byte{3}),
		),
		`DELETE FROM hoge WHERE a = 1 AND b = "2" AND c != B"\x03"`,
	)
}

func TestDeleteWithNoWhereClause(t *testing.T) {
	_, err := memeduck.Delete("hoge").SQL()
	assert.Error(t, err)
}
