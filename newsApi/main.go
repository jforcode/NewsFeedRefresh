package newsApi

import (
	"database/sql"
	"errors"
)

type NewsApi struct {
	db  *MainDb
	api *Api
}

func (main *NewsApi) Init(db *sql.DB, apiKey, apiUrl string) error {
	if db == nil || apiKey == "" || apiUrl == "" {
		return errors.New("Invalid parameters")
	}

	main.db = &MainDb{}
	main.api = &Api{}

	main.db.Init(db)
	main.api.Init(apiUrl, apiKey)

	return nil
}

func StartFetch() {
	// need params:
	// - apiBaseUrl, apiKey
	// db flags:
	// - number of transactions remaining today.
	// - if number of transactions was last updated yesterday, reset it to 1000
}
