package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jforcode/NewsApiSDK"
	"github.com/jforcode/NewsFeedRefresh/dao"
	"github.com/jforcode/Util"
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
	flags := make(map[string]string)
	flags["parseTime"] = "true"

	db, err := util.Db.GetDb(user, password, host, database, flags)
	if err != nil {
		panic(err)
	}

	api := newsApi.NewNewsApi(apiUrl, apiKey)
	dao := dao.New(db)
	dailyRefresher := newsApi.NewRefresher(api)

	refresher := &Refresher{
		api:                    api,
		dao:                    dao,
		dailyRefresher:         dailyRefresher,
		sourceName:             "news_api",
		defaultNumTransactions: 1000,
	}

	err = refresher.DoInitialChecks()
	if err != nil {
		panic(err)
	}

	err = refresher.StartRefresh()
	if err != nil {
		panic(err)
	}
}
