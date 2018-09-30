package main

import (
	"database/sql"
	"fmt"

	"github.com/jforcode/NewsFeedRefresh/modules/common"
	"github.com/jforcode/NewsFeedRefresh/modules/feedSrv"
	"github.com/jforcode/NewsFeedRefresh/modules/newsApi"
	"github.com/magiconair/properties"
)

type Main struct {
	apiMain  *newsApi.Main
	feedMain *feedSrv.Main
}

func main() {
	p := properties.MustLoadFile("app.properties", properties.UTF8)
	dataSource := p.GetString("datasource", "")
	apiKey := p.GetString("apiKey", "")
	apiUrl := p.GetString("apiUrl", "")
	requestLimit := p.GetInt("requestLimit", 0)

	db, err := sql.Open("mysql", dataSource)
	if err != nil {
		fmt.Println("Error: " + err.Error())
		return
	}

	apiMain, err := newsApi.Init(db, apiKey, apiUrl, requestLimit)
	if err != nil {
		fmt.Println("Error: " + err.Error())
		return
	}

	feedMain, err := feedSrv.Init(db)
	if err != nil {
		fmt.Println("Error: " + err.Error())
		return
	}

	main := &Main{}
	main.Init(apiMain, feedMain)
	main.Start()
}

func (main *Main) Init(apiMain *newsApi.Main, feedMain *feedSrv.Main) {
	main.apiMain = apiMain
	main.feedMain = feedMain
}

func (main *Main) Start() {
	chSources := make(chan [](*common.Source))
	chArticles := make(chan [](*common.Article))

	go main.apiMain.StartFetch(chSources, chArticles)
	select {
	case <-chSources:
		go main.feedMain.SaveSources(<-chSources)
	case <-chArticles:
		go main.feedMain.SaveArticles(<-chArticles)

		// TODO: probably error channels
	}
}
