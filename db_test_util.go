package main

import (
	"database/sql"
)

type DbTestUtil struct {
	db *sql.DB
}

func (dbTest *DbTestUtil) Init(db *sql.DB) {
	dbTest.db = db
}

func (dbTest *DbTestUtil) ClearTables(tables ...string) {
	for _, table := range tables {
		_, err := dbTest.db.Exec("DELETE FROM " + table)
		if err != nil {
			panic(err)
		}
	}
}

func (dbTest *DbTestUtil) AssertRowCount(table string, where string, params []interface{}, expectedCount int) {
	query := "SELECT COUNT(*) FROM " + table
	if where != "" {
		query += " WHERE " + where
	}

	rows, err := dbTest.db.Query(query, params...)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	if rows.Next() {
		var count int
		rows.Scan(&count)
		if !(count == expectedCount) {
			panic("Assertion failed. Expected count: " + string(expectedCount) + ", Actual count: " + string(count))
		}
	}
}
