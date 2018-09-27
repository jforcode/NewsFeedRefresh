package feedSrv

import (
	"database/sql"
	"errors"
)

type Main struct {
	db *sql.DB
}

func Init(db *sql.DB) (*Main, error) {
	if db == nil {
		return nil, errors.New("Invalid parameters")
	}

	main := Main{db: db}
	return &main, nil
}
