package memeduck_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/genkami/memeduck"
)

func testSelect(t *testing.T, stmt *memeduck.SelectStmt, expected string) {
	actual, err := stmt.SQL()
	assert.Nil(t, err, expected)
	assert.Equal(t, expected, actual)
}

func TestSelect(t *testing.T) {
	testSelect(t,
		memeduck.Select("hoge", []string{"a", "b"}),
		`SELECT a, b FROM hoge`,
	)
	testSelect(t,
		memeduck.Select("hoge", []string{"a", "b"}).Where(
			memeduck.Bool(true),
		),
		`SELECT a, b FROM hoge WHERE TRUE`,
	)
	testSelect(t,
		memeduck.Select("hoge", []string{"a", "b"}).Where(
			memeduck.Eq(1, 2),
		),
		`SELECT a, b FROM hoge WHERE 1 = 2`,
	)
	testSelect(t,
		memeduck.Select("hoge", []string{"a", "b"}).Where(
			memeduck.Ne(memeduck.Ident("a"), "foo"),
		),
		`SELECT a, b FROM hoge WHERE a != "foo"`,
	)
	testSelect(t,
		memeduck.Select("hoge", []string{"a", "b"}).Where(
			memeduck.Eq(memeduck.Ident("a"), 1),
			memeduck.Eq(memeduck.Ident("b"), "2"),
		),
		`SELECT a, b FROM hoge WHERE a = 1 AND b = "2"`,
	)
	testSelect(t,
		memeduck.Select("hoge", []string{"a", "b"}).Where(
			memeduck.Eq(memeduck.Ident("a"), 1),
			memeduck.Eq(memeduck.Ident("b"), "2"),
			memeduck.Ne(memeduck.Ident("c"), []byte{3}),
		),
		`SELECT a, b FROM hoge WHERE a = 1 AND b = "2" AND c != B"\x03"`,
	)
}

func TestSelectWithoutColumn(t *testing.T) {
	_, err := memeduck.Select("hoge", []string{}).SQL()
	assert.Error(t, err)
}
