package newsApi

import (
	"database/sql"
	"errors"
	"strconv"
	"time"

	"github.com/jforcode/NewsFeedRefresh/modules/common"
	"github.com/robfig/cron"
)

const (
	flagRemainingRequests string = "remaining_requests"
	flagRefreshSources    string = "refresh_sources"
)

type Main struct {
	db  *MainDb
	api *Api
}

func Init(db *sql.DB, apiKey, apiUrl string, requestLimit int) (*Main, error) {
	if db == nil || apiKey == "" || apiUrl == "" {
		return nil, errors.New("Invalid parameters")
	}

	main := Main{db: &MainDb{}, api: &Api{}}

	main.db.Init(db)
	main.api.Init(apiUrl, apiKey)

	// TODO: replace it with native timer

	main.refreshNumTransactions(requestLimit)
	return &main, nil
}

func (main *Main) refreshNumTransactions(requestLimit int) error {
	c := cron.New()
	c.AddFunc("@daily", func() {
		main.db.SetFlag(flagRemainingRequests, strconv.Itoa(requestLimit))
	})

	remTransFlag, err := main.db.GetFlag(flagRemainingRequests)
	if err != nil {
		return err
	}

	today := time.Now().UTC()
	firstMoment := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	if remTransFlag.UpdatedAt.Before(firstMoment) {
		main.db.SetFlag(flagRemainingRequests, strconv.Itoa(requestLimit))
	}

	return nil
}

func (main *Main) refreshSources() {
	// TODO: refreshing sources, monthly/weekly.
	// should the duration be a parameter/property?
	// refresh flags code can be combined, if passing duration as parameter.
}

func (main *Main) StartFetch(chSources chan [](*common.Source), chArticles chan [](*common.Article)) {

	// get the number of remaining transactions from db
	// get the time remaining till end of day
	// get sources, and also send it to save.
	// based on the above data, spread out fetching of articles over the time and sources, with pagination
	// and then send it to save
}
