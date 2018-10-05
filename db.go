package main

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"

	"github.com/golang/glog"
)

type DbMain struct {
	db *sql.DB
}

func (main *DbMain) Init(db *sql.DB) {
	main.db = db
}

// TODO: make typeTo as a proper type with constants
func (main *DbMain) GetFlag(key, typeTo string) (*Flag, error) {
	prefix := "main.DbMain.GetFlag"
	query := "SELECT _id, flag_key, flag_value, created_at, updated_at, status FROM news_api_flags WHERE flag_key = ?"

	rows, err := main.db.Query(query, key)
	if err != nil {
		return nil, errors.New(prefix + " (query): " + err.Error())
	}
	defer rows.Close()

	if rows.Next() {
		flag := &Flag{}
		var value string
		rows.Scan(&flag.Id_, &flag.Key, &value, &flag.CreatedAt, &flag.UpdatedAt, &flag.Status)

		switch typeTo {
		case "string":
			flag.Value = value

		case "int":
			flag.Value, err = strconv.Atoi(value)
			if err != nil {
				return nil, errors.New(prefix + "Invalid type int: " + err.Error())
			}

		case "bool":
			flag.Value, err = strconv.ParseBool(value)
			if err != nil {
				return nil, errors.New(prefix + "Invalid type bool: " + err.Error())
			}
		}

		return flag, nil
	}

	return nil, nil
}

func (main *DbMain) SetFlag(key, value, typeTo string) error {
	prefix := "main.DbMain.SetFlag"
	flag, err := main.GetFlag(key, typeTo)
	if err != nil {
		return errors.New(prefix + " (get flag): " + err.Error())
	}

	if flag != nil {
		return main.updateFlag(key, value)
	} else {
		return main.createFlag(key, value)
	}
}

func (main *DbMain) updateFlag(key, value string) error {
	glog.Infoln("Creating flag with key: " + key + ", value: " + value)
	prefix := "main.DbMain.createFlag"
	query := "UPDATE news_api_flags SET flag_value = ? WHERE flag_key = ?"
	stmt, err := main.db.Prepare(query)
	if err != nil {
		return errors.New(prefix + " (prepare): " + err.Error())
	}
	defer stmt.Close()

	_, err = stmt.Exec(value, key)
	if err != nil {
		return errors.New(prefix + " (exec): " + err.Error())
	}

	return nil
}

func (main *DbMain) createFlag(key, value string) error {
	glog.Infoln("Updating flag with key: " + key + ", value: " + value)
	prefix := "main.DbMain.updateFlag"
	query := "INSERT INTO news_api_flags (flag_key, flag_value) VALUES (?, ?)"
	stmt, err := main.db.Prepare(query)
	if err != nil {
		return errors.New(prefix + " (Prepare): " + err.Error())
	}
	defer stmt.Close()

	_, err = stmt.Exec(key, value)
	if err != nil {
		return errors.New(prefix + " (Prepare): " + err.Error())
	}

	return nil
}

func (main *DbMain) GetSources() ([]*Source, error) {
	prefix := "main.DbMain.GetSources"
	query := "SELECT _id, api_source_name, s_id, name, description, url, category, language, country, created_at, updated_at, status FROM sources"
	rows, err := main.db.Query(query)
	if err != nil {
		return nil, errors.New(prefix + " (Query): " + err.Error())
	}

	sources := make([]*Source, 0)
	for rows.Next() {
		var source Source
		err := rows.Scan(&source.Id_, &source.ApiSourceName, &source.SourceId, &source.Name, &source.Description, &source.Url, &source.Category, &source.Language, &source.Country, &source.CreatedAt, &source.UpdatedAt, &source.Status)
		if err != nil {
			glog.Errorln(err)
		}

		sources = append(sources, &source)
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
	prefix := "main.DbMain.GetArticles"
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

func (main *DbMain) batchInsert(query string, parameters []interface{}) (int64, error) {
	prefix := "main.DbMain.batchInsert"
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
