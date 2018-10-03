package main

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang/glog"
)

type Refresher struct {
	api    *NewsApi
	dbMain *DbMain
	util   *Util
}

func (refr *Refresher) StartRefresh() error {
	prefix := "main.Refresher.StartRefresh"

	glog.Infoln("Checking wether to update remaining requests")
	err := refr.CheckRemainingRequests()
	if err != nil {
		return errors.New(prefix + " (check requests): " + err.Error())
	}

	glog.Infoln("Checking wether to refresh sources")
	err = refr.CheckSources()
	if err != nil {
		return errors.New(prefix + " (check sources): " + err.Error())
	}

	flagNumRequests, err := refr.dbMain.GetFlag("remaining_requests", "int")
	if err != nil {
		return errors.New(prefix + " (get flag requests): " + err.Error())
	}
	remainingRequests := 1000
	if flagNumRequests != nil {
		remainingRequests = flagNumRequests.Value.(int)
	}

	sources, err := refr.dbMain.GetSources()
	if err != nil {
		return errors.New(prefix + " (get sources): " + err.Error())
	}

	chArticles := make(chan []*Article)

	glog.Infoln("Fetching articles")
	go refr.FetchArticles(sources, remainingRequests, chArticles)
	for articles := range chArticles {
		go refr.dbMain.SaveArticles(articles)
	}

	return nil
}

func (refr *Refresher) CheckSources() error {
	prefix := "main.Refresher.CheckSources"
	today := time.Now()
	monthStart := time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, today.Location())

	flagSrcRefreshed, err := refr.dbMain.GetFlag("sources_refreshed", "bool")
	if err != nil {
		return errors.New(prefix + " (get flag sources): " + err.Error())
	}

	if flagSrcRefreshed == nil || flagSrcRefreshed.UpdatedAt.Before(monthStart) {
		glog.Infoln("refreshing sources. ", flagSrcRefreshed)
		sources, err := refr.api.FetchSources()
		if err != nil {
			return errors.New(prefix + " (fetch sources): " + err.Error())
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

func (refr *Refresher) CheckRemainingRequests() error {
	prefix := "main.Refresher.CheckRemainingRequests"
	today := time.Now()
	dayStart := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())

	flagNumRequests, err := refr.dbMain.GetFlag("remaining_requests", "int")
	if err != nil {
		return errors.New(prefix + " (get remaining requests): " + err.Error())
	}

	if flagNumRequests == nil || flagNumRequests.UpdatedAt.Before(dayStart) {
		glog.Infoln("resetting remaining requests for today.", flagNumRequests)
		err := refr.dbMain.SetFlag("remaining_requests", "1000", "int")
		if err != nil {
			return errors.New(prefix + " (set remaining requests): " + err.Error())
		}
	}

	return nil
}

// TODO: should use struct for variables, and in a config, so that it can be changed easily
func (refr *Refresher) FetchArticles(sources []*Source, remainingRequests int, chArticles chan []*Article) {
	prefix := "main.Refresher.FetchArticles"
	batchSize := 20
	lenSources := len(sources)
	sourceIds := make([]string, refr.util.MinInt(batchSize, lenSources))
	lenIds := len(sourceIds)
	pageNum := 1
	pageSize := 100
	today := time.Now()
	lastMoment := time.Date(today.Year(), today.Month(), today.Day()+1, 0, -30, 0, 0, today.Location())

	for {
		for index, source := range sources {
			sourceIds[index%lenIds] = source.Name

			if (index+1)%batchSize == 0 || index == lenSources-1 {
				if time.Now().After(lastMoment) || remainingRequests <= 0 {
					glog.Infoln("Exiting ", lastMoment, remainingRequests)
					close(chArticles)
					return
				}

				glog.Infoln("Making a request for articles. index: ", index, ", remaining requests: ", remainingRequests, ", pageNum: ", pageNum, ", sourceIds: ", sourceIds)
				remainingRequests--
				articles, err := refr.api.FetchArticles(sourceIds, pageNum, pageSize)
				if err != nil {
					fmt.Println("Error: " + prefix + " (fetch articles): " + err.Error())
				} else {
					chArticles <- articles
				}

				sourceIds = make([]string, refr.util.MinInt(batchSize, lenSources-index))
				lenIds = len(sourceIds)
			}
		}

		pageNum++
		refr.dbMain.SetFlag("remaining_requests", strconv.Itoa(remainingRequests), "int")
		time.Sleep(1 * time.Minute)
	}
}
