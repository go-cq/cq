package types_test

import (
	_ "github.com/johnnadratowski/cq"
	. "gopkg.in/check.v1"

	"github.com/johnnadratowski/cq/types"
)

func (s *TypesSuite) TestQueryRelationship(c *C) {
	stmt := prepareTest(`create (:Test)-[r:TEST_TYPE {foo:"bar", i:1}]->(:Test) return r`)
	rows, err := stmt.Query()
	c.Assert(err, IsNil)

	rows.Next()
	var test types.Relationship
	err = rows.Scan(&test)
	c.Assert(err, IsNil)
	t1 := types.Relationship{}
	t1.Properties = map[string]types.CypherValue{}
	t1.Properties["foo"] = types.CypherValue{Type: types.CypherString, Val: "bar"}
	t1.Properties["i"] = types.CypherValue{Type: types.CypherInt, Val: 1}
	c.Assert(test.Properties, DeepEquals, t1.Properties)
	c.Assert(test.Type, Equals, "TEST_TYPE")
}
