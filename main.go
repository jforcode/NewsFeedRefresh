package main

import (
	"database/sql"

	"github.com/jforcode/NewsFeedRefresh/newsApi"
	"github.com/magiconair/properties"
)

func main() {
	p := properties.MustLoadFile("app.properties", properties.UTF8)
	dataSource := p.GetString("datasource", "")
	apiKey := p.GetString("apiKey", "")
	apiUrl := p.GetString("apiUrl", "")
	db, _ := initDb(dataSource)

	newsApi.Init(db, apiKey, apiUrl)
	newsApi.StartFetch()
}

func initDb(dataSource string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dataSource)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
