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
			memeduck.Ne(memeduck.Ident("a"), memeduck.Param("a")),
		),
		`SELECT a, b FROM hoge WHERE a != @a`,
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
	testSelect(t,
		memeduck.Select("hoge", []string{"a", "b"}).Where(
			memeduck.Eq(memeduck.Ident("a", "b"), 1),
			memeduck.Eq(memeduck.Ident("a", "c"), "2"),
		),
		`SELECT a, b FROM hoge WHERE a.b = 1 AND a.c = "2"`,
	)
}

func TestSelectWithAsStruct(t *testing.T) {
	testSelect(t,
		memeduck.Select("hoge", []string{"a", "b"}).Where(
			memeduck.Eq(memeduck.Ident("c"), memeduck.Param("id")),
		).AsStruct(),
		`SELECT AS STRUCT a, b FROM hoge WHERE c = @id`,
	)
}

func TestSelectWithOrderBy(t *testing.T) {
	testSelect(t,
		memeduck.Select("hoge", []string{"a", "b"}).Where(
			memeduck.Eq(memeduck.Ident("a"), 123),
		).OrderBy("a", memeduck.ASC),
		`SELECT a, b FROM hoge WHERE a = 123 ORDER BY a ASC`,
	)
	testSelect(t,
		memeduck.Select("hoge", []string{"a", "b"}).Where(
			memeduck.Eq(memeduck.Ident("a"), 123),
		).OrderBy("a", memeduck.ASC).
			OrderBy("b", memeduck.DESC),
		`SELECT a, b FROM hoge WHERE a = 123 ORDER BY a ASC, b DESC`,
	)
}

func TestSelectWithLimit(t *testing.T) {
	testSelect(t,
		memeduck.Select("hoge", []string{"a", "b"}).
			Limit(10),
		`SELECT a, b FROM hoge LIMIT 10`,
	)
	testSelect(t,
		memeduck.Select("hoge", []string{"a", "b"}).
			OrderBy("a", memeduck.ASC).
			Limit(10),
		`SELECT a, b FROM hoge ORDER BY a ASC LIMIT 10`,
	)
	testSelect(t,
		memeduck.Select("hoge", []string{"a", "b"}).Where(
			memeduck.Eq(memeduck.Ident("a"), 123),
		).OrderBy("a", memeduck.ASC).
			Limit(10),
		`SELECT a, b FROM hoge WHERE a = 123 ORDER BY a ASC LIMIT 10`,
	)
}

func TestSelectWithLimitOffset(t *testing.T) {
	testSelect(t,
		memeduck.Select("hoge", []string{"a", "b"}).
			LimitOffset(10, 3),
		`SELECT a, b FROM hoge LIMIT 10 OFFSET 3`,
	)
	testSelect(t,
		memeduck.Select("hoge", []string{"a", "b"}).
			OrderBy("a", memeduck.ASC).
			LimitOffset(10, 3),
		`SELECT a, b FROM hoge ORDER BY a ASC LIMIT 10 OFFSET 3`,
	)
	testSelect(t,
		memeduck.Select("hoge", []string{"a", "b"}).Where(
			memeduck.Eq(memeduck.Ident("a"), 123),
		).OrderBy("a", memeduck.ASC).
			LimitOffset(10, 3),
		`SELECT a, b FROM hoge WHERE a = 123 ORDER BY a ASC LIMIT 10 OFFSET 3`,
	)
}

func TestSelectWithoutColumn(t *testing.T) {
	_, err := memeduck.Select("hoge", []string{}).SQL()
	assert.Error(t, err)
}

func TestSelectWithSubQuery(t *testing.T) {
	testSelect(t,
		memeduck.Select("hoge", []string{"a", "b"}).
			SubQuery(
				memeduck.ScalarSubQuery(memeduck.Select("fuga", []string{"c"}).Where(memeduck.Eq(3, 4))).As("fuga"),
			).
			Where(
				memeduck.Eq(1, 2),
			),
		`SELECT a, b, (SELECT c FROM fuga WHERE 3 = 4) AS fuga FROM hoge WHERE 1 = 2`,
	)

	testSelect(t,
		memeduck.Select("hoge", []string{"a", "b"}).
			SubQuery(
				memeduck.ArraySubQuery(memeduck.Select("fuga", []string{"c", "d"}).Where(memeduck.Eq(3, 4)).AsStruct()).As("fuga"),
			).
			Where(
				memeduck.Eq(1, 2),
			),
		`SELECT a, b, ARRAY(SELECT AS STRUCT c, d FROM fuga WHERE 3 = 4) AS fuga FROM hoge WHERE 1 = 2`,
	)
}
