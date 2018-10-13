package main

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang/glog"
	"github.com/jforcode/NewsApiSDK"
	"github.com/jforcode/NewsFeedRefresh/dao"
)

type IRefresher interface {
	DoInitialChecks() error
	StartRefresh() error
}

type Refresher struct {
	api                    newsApi.INewsApi
	dailyRefresher         newsApi.IRefresher
	dao                    dao.IDao
	sourceName             string
	defaultNumTransactions int
}

func (refr *Refresher) DoInitialChecks() error {
	prefix := "main.Refresher.Init"

	err := refr.checkRemainingRequests()
	if err != nil {
		return errors.New(prefix + " (check requests): " + err.Error())
	}

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

	sources, err := refr.dao.GetSources()
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
		articles := make([]*dao.Article, len(apiArticles))
		for index, apiArticle := range apiArticles {
			articles[index] = refr.convertArticle(apiArticle)
		}

		_, err := refr.dao.SaveArticles(articles)
		if err != nil {
			glog.Errorln(prefix+" (save articles): "+err.Error(), articles)
		}

	case updatedRequests := <-chNumRequestsUpdated:
		remainingRequests -= updatedRequests
		err := refr.setRemainingRequests(remainingRequests)
		if err != nil {
			glog.Errorln(prefix+" (set remaining requests): "+err.Error(), remainingRequests)
		}

	case err := <-chError:
		glog.Errorln(err)

	}

	return nil
}

func (refr *Refresher) checkSources() error {
	prefix := "main.Refresher.checkSources"

	today := time.Now().UTC()
	monthStart := time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, today.Location())

	flagSrcRefreshed, err := refr.dao.GetFlag("sources_refreshed", "bool")
	if err != nil {
		return errors.New(prefix + " (get flag sources): " + err.Error())
	}

	remainingRequests, err := refr.getRemainingRequests()
	if err != nil {
		return errors.New(prefix + " (get remaining requests): " + err.Error())
	}

	if (flagSrcRefreshed == nil || flagSrcRefreshed.UpdatedAt.Before(monthStart)) && remainingRequests > 0 {
		_, err := refr.dao.ClearSources(refr.sourceName)
		if err != nil {
			return errors.New(prefix + " (clear sources): " + err.Error())
		}

		apiSourcesResponse, err := refr.api.FetchSources(&newsApi.FetchSourcesParams{})
		if err != nil {
			return errors.New(prefix + " (fetch sources): " + err.Error())
		}

		err = refr.setRemainingRequests(remainingRequests - 1)
		if err != nil {
			return errors.New(prefix + " (set remaining requests): " + err.Error())
		}

		sources := make([]*dao.Source, len(apiSourcesResponse.Sources))
		for index, apiSource := range apiSourcesResponse.Sources {
			sources[index] = refr.convertSource(apiSource)
		}

		_, err = refr.dao.SaveSources(sources)
		if err != nil {
			return errors.New(prefix + " (save sources): " + err.Error())
		}

		err = refr.dao.SetFlag("sources_refreshed", "TRUE", "bool")
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

	flagNumRequests, err := refr.dao.GetFlag("remaining_requests", "int")
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
	flagNumRequests, err := refr.dao.GetFlag("remaining_requests", dao.FlagTypeInt)
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
	return refr.dao.SetFlag("remaining_requests", strconv.Itoa(remainingRequests), dao.FlagTypeInt)
}

func (refr *Refresher) convertSource(apiSource *newsApi.ApiSource) *dao.Source {
	return &dao.Source{
		ApiSourceName: refr.sourceName,
		SourceId:      apiSource.Id,
		Name:          apiSource.Name,
		Description:   apiSource.Description,
		Url:           apiSource.URL,
		Category:      apiSource.Category,
		Language:      apiSource.Language,
		Country:       apiSource.Country,
	}
}

func (refr *Refresher) convertArticle(apiArticle *newsApi.ApiArticle) *dao.Article {
	return &dao.Article{
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
}
