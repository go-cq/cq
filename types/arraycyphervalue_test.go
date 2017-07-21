package types_test

import (
	"errors"

	_ "github.com/johnnadratowski/cq"
	"github.com/johnnadratowski/cq/types"
	. "gopkg.in/check.v1"
)

func (s *TypesSuite) TestQueryCypherValueArray(c *C) {
	rows := prepareAndQuery("return [1.1,2.1,'asdf']")
	rows.Next()
	var test types.ArrayCypherValue
	err := rows.Scan(&test)
	c.Assert(err, IsNil)
	c.Assert(test.Val, DeepEquals,
		[]types.CypherValue{
			types.CypherValue{types.CypherFloat64, 1.1},
			types.CypherValue{types.CypherFloat64, 2.1},
			types.CypherValue{types.CypherString, "asdf"}})
}

func (s *TypesSuite) TestQueryNullCypherValueArray(c *C) {
	rows := prepareAndQuery("return null")
	rows.Next()
	var test types.ArrayCypherValue
	err := rows.Scan(&test)
	c.Assert(err, DeepEquals, errors.New("sql: Scan error on column index 0: cq: scan value is null"))
}
