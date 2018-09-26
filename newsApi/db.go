package newsApi

import (
	"database/sql"
)

type MainDb struct {
	db *sql.DB
}

func (srv *MainDb) Init(db *sql.DB) {
	srv.db = db
}

// TODO:
func (srv *MainDb) GetFlag(key string) (*Flag, error) {
	return nil, nil
}

// TODO
func (srv *MainDb) SetFlag(flag *Flag) error {
	return nil
}
