package common

import "time"

type DbRecord struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	Status    string
}
