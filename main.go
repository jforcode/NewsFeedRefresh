package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jforcode/DbUtil"
	"github.com/jforcode/NewsApi"
	"github.com/magiconair/properties"
)

func main() {
	p := properties.MustLoadFile("app.properties", properties.UTF8)
	apiKey := p.GetString("apiKey", "")
	apiUrl := p.GetString("apiUrl", "")

	user := p.GetString("user", "")
	password := p.GetString("password", "")
	host := p.GetString("host", "")
	database := p.GetString("db", "")
	params := make(map[string]string)
	params["parseTime"] = "true"

	db, err := dbUtil.GetDb(user, password, host, database, params)
	if err != nil {
		panic(err)
	}

	api := &newsApi.NewsApi{}
	api.Init(apiUrl, apiKey)

	dbMain := &DbMain{}
	dbMain.Init(db)

	refresher := &Refresher{}
	err = refresher.Init(api, dbMain)
	if err != nil {
		panic(err)
	}

	err = refresher.StartRefresh()
	if err != nil {
		panic(err)
	}
}
