package main

import (
	"database/sql"
	"errors"
	"strings"
)

type DbMain struct {
	db *sql.DB
}

func (main *DbMain) Init(db *sql.DB) error {
	if db == nil {
		return errors.New("Invalid parameters")
	}

	main.db = db
	return nil
}

func (main *DbMain) GetFlag(key string) (*Flag, error) {
	query := "SELECT _id, flag_key, flag_value, created_at, updated_at, status FROM news_api_flags WHERE flag_key = ?"
	rows, err := main.db.Query(query, key)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		flag := &Flag{}
		rows.Scan(flag.Id_, flag.Key, flag.Value, flag.CreatedAt, flag.UpdatedAt, flag.Status)
		return flag, nil
	}

	return nil, nil
}

func (main *DbMain) SetFlag(key string, value string) error {
	flag, err := main.GetFlag(key)
	if err != nil {
		return err
	}
	if flag != nil {
		return main.updateFlag(key, value)
	} else {
		return main.createFlag(key, value)
	}
}

func (main *DbMain) createFlag(key, value string) error {
	query := "UPDATE news_api_flags SET flag_value = ? WHERE flag_key = ?"
	stmt, err := main.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(value, key)
	if err != nil {
		return err
	}

	return nil
}

func (main *DbMain) updateFlag(key, value string) error {
	query := "INSERT INTO news_api_flags (flag_key, flag_value) VALUES (?, ?)"
	stmt, err := main.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(key, value)
	if err != nil {
		return err
	}

	return nil
}

func (main *DbMain) GetSources() ([]*Source, error) {
	query := "SELECT _id, api_source_name, s_id, name, description, url, category, language, country, created_at, updated_at, status FROM sources"
	rows, err := main.db.Query(query)
	if err != nil {
		return nil, err
	}

	sources := make([]*Source, 0)
	for rows.Next() {
		source := &Source{}
		rows.Scan(source.Id_, source.ApiSourceName, source.SourceId, source.Description, source.Url, source.Category, source.Language, source.Country, source.CreatedAt, source.UpdatedAt, source.Status)
		sources = append(sources, source)
	}

	return sources, nil
}

func (main *DbMain) SaveSources(sources []*Source) (int64, error) {
	query := "INSERT INTO sources (api_source_name, s_id, name, description, url, category, language, country) VALUES "
	parameterHolders := make([]string, len(sources))
	parameters := make([]interface{}, 0)

	for index, source := range sources {
		parameterHolders[index] = "(?, ?, ?, ?, ?, ?, ?, ?)"
		parameters = append(parameters, source.ApiSourceName, source.SourceId, source.Name, source.Description, source.Url, source.Category, source.Language, source.Country)
	}

	query = query + strings.Join(parameterHolders, ",")
	return main.batchInsert(query, parameters)
}

func (main *DbMain) GetArticles() ([]*Article, error) {
	query := "SELECT _id, api_source_name, source_id, source_name, author, title, description, url, url_to_image, published_at, created_at, updated_at, status FROM articles"
	rows, err := main.db.Query(query)
	if err != nil {
		return nil, err
	}

	articles := make([]*Article, 0)
	for rows.Next() {
		article := &Article{}
		rows.Scan(article.Id_, article.ApiSourceName, article.SourceId, article.SourceName, article.Author, article.Title, article.Description, article.Url, article.UrlToImage, article.PublishedAt, article.CreatedAt, article.UpdatedAt, article.Status)
		articles = append(articles, article)
	}

	return articles, nil
}

func (main *DbMain) SaveArticles(articles []*Article) (int64, error) {
	query := "INSERT INTO articles (api_source_name, source_id, source_name, author, title, description, url, url_to_image, published_at) VALUES "
	parameterHolders := make([]string, len(articles))
	parameters := make([]interface{}, 0)

	for index, article := range articles {
		parameterHolders[index] = "(?, ?, ?, ?, ?, ?, ?, ?, ?)"
		parameters = append(parameters, article.ApiSourceName, article.SourceId, article.SourceName, article.Author, article.Title, article.Description, article.Url, article.UrlToImage, article.PublishedAt)
	}

	query = query + strings.Join(parameterHolders, ",")
	return main.batchInsert(query, parameters)
}

func (main *DbMain) batchInsert(query string, parameters ...interface{}) (int64, error) {
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
