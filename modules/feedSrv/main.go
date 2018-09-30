package feedSrv

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/jforcode/NewsFeedRefresh/modules/common"
)

type Main struct {
	db *sql.DB
}

func Init(db *sql.DB) (*Main, error) {
	if db == nil {
		return nil, errors.New("Invalid parameters")
	}

	main := Main{db: db}
	return &main, nil
}

func (main *Main) SaveSources(sources [](*common.Source)) (int64, error) {
	query := "INSERT INTO source (api_source_name, s_id, name, description, url, category, language, country) VALUES "
	parameterHolders := make([]string, len(sources))
	parameters := make([]interface{}, 0)

	for index, source := range sources {
		parameterHolders[index] = "(?, ?, ?, ?, ?, ?, ?, ?)"
		parameters = append(parameters, source.ApiSourceName, source.SourceId, source.Name, source.Description, source.Url, source.Category, source.Language, source.Country)
	}

	query = query + strings.Join(parameterHolders, ",")
	return main.batchInsert(query, parameters)
}

func (main *Main) SaveArticles(articles [](*common.Article)) (int64, error) {
	query := "INSERT INTO article (api_source_name, source_id, source_name, author, title, description, url, url_to_image, published_at) VALUES "
	parameterHolders := make([]string, len(articles))
	parameters := make([]interface{}, 0)

	for index, article := range articles {
		parameterHolders[index] = "(?, ?, ?, ?, ?, ?, ?, ?, ?)"
		parameters = append(parameters, article.ApiSourceName, article.SourceId, article.SourceName, article.Author, article.Title, article.Description, article.Url, article.UrlToImage, article.PublishedAt)
	}

	query = query + strings.Join(parameterHolders, ",")
	return main.batchInsert(query, parameters)
}

func (main *Main) batchInsert(query string, parameters ...interface{}) (int64, error) {
	stmt, err := main.db.Prepare(query)
	if err != nil {
		return -1, err
	}

	res, err := stmt.Exec(parameters)
	if err != nil {
		return -1, err
	}

	numInserted, err := res.RowsAffected()
	if err != nil {
		return -1, err
	}

	return numInserted, err
}
