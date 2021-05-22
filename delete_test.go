package memeduck_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/genkami/memeduck"
	"github.com/genkami/memeduck/where"
)

func testDelete(t *testing.T, stmt *memeduck.DeleteStmt, expected string) {
	actual, err := stmt.SQL()
	assert.Nil(t, err, expected)
	assert.Equal(t, expected, actual)
}

func TestDeleteWithBool(t *testing.T) {
	testDelete(t,
		memeduck.Delete("hoge", where.Bool(true)),
		`DELETE FROM hoge WHERE TRUE`,
	)
}

func TestDeleteWithBinaryOp(t *testing.T) {
	testDelete(t,
		memeduck.Delete("hoge", where.Op(1, where.EQ, 2)),
		`DELETE FROM hoge WHERE 1 = 2`,
	)
}
