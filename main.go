package main

import (
	logger "github.com/apsdehal/go-logger"
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
	log, err := logger.New("main", 1)
	if err != nil {
		panic(err)
	}

	refresher := &Refresher{
		api:                    api,
		dao:                    dao,
		dailyRefresher:         dailyRefresher,
		log:                    log,
		debugAndErrorFormat:    "%s -> %s : %+v",
		sourceName:             "news_api",
		defaultNumTransactions: 1000,
	}

	err = refresher.DoInitialChecks()
	if err != nil {
		log.Fatalf("%s : Fatal error: %+v", "doing initial checks", err)
		panic(err)
	}

	err = refresher.StartRefresh()
	if err != nil {
		log.FatalF("%s: Fatal error: %+v", "starting refresh", err)
		panic(err)
	}
}
