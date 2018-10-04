package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
	"github.com/magiconair/properties"
)

func main() {
	p := properties.MustLoadFile("app.properties", properties.UTF8)
	dataSource := p.GetString("datasource", "")
	apiKey := p.GetString("apiKey", "")
	apiUrl := p.GetString("apiUrl", "")

	db, err := sql.Open("mysql", dataSource)
	if err != nil {
		panic(err)
	}

	api := &NewsApi{}
	api.Init(apiUrl, apiKey, "news_api")

	dbMain := &DbMain{}
	dbMain.Init(db)

	util := &Util{}

	refresher := &Refresher{api, dbMain, util}

	glog.Infoln("Starting Refresh")
	err = refresher.StartRefresh()
	if err != nil {
		panic(err)
	}
}

func GetDb(user, password, host, database string) (*sql.DB, error) {
	datasource := user +
		":" + password +
		"@" + host +
		"/" + database +
		"?" + "parseTime=true"

	db, err := sql.Open("mysql", datasource)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
