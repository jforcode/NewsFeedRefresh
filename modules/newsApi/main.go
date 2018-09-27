package newsApi

import (
	"database/sql"
	"errors"
)

type Main struct {
	db  *MainDb
	api *Api
}

func Init(db *sql.DB, apiKey, apiUrl string) (*Main, error) {
	if db == nil || apiKey == "" || apiUrl == "" {
		return nil, errors.New("Invalid parameters")
	}

	main := Main{db: &MainDb{}, api: &Api{}}

	main.db.Init(db)
	main.api.Init(apiUrl, apiKey)

	return &main, nil
}

func (main *Main) StartFetch() {
	// need params:
	// - apiBaseUrl, apiKey
	// db flags:
	// - number of transactions remaining today.
	// - if number of transactions was last updated yesterday, reset it to 1000
}
