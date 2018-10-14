package main

import (
	"strconv"
	"time"

	"github.com/apsdehal/go-logger"
	"github.com/jforcode/DeepError"
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
	log                    *logger.Logger
	debugAndErrorFormat    string
	sourceName             string
	defaultNumTransactions int
}

func (refr *Refresher) DoInitialChecks() error {
	fnName := "main.Refresher.Init"

	refr.log.Info("Checking remaining requests")
	err := refr.checkRemainingRequests()
	if err != nil {
		return deepError.New(fnName, "checking requests", err)
	}

	refr.log.Info("Checking sources")
	err = refr.checkSources()
	if err != nil {
		return deepError.New(fnName, " (check sources): ", err)
	}

	return nil
}

func (refr *Refresher) StartRefresh() error {
	fnName := "main.Refresher.StartRefresh"

	refr.log.Info("Getting remaining requests")
	remainingRequests, err := refr.getRemainingRequests()
	if err != nil {
		return deepError.New(fnName, " (get remaining requests): ", err)
	}
	refr.log.Debugf(refr.debugAndErrorFormat, fnName, "Got remaining requests", remainingRequests)

	refr.log.Info("Getting sources")
	sources, err := refr.dao.GetSources()
	if err != nil {
		return deepError.New(fnName, " (get sources): ", err)
	}
	refr.log.Debugf(refr.debugAndErrorFormat, fnName, "Got sources", sources)

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
		LastMomentMinutes: 5,
	}

	refr.log.Info("Doing daily refresh")
	refr.log.Debugf(refr.debugAndErrorFormat, fnName, "Refresher config", refresherConfig)
	go refr.dailyRefresher.DailyRefresh(refresherConfig, chArticles, chNumRequestsUpdated, chError)

	for {
		select {
		case apiArticles := <-chArticles:
			refr.log.Debugf(refr.debugAndErrorFormat, fnName, "Received articles from channel", apiArticles)
			articles := make([]*dao.Article, len(apiArticles))
			for index, apiArticle := range apiArticles {
				articles[index] = refr.convertArticle(apiArticle)
			}

			refr.log.DebugF(refr.debugAndErrorFormat, fnName, "Saving articles", articles)
			_, err := refr.dao.SaveArticles(articles)
			if err != nil {
				refr.log.ErrorF(refr.debugAndErrorFormat, fnName, "Saving articles", err)
			}

		case updatedRequests := <-chNumRequestsUpdated:
			refr.log.Debugf("%s -> %s : %d", fnName, "Received updated requests from channel", updatedRequests)
			remainingRequests -= updatedRequests
			err := refr.setRemainingRequests(remainingRequests)
			if err != nil {
				refr.log.ErrorF("%s -> %s: %+v", fnName, "Setting remaining requests", err)
			}

		case err := <-chError:
			refr.log.ErrorF("%s -> %s: %+v", fnName, "Error from channel", err)

		}
	}

	return nil
}

func (refr *Refresher) checkSources() error {
	fnName := "main.Refresher.checkSources"

	today := time.Now().UTC()
	monthStart := time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, today.Location())
	refr.log.Debugf(refr.debugAndErrorFormat, fnName, "Start of month", monthStart)

	flagSrcRefreshed, err := refr.dao.GetFlag("sources_refreshed", "bool")
	if err != nil {
		return deepError.New(fnName, " (get flag sources): ", err)
	}

	remainingRequests, err := refr.getRemainingRequests()
	if err != nil {
		return deepError.New(fnName, " (get remaining requests): ", err)
	}

	if (flagSrcRefreshed == nil || flagSrcRefreshed.UpdatedAt.Before(monthStart)) && remainingRequests > 0 {
		_, err := refr.dao.ClearSources(refr.sourceName)
		if err != nil {
			return deepError.New(fnName, " (clear sources): ", err)
		}

		apiSourcesResponse, err := refr.api.FetchSources(&newsApi.FetchSourcesParams{})
		if err != nil {
			return deepError.New(fnName, " (fetch sources): ", err)
		}

		err = refr.setRemainingRequests(remainingRequests - 1)
		if err != nil {
			return deepError.New(fnName, " (set remaining requests): ", err)
		}

		sources := make([]*dao.Source, len(apiSourcesResponse.Sources))
		for index, apiSource := range apiSourcesResponse.Sources {
			sources[index] = refr.convertSource(apiSource)
		}

		_, err = refr.dao.SaveSources(sources)
		if err != nil {
			return deepError.New(fnName, " (save sources): ", err)
		}

		err = refr.dao.SetFlag("sources_refreshed", "TRUE", "bool")
		if err != nil {
			return deepError.New(fnName, " (set flag sources): ", err)
		}
	}

	return nil
}

func (refr *Refresher) checkRemainingRequests() error {
	fnName := "main.Refresher.checkRemainingRequests"
	today := time.Now().UTC()
	dayStart := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())

	flagNumRequests, err := refr.dao.GetFlag("remaining_requests", "int")
	if err != nil {
		return deepError.New(fnName, " (get remaining requests): ", err)
	}

	if flagNumRequests == nil || flagNumRequests.UpdatedAt.Before(dayStart) {
		err := refr.setRemainingRequests(refr.defaultNumTransactions)
		if err != nil {
			return deepError.New(fnName, "setting remaining requests", err)
		}
	}

	return nil
}

func (refr *Refresher) getRemainingRequests() (int, error) {
	fnName := "main.Refresher.getRemainingRequests"
	flagNumRequests, err := refr.dao.GetFlag("remaining_requests", dao.FlagTypeInt)
	if err != nil {
		return -1, deepError.New(fnName, "getting flag", err)
	}
	if flagNumRequests == nil {
		return -1, deepError.DeepErr{Function: fnName, Message: "flag for remaining requests should not be nil"}
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
