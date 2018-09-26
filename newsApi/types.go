package newsApi

import (
	"time"
)

type ArticleSource struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Source struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Category    string `json:"category"`
	Language    string `json:"language"`
	Country     string `json:"country"`
}

type Article struct {
	Source      ArticleSource `json:"source"`
	Author      string        `json:"author"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	URL         string        `json:"url"`
	URLToImage  string        `json:"urlToImage"`
	PublishedAt time.Time     `json:"publishedAt"`
}

type SourceApiResponse struct {
	Status       string   `json:"status"`
	Sources      []Source `json:"sources"`
	ErrorCode    string   `json:"code"`
	ErrorMessage string   `json:"message"`
}

type NewsApiResponse struct {
	Status       string    `json:"status"`
	TotalResults int       `json:"totalResults"`
	Articles     []Article `json:"articles"`
	ErrorCode    string    `json:"code"`
	ErrorMessage string    `json:"message"`
}
