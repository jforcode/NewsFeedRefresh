package main

import (
	"database/sql"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
	"github.com/jforcode/DbUtil"
)

func DbTestInit() (*sql.DB, *DbMain) {
	params := make(map[string]string)
	params["parseTime"] = "true"

	db, err := dbUtil.GetDb("root", "FORGIVEFeb@2018", "(127.0.0.1:3306)", "news_feed_test", params)
	if err != nil {
		panic(err)
	}

	dbMain := &DbMain{}
	dbMain.Init(db)

	return db, dbMain
}

func TestFlags(t *testing.T) {
	db, dbMain := DbTestInit()
	defer db.Close()

	dbUtil.ClearTables(db, "news_api_flags")

	t.Run("get non-existent flag", func(t *testing.T) {
		defer dbUtil.ClearTables(db, "news_api_flags")
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
				defer dbUtil.ClearTables(db, "news_api_flags")

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
				defer dbUtil.ClearTables(db, "news_api_flags")

				dbMain.SetFlag(test.key, test.value, test.typeTo)

				err := dbMain.SetFlag(test.key, test.updatedValue, test.typeTo)
				if !(err == nil) {
					t.FailNow()
				}

				flag, err := dbMain.GetFlag(test.key, test.typeTo)
				if !(err == nil && flag != nil && test.expected == flag.Value) {
					t.FailNow()
				}

				dbUtil.GetRowCount(db, "news_api_flags", "flag_key = ?", []interface{}{test.key})
			})
		}
	})

	// TODO: type mismatch in flags. conversion errors. setting in different type, fetching different type
}

func areSourcesEqual(source1 *Source, source2 *Source) bool {
	return source1.Name == source2.Name && source1.SourceId == source2.SourceId
}

func TestSources(t *testing.T) {
	db, dbMain := DbTestInit()
	defer db.Close()

	dbUtil.ClearTables(db, "sources")

	t.Run("save and get sources", func(t *testing.T) {
		sources := make([]*Source, 2)
		sources[0] = &Source{SourceId: "test", Name: "testSource"}
		sources[1] = &Source{SourceId: "test2", Name: "testSource2"}

		numSaved, err := dbMain.SaveSources(sources)
		if !(err == nil && numSaved == 2) {
			glog.Errorln(err)
			t.FailNow()
		}

		rowCount, err := dbUtil.GetRowCount(db, "sources", "", []interface{}{})
		if err != nil || rowCount != 2 {
			t.FailNow()
		}

		gotSources, err := dbMain.GetSources()
		if !(err == nil &&
			len(gotSources) == 2 &&
			areSourcesEqual(gotSources[0], sources[0]) &&
			areSourcesEqual(gotSources[1], sources[1])) {

			t.FailNow()
		}
	})
}

func areArticlesEqual(article1 *Article, article2 *Article) bool {
	return article1.SourceId == article2.SourceId &&
		article1.SourceName == article2.SourceName &&
		article1.Author == article2.Author &&
		article1.Title == article2.Title
}

func TestArticles(t *testing.T) {
	db, dbMain := DbTestInit()
	defer db.Close()

	err := dbUtil.ClearTables(db, "articles")
	if err != nil {
		t.FailNow()
	}

	t.Run("save and get articles", func(t *testing.T) {
		articles := make([]*Article, 2)
		articles[0] = &Article{SourceId: "test", SourceName: "testSource", Author: "", Title: "Test artcile 1", PublishedAt: time.Now().UTC()}
		articles[1] = &Article{SourceId: "test", SourceName: "testSource", Author: "Jeevan", Title: "Test article 2", PublishedAt: time.Now().UTC()}

		numSaved, err := dbMain.SaveArticles(articles)
		if !(err == nil && numSaved == 2) {
			glog.Errorln(err)
			t.FailNow()
		}

		rowCount, err := dbUtil.GetRowCount(db, "articles", "", []interface{}{})
		if !(err == nil && rowCount == 2) {
			t.FailNow()
		}

		gotArticles, err := dbMain.GetArticles()
		if !(err == nil &&
			len(gotArticles) == 2 &&
			areArticlesEqual(gotArticles[0], articles[0]) &&
			areArticlesEqual(gotArticles[1], articles[1])) {

			glog.Errorf("\nDb 0: %+v\nAr 0: %+v\nDb 1: %+v\nAr 1: %+v", gotArticles[0], articles[0], gotArticles[1], articles[1])
			t.FailNow()
		}
	})
}
