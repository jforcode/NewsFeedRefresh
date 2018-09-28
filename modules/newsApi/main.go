package newsApi

import (
	"database/sql"
	"errors"
	"strconv"
	"time"

	"github.com/robfig/cron"
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
	c := cron.New()
	c.AddFunc("@daily", func() {
		main.db.SetFlag("remaining_requests", strconv.Itoa(requestLimit))
	})

	remTransFlag, err := main.db.GetFlag("remaining_requests")
	if err != nil {
		return nil, err
	}

	today := time.Now().UTC()
	tomorrow := time.Date(today.Year(), today.Month(), today.Day()+1, 0, 0, 0, 0, today.Location())
	if remTransFlag.UpdatedAt.Before(tomorrow) {
		main.db.SetFlag("remaining_requests", strconv.Itoa(requestLimit))
	}

	return &main, nil
}

func (main *Main) StartFetch(chSources chan []ApiSource, chArticles chan []ApiArticle) {
	// get the number of remaining transactions from db
	// get the time remaining till end of day
	// get sources, and also send it to save.
	// based on the above data, spread out fetching of articles over the time and sources, with pagination
	// and then send it to save
}