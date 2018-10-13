package dao

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jforcode/Util"
)

type IDao interface {
	GetArticles() ([]*Article, error)
	SaveArticles(articles []*Article) (int64, error)

	GetSources() ([]*Source, error)
	SaveSources(sources []*Source) (int64, error)
	ClearSources(apiSourceName string) (int64, error)

	GetFlag(key string, typeTo FlagType) (*Flag, error)
	SetFlag(key, value string, typeTo FlagType) error
	updateFlag(key, value string) error
	createFlag(key, value string) error
}

type Dao struct {
	db *sql.DB
}

func New(db *sql.DB) IDao {
	return &Dao{db}
}

type dbRecord struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	Status    string
}

func dbTestInit() (*sql.DB, IDao) {
	params := make(map[string]string)
	params["parseTime"] = "true"

	db, err := util.Db.GetDb("root", "FORGIVEFeb@2018", "(127.0.0.1:3306)", "news_feed_test", params)
	if err != nil {
		panic(err)
	}

	dao := New(db)

	return db, dao
}
