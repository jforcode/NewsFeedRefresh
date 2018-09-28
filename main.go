package main

import (
	"database/sql"
	"fmt"

	"github.com/jforcode/NewsFeedRefresh/modules/feedSrv"
	"github.com/jforcode/NewsFeedRefresh/modules/newsApi"
	"github.com/magiconair/properties"
)

func main() {
	p := properties.MustLoadFile("app.properties", properties.UTF8)
	dataSource := p.GetString("datasource", "")
	apiKey := p.GetString("apiKey", "")
	apiUrl := p.GetString("apiUrl", "")
	requestLimit := p.GetInt("requestLimit", 0)

	db, err := initDb(dataSource)
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

	chSources := make(chan []newsApi.ApiSource)
	chArticles := make(chan []newsApi.ApiArticle)
	go apiMain.StartFetch(chSources, chArticles)

	select {
	case <-chSources:
		sources := <-chSources
		// convert, then save
		fmt.Println(sources)

	case <-chArticles:
		articles := <-chArticles
		fmt.Println(articles)
	}

	fmt.Println(feedMain)
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
