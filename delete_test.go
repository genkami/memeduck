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

func TestDeleteWithBool(t *testing.T) {
	testDelete(t,
		memeduck.Delete("hoge", memeduck.Bool(true)),
		`DELETE FROM hoge WHERE TRUE`,
	)
}

func TestDeleteWithBinaryOp(t *testing.T) {
	testDelete(t,
		memeduck.Delete("hoge", memeduck.Op(1, memeduck.EQ, 2)),
		`DELETE FROM hoge WHERE 1 = 2`,
	)
	testDelete(t,
		memeduck.Delete("hoge", memeduck.Op(memeduck.Ident("a"), memeduck.NE, "foo")),
		`DELETE FROM hoge WHERE a != "foo"`,
	)
}
