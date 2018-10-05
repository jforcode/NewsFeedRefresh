package main

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang/glog"
)

// TODO: make all parameters as config.. like number of requests in a day, tiem for refreshing sources, batchsize etc.

type Refresher struct {
	api    *NewsApi
	dbMain *DbMain
	util   *Util
}

func (refr *Refresher) Init(api *NewsApi, dbMain *DbMain, util *Util) error {
	prefix := "main.Refresher.Init"

	refr.api = api
	refr.dbMain = dbMain
	refr.util = util

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

	return nil
}

func (refr *Refresher) StartRefresh() error {
	prefix := "main.Refresher.StartRefresh"

	remainingRequests, err := refr.GetRemainingRequests()
	if err != nil {
		return errors.New(prefix + " (get remaining requests): " + err.Error())
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
	today := time.Now().UTC()
	monthStart := time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, today.Location())

	flagSrcRefreshed, err := refr.dbMain.GetFlag("sources_refreshed", "bool")
	if err != nil {
		return errors.New(prefix + " (get flag sources): " + err.Error())
	}

	remainingRequests, err := refr.GetRemainingRequests()
	if err != nil {
		return errors.New(prefix + " (get remaining request): " + err.Error())
	}

	if (flagSrcRefreshed == nil || flagSrcRefreshed.UpdatedAt.Before(monthStart)) && remainingRequests > 0 {
		glog.Infoln("refreshing sources. ", flagSrcRefreshed)
		sources, err := refr.api.FetchSources()
		if err != nil {
			return errors.New(prefix + " (fetch sources): " + err.Error())
		}

		err = refr.SetRemainingRequests(remainingRequests - 1)
		if err != nil {
			return errors.New(prefix + " (set remaining requests): " + err.Error())
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
	today := time.Now().UTC()
	dayStart := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())

	flagNumRequests, err := refr.dbMain.GetFlag("remaining_requests", "int")
	if err != nil {
		return errors.New(prefix + " (get remaining requests): " + err.Error())
	}

	if flagNumRequests == nil || flagNumRequests.UpdatedAt.Before(dayStart) {
		glog.Infoln("resetting remaining requests for today.", flagNumRequests)

		err := refr.SetRemainingRequests(1000)
		if err != nil {
			return errors.New(prefix + " (set remaining requests): " + err.Error())
		}
	}

	return nil
}

// TODO: should use struct for variables, and in a config, so that it can be changed easily
// TODO: channel for errors
func (refr *Refresher) FetchArticles(sources []*Source, remainingRequests int, chArticles chan []*Article) {
	prefix := "main.Refresher.FetchArticles"
	batchSize := 20
	lenSources := len(sources)
	lenIds := refr.util.MinInt(batchSize, lenSources)
	sourceIds := make([]string, lenIds)
	pageNum := 1
	pageSize := 100
	today := time.Now().UTC()
	lastMoment := time.Date(today.Year(), today.Month(), today.Day()+1, 0, -30, 0, 0, today.Location())
	// without this, sourceIds will be out of order, which is OK, no impact whatsoever.
	// this is just so that code works as we expect.
	firstIndForBatch := 0

	for {
		for index, source := range sources {
			sourceIds[index-firstIndForBatch] = source.Name

			if index-firstIndForBatch+1 == lenIds {
				if time.Now().After(lastMoment) || remainingRequests <= 0 {
					glog.Infoln("Exiting ", time.Now(), lastMoment, remainingRequests)
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

				lenIds = refr.util.MinInt(batchSize, lenSources-index-1)
				sourceIds = make([]string, lenIds)
				firstIndForBatch = index + 1
			}
		}

		pageNum++
		err := refr.SetRemainingRequests(remainingRequests)
		if err != nil {
			fmt.Println(prefix + " (set remaining requests) " + err.Error())
		}

		time.Sleep(1 * time.Minute)
	}
}

func (refr *Refresher) GetRemainingRequests() (int, error) {
	prefix := "main.Refresher.GetRemainingRequests"
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

func (refr *Refresher) SetRemainingRequests(remainingRequests int) error {
	return refr.dbMain.SetFlag("remaining_requests", strconv.Itoa(remainingRequests), "int")
}
