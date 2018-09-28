package feedSrv

import (
	"time"

	"github.com/jforcode/NewsFeedRefresh/modules/common"
)

type Article struct {
	common.DbRecord
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
	common.DbRecord
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
