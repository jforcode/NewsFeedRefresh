package main

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

type Refresher struct {
	api    *NewsApi
	dbMain *DbMain
	util   *Util
}

func (refr *Refresher) StartRefresh() error {
	prefix := "main.Refresher.StartRefresh"
	err := refr.CheckSources()
	if err != nil {
		return errors.New(prefix + " (check sources): " + err.Error())
	}

	err = refr.CheckRemainingRequests()
	if err != nil {
		return errors.New(prefix + " (check requests): " + err.Error())
	}

	flagNumRequests, err := refr.dbMain.GetFlag("remaining_requests")
	if err != nil {
		return errors.New(prefix + " (get flag requests): " + err.Error())
	}
	remainingRequests := refr.util.GetInt(flagNumRequests, 1000)

	sources, err := refr.dbMain.GetSources()
	if err != nil {
		return errors.New(prefix + " (get sources): " + err.Error())
	}

	chArticles := make(chan []*Article)

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

	flagSrcRefreshed, err := refr.dbMain.GetFlag("sources_refreshed")
	if err != nil {
		return errors.New(prefix + " (get flag sources): " + err.Error())
	}

	if flagSrcRefreshed == nil || flagSrcRefreshed.UpdatedAt.Before(monthStart) {
		sources, err := refr.api.FetchSources()
		if err != nil {
			return errors.New(prefix + " (fetch sources): " + err.Error())
		}

		_, err = refr.dbMain.SaveSources(sources)
		if err != nil {
			return errors.New(prefix + " (save sources): " + err.Error())
		}

		err = refr.dbMain.SetFlag("sources_refreshed", "TRUE")
		if err != nil {
			return errors.New(prefix + " (set flag sources): " + err.Error())
		}
	}

	return nil
}

func (refr *Refresher) CheckRemainingRequests() error {
	today := time.Now()
	dayStart := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())

	flagNumRequests, err := refr.dbMain.GetFlag("remaining_requests")
	if err != nil {
		return err
	}

	if flagNumRequests == nil || flagNumRequests.UpdatedAt.Before(dayStart) {
		err := refr.dbMain.SetFlag("remaining_requests", "1000")
		if err != nil {
			return err
		}
	}

	return nil
}

// TODO: should use struct for variables, and in a config, so that it can be changed easily
func (refr *Refresher) FetchArticles(sources []*Source, remainingRequests int, chArticles chan []*Article) {
	batchSize := 20
	lenSources := len(sources)
	sourceIds := make([]string, refr.util.MinInt(batchSize, lenSources))
	lenIds := len(sourceIds)
	pageNum := 1
	pageSize := 100
	today := time.Now()
	// 30 minutes before end of day
	lastMoment := time.Date(today.Year(), today.Month(), today.Day()+1, 0, -30, 0, 0, today.Location())

	for {
		for index, source := range sources {
			sourceIds[index%lenIds] = source.Name

			if (index+1)%batchSize == 0 || index == lenSources-1 {
				if time.Now().After(lastMoment) || remainingRequests <= 0 {
					close(chArticles)
					return
				}

				remainingRequests--
				articles, err := refr.api.FetchArticles(sourceIds, pageNum, pageSize)
				if err != nil {
					fmt.Println(err.Error())
				} else {
					chArticles <- articles
				}

				sourceIds = make([]string, refr.util.MinInt(batchSize, lenSources-index))
				lenIds = len(sourceIds)
			}
		}

		pageNum++
		refr.dbMain.SetFlag("remaining_requests", strconv.Itoa(remainingRequests))
		time.Sleep(1 * time.Hour)

	}
}
