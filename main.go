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

	main := &Refresher{api, dbMain, util}

	glog.Infoln("Starting Refresh")
	err = main.StartRefresh()
	if err != nil {
		panic(err)
	}
}
