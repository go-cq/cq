package types_test

import (
	"errors"

	_ "github.com/Unified/golang-lib/lib/neo/drivers/cq"
	. "gopkg.in/check.v1"

	"github.com/Unified/golang-lib/lib/neo/drivers/cq/types"
)

func (s *TypesSuite) TestQueryMapStringStringParam(c *C) {
	stmt := prepareTest("with {0} as test return test")
	rows, err := stmt.Query(
		types.MapStringString{Val: map[string]string{"key1": "1", "key2": "2"}})
	c.Assert(err, IsNil)

	rows.Next()
	var test types.MapStringString
	err = rows.Scan(&test)
	c.Assert(err, IsNil)
	c.Assert(test.Val, DeepEquals, map[string]string{"key1": "1", "key2": "2"})
}

func (s *TypesSuite) TestQueryStringStringMapParam(c *C) {
	stmt := prepareTest("with {0} as test return test")
	rows, err := stmt.Query(map[string]string{"key1": "1", "key2": "2"})
	c.Assert(err, IsNil)

	rows.Next()
	var test types.MapStringString
	err = rows.Scan(&test)
	c.Assert(err, IsNil)
	c.Assert(test.Val, DeepEquals, map[string]string{"key1": "1", "key2": "2"})
}

func (s *TypesSuite) TestQueryMapStringString(c *C) {
	rows := prepareAndQuery(`return {key1:"1",key2:"2"}`)
	rows.Next()
	var test types.MapStringString
	err := rows.Scan(&test)
	c.Assert(err, IsNil)
	c.Assert(test.Val, DeepEquals, map[string]string{"key1": "1", "key2": "2"})
}

func (s *TypesSuite) TestQueryBadMapStringString(c *C) {
	rows := prepareAndQuery(`return {key1:1}`)
	rows.Next()
	var test types.MapStringString
	err := rows.Scan(&test)
	c.Assert(err, DeepEquals, errors.New("sql: Scan error on column index 0: cq: invalid Scan value for *types.MapStringString: map[string]types.CypherValue"))
}

func (s *TypesSuite) TestQueryNullMapStringString(c *C) {
	rows := prepareAndQuery("return null")
	rows.Next()
	var test types.MapStringString
	err := rows.Scan(&test)
	c.Assert(err, DeepEquals, errors.New("sql: Scan error on column index 0: cq: scan value is null"))
}
