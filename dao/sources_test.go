package dao

import (
	"testing"
	"time"

	"github.com/golang/glog"
	"github.com/jforcode/Util"
)

func areArticlesEqual(article1 *Article, article2 *Article) bool {
	return article1.SourceId == article2.SourceId &&
		article1.SourceName == article2.SourceName &&
		article1.Author == article2.Author &&
		article1.Title == article2.Title
}

func TestArticles(t *testing.T) {
	db, dbMain := dbTestInit()
	defer db.Close()

	err := util.Db.ClearTables(db, "articles")
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

		rowCount, err := util.Db.GetRowCount(db, "articles", "", []interface{}{})
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
