package dao

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jforcode/Util"
)

type Dao struct {
	db *sql.DB
}

func New(db *sql.DB) *Dao {
	return &Dao{db}
}

type dbRecord struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	Status    string
}

func dbTestInit() (*sql.DB, *Dao) {
	params := make(map[string]string)
	params["parseTime"] = "true"

	db, err := util.db.GetDb("root", "FORGIVEFeb@2018", "(127.0.0.1:3306)", "news_feed_test", params)
	if err != nil {
		panic(err)
	}

	dao := New(db)

	return db, dao
}
