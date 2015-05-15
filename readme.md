# cq - cypher queries for neo4j
cq.v2 is the neo4j 2.3 focused development branch. It will support packstream
and NDP (Neo4j Data Protocol) for the socket-based connections, as well as the old HTTP-based
transactional Cypher API. Also, in addition to supporting the stdlib
database/sql API, it will also support a slightly lower-level Neo4j-specific API (at least,
for the NDP packages).

If you'd like to use the new [gopkg.in](http://godoc.org/gopkg.in/docs.v1) semantic versioning system:

[![Build Status](https://travis-ci.org/go-cq/cq.svg?branch=master)](https://travis-ci.org/go-cq/cq?branch=v2)
[![Coverage Status](https://img.shields.io/coveralls/go-cq/cq.svg)](https://coveralls.io/r/go-cq/cq?branch=v2)
[![Waffle](https://badge.waffle.io/go-cq/cq.png?label=ready)](https://waffle.io/go-cq/cq)
[![Gitter chat](https://badges.gitter.im/go-cq/cq.png)](https://gitter.im/go-cq/cq)

# NDP API

```go
import "gopkg.in/cq.v2/ndp.v1"
```

The NDP API should be the way to go moving forward. It should be significantly
faster than the HTTP API. It will probably only be in "beta" for Neo4j 2.3.

## NDP API [minimum viable snippet](http://blog.fogus.me/2012/08/23/minimum-viable-snippet/)

```go
package main

import (
    "log"

    "gopkg.in/cq.v2/ndp.v1" 
)

type Person struct {
    name string
    age  int
}

func main() {
    session, err := ndp.NewSession("localhost", 7687)
    if err != nil {
        log.Fatal(err)
    }

    results, err := session.Run("MATCH (a:Person)-[:FOLLOWS]->(p:Person) where a.name = {0} RETURN a.name, a.age", "wefreema")
    if err != nil {
        log.Fatal(err)
    }
    for _, result := range results {
        person := Person{}
        err = result.Scan(&person)
        if err != nil {
            log.Fatal(err)
        }
        fmt.Printf("%s is %d years old\n", person.name, person.age)
    }
}
```

# database/sql API

The database/sql API can be used with NDP (neo4j 2.3+), or HTTP (neo4j 2.0+).

```go
import "gopkg.in/cq.v2/ndp.v1/stdlib"

// or, if you're running neo4j <= 2.2.x (or don't want to use NDP) 
import "gopkg.in/cq.v2/http/stdlib"
```

See the [excellent database/sql tutorial](http://go-database-sql.org/index.html) from [VividCortex](https://vividcortex.com/), as well as the [package documentation for database/sql](http://golang.org/pkg/database/sql/) for an introduction to the idiomatic go database access.

You can (and should) use parameters, but the placeholders must be numbers in sequence, e.g. `{0}`, `{1}`, `{2}`, and then you must put them in order in the calls to `Query`/`Exec`. If you'd like to use named parameters, you can use the [sqlx](https://github.com/jmoiron/sqlx) library along with cq. Please let me know if any issues arise from using sqlx with cq--it is not thoroughly tested.

## database/SQL API [minimum viable snippet](http://blog.fogus.me/2012/08/23/minimum-viable-snippet/)

```go
package main

import (
	"database/sql"
	"log"
	
	_ "gopkg.in/cq.v2/ndp.v1/stdlib"
)

func main() {
    db, err := sql.Open("cq", "localhost:7687;user;pass")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	stmt, err := db.Prepare(`
		match (n:Person)-[:FOLLOWS]->(m:Person) 
		where n.screenName = {0} 
		return m.name, m.age
		limit 10
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query("wefreema")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var name string
    var age int
	for rows.Next() {
		err := rows.Scan(&name, &age)
		if err != nil {
			log.Fatal(err)
		}
        person := Person{name, age}
		log.Println(friend)
	}
}
```

## deployment on Heroku w/ GrapheneDB ***needs updating for v2

There is a repo with a template app for Heroku [here](https://github.com/wfreeman/cq-example).
Use this Heroku deploy button to push the template project on to a new app on your Heroku account.

[![Deploy](https://www.herokucdn.com/deploy/button.png)](https://heroku.com/deploy?template=https://github.com/wfreeman/cq-example)

## transactional API
The transactional API using `db.Begin()` is optimized for sending many queries to the [transactional Cypher endpoint](http://docs.neo4j.org/chunked/milestone/rest-api-transactional.html), in that it will batch them up and send them in chunks by default. Currently only supports `stmt.Exec()` within a transaction, will work on supporting `stmt.Query()` next and queueing up results.

#### transactional API example
```go
func main() {
	db, err := sql.Open("cq-http", "http://localhost:7474;user;pass")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	
	stmt, err := tx.Prepare("create (:User {screenName:{0}})")	
	if err != nil {
		log.Fatal(err)
	}
	
	stmt.Exec("wefreema")
	stmt.Exec("JnBrymn")
	stmt.Exec("technige")
	
	err := tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}
```

## types subpackage

database/sql out of the box doesn't implement many types to pass in as parameters or Scan() out of rows. Custom Cypher types are implemented in the `cq/types` subpackage (`import "gopkg.in/cq.v2/http/types"`). These custom types allow users of cq to `Scan()` types out of results, as well as pass types in as parameters.

| Go type			| Can be <br/>query parameter?	| cq wrapper, for Scan	| CypherType uint8 |
|:------------------ |:------------------:|:--------------------- | --------------------- |
| `nil`						| yes						| `CypherValue`				| `CypherNull`						|
| `bool`						| yes						| use go `bool`				| `CypherBoolean`					|
| `string`					| yes						| use go `string`				| `CypherString`					|
| `int`						| yes						| use go `int`					| `CypherInt`					|
| `int64`					| yes						| use go `int64`				| `CypherInt64`					|
| `float64`					| yes						| use go `float64`			| `CypherFloat64`					|
| `time.Time`				| yes						| `NullTime`			| `NullTime`					|
| `types.Node`				| no						| `Node`							| `CypherNode`						|
| `types.Relationship`	| no						| `Relationship`				| `CypherRelationship`			|
| `types.CypherValue`	| yes						| `CypherValue`				| `CypherValueType`			|
| N/A							| no						| not implemented				| `CypherPath`						|
| `[]string`				| yes						| `ArrayString`				| `CypherArrayString` |
| `[]int`					| yes						| `ArrayInt`					| `CypherArrayInt` |
| `[]int64`					| yes						| `ArrayInt64`					| `CypherArrayInt64` |
| `[]float64`				| yes						| `ArrayFloat64`				| `CypherArrayFloat64`	|
| `[]types.CypherValue`	| yes						| `ArrayCypherValue`			| `CypherArrayCypherValue`	|
| `map[string]string`	| yes						| `MapStringString`			| `CypherMapStringString`			|
| `map[string]types.CypherValue`| yes			| `MapStringCypherValue`	| `CypherMapStringCypherValue`				|


## transactional API benchmarks
Able to get sustained times of 20k+ cypher statements per second, even with multiple nodes per create... on a 2011 vintage macbook.

```
(master ✓) wes-macbook:cq go test -bench=".*Transaction.*" -test.benchtime=10s
PASS
BenchmarkTransactional10SimpleCreate	  100000	    150630 ns/op
BenchmarkTransactional100SimpleCreate	  500000	     39202 ns/op
BenchmarkTransactional1000SimpleCreate	 1000000	     27320 ns/op
BenchmarkTransactional10000SimpleCreate	  500000	     28524 ns/op
ok  	github.com/wfreeman/cq	79.973s
```


## thanks
Thanks to issue reporters and [contributors](https://github.com/go-cq/cq/graphs/contributors)!

## license

MIT license. See license file.

