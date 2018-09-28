package newsApi

import (
	"time"

	"github.com/jforcode/NewsFeedRefresh/modules/common"
)

// db types
type Flag struct {
	common.DbRecord
	Id_   int
	Key   string
	Value string
}

type ApiError struct {
	status  string
	code    string
	message string
}

func (apiErr *ApiError) Error() string {
	return apiErr.code + " (" + apiErr.status + ") - " + apiErr.message
}

// NewsAPI response types
type ApiSource struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Category    string `json:"category"`
	Language    string `json:"language"`
	Country     string `json:"country"`
}

type ApiSourcesResponse struct {
	Status       string      `json:"status"`
	Sources      []ApiSource `json:"sources"`
	ErrorCode    string      `json:"code"`
	ErrorMessage string      `json:"message"`
}

type ApiArticleSource struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type ApiArticle struct {
	Source      ApiArticleSource `json:"source"`
	Author      string           `json:"author"`
	Title       string           `json:"title"`
	Description string           `json:"description"`
	URL         string           `json:"url"`
	URLToImage  string           `json:"urlToImage"`
	PublishedAt time.Time        `json:"publishedAt"`
}

type ApiArticlesResponse struct {
	Status       string       `json:"status"`
	TotalResults int          `json:"totalResults"`
	Articles     []ApiArticle `json:"articles"`
	ErrorCode    string       `json:"code"`
	ErrorMessage string       `json:"message"`
}
