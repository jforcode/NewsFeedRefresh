package common

import (
	"time"
)

type DbRecord struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	Status    string
}

type Article struct {
	DbRecord
	Id_           int
	ApiSourceName string
	SourceId      string
	SourceName    string
	Author        string
	Title         string
	Description   string
	Url           string
	UrlToImage    string
	PublishedAt   time.Time
}

type Source struct {
	DbRecord
	Id_           int
	ApiSourceName string
	SourceId      string
	Name          string
	Description   string
	Url           string
	Category      string
	Language      string
	Country       string
}
