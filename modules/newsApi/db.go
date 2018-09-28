package newsApi

import (
	"database/sql"
)

type MainDb struct {
	db *sql.DB
}

func (srv *MainDb) Init(db *sql.DB) error {
	srv.db = db
	return nil
}

// TODO:
func (srv *MainDb) GetFlag(key string) (*Flag, error) {
	return nil, nil
}

// TODO
func (srv *MainDb) SetFlag(key, value string) error {
	return nil
}
