package cq

import (
	"database/sql/driver"
	"log"
	"testing"
)

func prepareTest(query string) driver.Stmt {
	db := openTest()
	if db == nil {
		log.Fatal("can't connect to test db: ", testURL)
	}
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Print(err)
	}
	return stmt
}

func TestQuerySimple(t *testing.T) {
	stmt := prepareTest("return 1")
	rows, err := stmt.Query([]driver.Value{})
	if err != nil {
		t.Fatal(err)
	}
	dest := make([]driver.Value, 1)

	err = rows.Next(dest)
	if err != nil {
		t.Fatal(err)
	}

	if rows.Columns()[0] != "1" {
		t.Fatal("column doesn't match")
	}

	err = rows.Next(dest)
	if err == nil {
		t.Fatal("doesn't end after first row")
	}
}
