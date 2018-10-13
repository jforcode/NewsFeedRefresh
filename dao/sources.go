package dao

import (
	"errors"
	"strings"

	"github.com/golang/glog"
)

type Source struct {
	dbRecord
	Id_           int
	ApiSourceName string
	SourceId      string
	Name          string
	Description   string
	Url           string
	Category      string
	Language      string
	Country       string
}

func (main *Dao) GetSources() ([]*Source, error) {
	prefix := "main.Dao.GetSources"
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

func (main *Dao) SaveSources(sources []*Source) (int64, error) {
	query := "INSERT INTO sources (api_source_name, s_id, name, description, url, category, language, country) VALUES "
	parameterHolders := make([]string, len(sources))
	parameters := make([]interface{}, 0)

	for index, source := range sources {
		parameterHolders[index] = "(?, ?, ?, ?, ?, ?, ?, ?)"
		parameters = append(parameters, source.ApiSourceName, source.SourceId, source.Name, source.Description, source.Url, source.Category, source.Language, source.Country)
	}

	query = query + strings.Join(parameterHolders, ",")
	return main.prepareAndExec(query, parameters...)
}

func (main *Dao) ClearSources(apiSourceName string) (int64, error) {
	query := "DELETE FROM sources WHERE api_source_name = ?"
	return main.prepareAndExec(query, apiSourceName)
}
