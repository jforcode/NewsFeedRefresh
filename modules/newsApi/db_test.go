package newsApi

import (
	"database/sql"
	"fmt"
	"testing"
)

func freshDb() *sql.DB {
	db, err := sql.Open("mysql", "testing data source")
	if err != nil {
		fmt.Println("Error in initializing db: " + err.Error())
	}

	return db
}

// TODO: mock a db, and test with that
// or test it with a test db
func TestGetFlags(t *testing.T) {
	t.Run("get flags", func(t *testing.T) {
		t.Run("called with empty key", func(t *testing.T) {
			// should return nil
		})

		t.Run("existent flag key", func(t *testing.T) {
			// should return that value
		})

		t.Run("non existent flag key", func(t *testing.T) {
			// should return nil
		})
	})
}

func TestSetFlags(t *testing.T) {
	t.Run("set flags", func(t *testing.T) {
		t.Run("called with empty key", func(t *testing.T) {

		})

		t.Run("flag key is too long", func(t *testing.T) {

		})

		t.Run("non existent flag key", func(t *testing.T) {
			// should create
		})

		t.Run("existent flag key", func(t *testing.T) {
			// should update the value
		})
	})
}
