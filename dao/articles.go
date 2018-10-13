package dao

import (
	"errors"
	"strings"
	"time"

	"github.com/golang/glog"
)

type Article struct {
	dbRecord
	Id_           int
	ApiSourceName string
	SourceId      string
	SourceName    string
	Author        string
	Title         string
	Description   string
	Url           string
	UrlToImage    string
	PublishedAt   time.Time
}

func (main *Dao) GetArticles() ([]*Article, error) {
	prefix := "main.Dao.GetArticles"
	query := "SELECT _id, api_source_name, source_id, source_name, author, title, description, url, url_to_image, published_at, created_at, updated_at, status FROM articles"
	rows, err := main.db.Query(query)
	if err != nil {
		return nil, errors.New(prefix + " (Query): " + err.Error())
	}

	articles := make([]*Article, 0)
	for rows.Next() {
		var article Article
		err := rows.Scan(&article.Id_, &article.ApiSourceName, &article.SourceId, &article.SourceName, &article.Author, &article.Title, &article.Description, &article.Url, &article.UrlToImage, &article.PublishedAt, &article.CreatedAt, &article.UpdatedAt, &article.Status)
		if err != nil {
			glog.Errorln(err)
		}
		articles = append(articles, &article)
	}

	return articles, nil
}

func (main *Dao) SaveArticles(articles []*Article) (int64, error) {
	query := "INSERT INTO articles (api_source_name, source_id, source_name, author, title, description, url, url_to_image, published_at) VALUES "
	parameterHolders := make([]string, len(articles))
	parameters := make([]interface{}, 0)

	for index, article := range articles {
		parameterHolders[index] = "(?, ?, ?, ?, ?, ?, ?, ?, ?)"
		parameters = append(parameters, article.ApiSourceName, article.SourceId, article.SourceName, article.Author, article.Title, article.Description, article.Url, article.UrlToImage, article.PublishedAt)
	}

	query = query + strings.Join(parameterHolders, ",")
	return main.prepareAndExec(query, parameters...)
}

func (main *Dao) prepareAndExec(query string, parameters ...interface{}) (int64, error) {
	prefix := "main.Dao.prepareAndExec"
	stmt, err := main.db.Prepare(query)
	if err != nil {
		return -1, errors.New(prefix + " (Prepare): " + err.Error())
	}

	res, err := stmt.Exec(parameters...)
	if err != nil {
		return -1, errors.New(prefix + " (Exec): " + err.Error())
	}

	numInserted, err := res.RowsAffected()
	if err != nil {
		return -1, errors.New(prefix + " (Rows Affected): " + err.Error())
	}

	return numInserted, nil
}
