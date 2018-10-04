package main

import (
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func GetTestDb() (*DbMain, error) {
	db, err := sql.Open("mysql", "root:FORGIVEFeb@2018@tcp(127.0.0.1:3306)/news_feed_test")
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	ret := &DbMain{}
	ret.Init(db)

	return ret, nil
}

func clearTables(db *sql.DB, tables ...string) {
	for _, table := range tables {
		_, err := db.Exec("DELETE FROM " + table)
		if err != nil {
			panic(err)
		}
	}
}

func assertRowCount(db *sql.DB, table string, where string, params []interface{}, expectedCount int) {
	query := "SELECT COUNT(*) FROM " + table + " WHERE " + where
	rows, err := db.Query(query, params...)
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

func TestFlags(t *testing.T) {
	dbMain, err := GetTestDb()
	if err != nil {
		panic(err)
	}

	clearTables(dbMain.db, "news_api_flags")

	t.Run("get non-existent flag", func(t *testing.T) {
		defer clearTables(dbMain.db, "news_api_flags")
		flag, err := dbMain.GetFlag("test", "int")

		if !(err == nil && flag == nil) {
			t.FailNow()
		}
	})

	t.Run("set and get non-existent flag", func(t *testing.T) {
		tests := []struct {
			name     string
			key      string
			value    string
			expected interface{}
			typeTo   string
		}{
			{"string", "test", "asdf", "asdf", "string"},
			{"int", "test", "1234", 1234, "int"},
			{"boolean", "test", "TRUE", true, "bool"},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				defer clearTables(dbMain.db, "news_api_flags")

				err := dbMain.SetFlag(test.key, test.value, test.typeTo)
				if !(err == nil) {
					t.FailNow()
				}

				flag, err := dbMain.GetFlag(test.key, test.typeTo)
				if !(err == nil && flag != nil && test.expected == flag.Value) {
					t.FailNow()
				}
			})
		}
	})

	t.Run("update flag", func(t *testing.T) {
		tests := []struct {
			name         string
			key          string
			value        string
			updatedValue string
			expected     interface{}
			typeTo       string
		}{
			{"string", "test", "asdf", "def", "def", "string"},
			{"int", "test", "1234", "5678", 5678, "int"},
			{"boolean", "test", "TRUE", "FALSE", false, "bool"},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				defer clearTables(dbMain.db, "news_api_flags")

				dbMain.SetFlag(test.key, test.value, test.typeTo)

				err := dbMain.SetFlag(test.key, test.updatedValue, test.typeTo)
				if !(err == nil) {
					t.FailNow()
				}

				flag, err := dbMain.GetFlag(test.key, test.typeTo)
				if !(err == nil && flag != nil && test.expected == flag.Value) {
					t.FailNow()
				}

				assertRowCount(dbMain.db, "news_api_flags", "flag_key = ?", []interface{}{test.key}, 1)
			})
		}
	})
}

// TODO: type mismatch in flags. conversion errors. setting in different type, fetching different type
