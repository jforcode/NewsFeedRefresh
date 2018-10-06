package main

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang/glog"
	"github.com/jforcode/NewsApi"
)

type Refresher struct {
	api                    *newsApi.NewsApi
	dailyRefresher         *newsApi.Refresher
	dbMain                 *DbMain
	sourceName             string
	defaultNumTransactions int
}

func (refr *Refresher) Init(api *newsApi.NewsApi, dbMain *DbMain) error {
	prefix := "main.Refresher.Init"

	refr.api = api
	refr.dbMain = dbMain
	refr.sourceName = "news_api"
	refr.defaultNumTransactions = 1000

	glog.Infoln("Checking wether to update remaining requests")
	err := refr.checkRemainingRequests()
	if err != nil {
		return errors.New(prefix + " (check requests): " + err.Error())
	}

	glog.Infoln("Checking wether to refresh sources")
	err = refr.checkSources()
	if err != nil {
		return errors.New(prefix + " (check sources): " + err.Error())
	}

	return nil
}

func (refr *Refresher) StartRefresh() error {
	prefix := "main.Refresher.StartRefresh"

	remainingRequests, err := refr.getRemainingRequests()
	if err != nil {
		return errors.New(prefix + " (get remaining requests): " + err.Error())
	}

	sources, err := refr.dbMain.GetSources()
	if err != nil {
		return errors.New(prefix + " (get sources): " + err.Error())
	}

	sourceIds := make([]string, len(sources))
	for index, source := range sources {
		sourceIds[index] = source.SourceId
	}

	chArticles := make(chan []*newsApi.ApiArticle)
	chNumRequestsUpdated := make(chan int)
	chError := make(chan error)

	refresherConfig := &newsApi.RefresherConfig{
		RemainingRequests: remainingRequests,
		SourceIds:         sourceIds,
		PageSize:          100,
	}

	go refr.dailyRefresher.DailyRefresh(refresherConfig, chArticles, chNumRequestsUpdated, chError)
	select {
	case apiArticles := <-chArticles:
		articles := make([]*Article, len(apiArticles))
		for index, apiArticle := range apiArticles {
			articles[index] = refr.convertArticle(apiArticle)
		}

		go refr.dbMain.SaveArticles(articles)

	case updatedRequests := <-chNumRequestsUpdated:
		remainingRequests -= updatedRequests
		go refr.setRemainingRequests(remainingRequests)

	case err := <-chError:
		fmt.Println(err)

	}

	glog.Infoln("Fetching articles")

	return nil
}

func (refr *Refresher) checkSources() error {
	prefix := "main.Refresher.checkSources"
	today := time.Now().UTC()
	monthStart := time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, today.Location())

	flagSrcRefreshed, err := refr.dbMain.GetFlag("sources_refreshed", "bool")
	if err != nil {
		return errors.New(prefix + " (get flag sources): " + err.Error())
	}

	remainingRequests, err := refr.getRemainingRequests()
	if err != nil {
		return errors.New(prefix + " (get remaining requests): " + err.Error())
	}

	if (flagSrcRefreshed == nil || flagSrcRefreshed.UpdatedAt.Before(monthStart)) && remainingRequests > 0 {
		glog.Infoln("refreshing sources. ", flagSrcRefreshed)

		_, err := refr.dbMain.ClearSources(refr.sourceName)
		if err != nil {
			return errors.New(prefix + " (clear sources): " + err.Error())
		}

		apiSourcesResponse, err := refr.api.FetchSources()
		if err != nil {
			return errors.New(prefix + " (fetch sources): " + err.Error())
		}

		err = refr.setRemainingRequests(remainingRequests - 1)
		if err != nil {
			return errors.New(prefix + " (set remaining requests): " + err.Error())
		}

		sources := make([]*Source, len(apiSourcesResponse.Sources))
		for index, apiSource := range apiSourcesResponse.Sources {
			sources[index] = refr.convertSource(apiSource)
		}

		_, err = refr.dbMain.SaveSources(sources)
		if err != nil {
			return errors.New(prefix + " (save sources): " + err.Error())
		}

		err = refr.dbMain.SetFlag("sources_refreshed", "TRUE", "bool")
		if err != nil {
			return errors.New(prefix + " (set flag sources): " + err.Error())
		}
	}

	return nil
}

func (refr *Refresher) checkRemainingRequests() error {
	prefix := "main.Refresher.checkRemainingRequests"
	today := time.Now().UTC()
	dayStart := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())

	flagNumRequests, err := refr.dbMain.GetFlag("remaining_requests", "int")
	if err != nil {
		return errors.New(prefix + " (get remaining requests): " + err.Error())
	}

	if flagNumRequests == nil || flagNumRequests.UpdatedAt.Before(dayStart) {
		glog.Infoln("resetting remaining requests for today.", flagNumRequests)

		err := refr.setRemainingRequests(refr.defaultNumTransactions)
		if err != nil {
			return errors.New(prefix + " (set remaining requests): " + err.Error())
		}
	}

	return nil
}

func (refr *Refresher) getRemainingRequests() (int, error) {
	prefix := "main.Refresher.getRemainingRequests"
	flagNumRequests, err := refr.dbMain.GetFlag("remaining_requests", "int")
	if err != nil {
		return -1, errors.New(prefix + " (get flag) " + err.Error())
	}
	if flagNumRequests == nil {
		return -1, errors.New(prefix + " flag for remaining requests should not be nil")
	}

	remainingRequests := flagNumRequests.Value.(int)
	return remainingRequests, nil
}

func (refr *Refresher) setRemainingRequests(remainingRequests int) error {
	return refr.dbMain.SetFlag("remaining_requests", strconv.Itoa(remainingRequests), "int")
}

// TODO: json anyway to copy from one structure to another only some values
// basicaly any way to reduce the code here.
func (refr *Refresher) convertSource(apiSource *newsApi.ApiSource) *Source {
	source := Source{
		ApiSourceName: refr.sourceName,
		SourceId:      apiSource.Id,
		Name:          apiSource.Name,
		Description:   apiSource.Description,
		Url:           apiSource.URL,
		Category:      apiSource.Category,
		Language:      apiSource.Language,
		Country:       apiSource.Country,
	}

	return &source
}

func (refr *Refresher) convertArticle(apiArticle *newsApi.ApiArticle) *Article {
	article := Article{
		ApiSourceName: refr.sourceName,
		Author:        apiArticle.Author,
		Title:         apiArticle.Title,
		Description:   apiArticle.Description,
		Url:           apiArticle.URL,
		UrlToImage:    apiArticle.URLToImage,
		PublishedAt:   apiArticle.PublishedAt,
		SourceId:      apiArticle.Source.Id,
		SourceName:    apiArticle.Source.Name,
	}

	return &article
}
