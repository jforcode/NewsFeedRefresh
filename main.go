package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jforcode/NewsApiSDK"
	"github.com/jforcode/NewsFeedRefresh/dao"
	"github.com/magiconair/properties"
)

var (
	newsApiMain *newsApi.NewsApi
	daoMain     *dao.Dao
)

func main() {
	p := properties.MustLoadFile("app.properties", properties.UTF8)
	apiKey := p.GetString("apiKey", "")
	apiUrl := p.GetString("apiUrl", "")

	user := p.GetString("user", "")
	password := p.GetString("password", "")
	host := p.GetString("host", "")
	database := p.GetString("db", "")
	flags := make(map[string]string)
	flags["parseTime"] = "true"

	db, err := dbUtil.GetDb(user, password, host, database, flags)
	if err != nil {
		panic(err)
	}

	api := &newsApi.NewsApi{}
	api.Init(apiUrl, apiKey)

	dao := dao.New(db)

	refresher := &Refresher{}
	err = refresher.Init(api, dao)
	if err != nil {
		panic(err)
	}

	err = refresher.StartRefresh()
	if err != nil {
		panic(err)
	}
}
