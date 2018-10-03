package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type NewsApi struct {
	url        string
	key        string
	sourceName string
	client     *http.Client
}

func (api *NewsApi) Init(url, key, sourceName string) error {
	api.url = url
	api.key = key
	api.client = &http.Client{}
	api.sourceName = sourceName

	return nil
}

func (api *NewsApi) FetchSources() ([]*Source, error) {
	prefix := "main.NewsApi.FetchSources"
	bodyBytes, err := api.get("sources", nil)
	if err != nil {
		return nil, errors.New(prefix + " (Api Get): " + err.Error())
	}

	var sourceResponse ApiSourcesResponse
	if err := json.Unmarshal(bodyBytes, &sourceResponse); err != nil {
		return nil, errors.New(prefix + " (json unmarshal): " + err.Error())
	}

	sources := make([]*Source, len(sourceResponse.Sources))
	for index, apiSource := range sourceResponse.Sources {
		sources[index] = api.convertSource(&apiSource)
	}

	return sources, nil
}

func (api *NewsApi) FetchArticles(sourceIds []string, pageNum, pageSize int) ([]*Article, error) {
	prefix := "main.NewsApi.FetchArticles"
	params := make(map[string]string)
	params["sources"] = strings.Join(sourceIds, ",")
	params["page"] = strconv.Itoa(pageNum)
	params["pageSize"] = strconv.Itoa(pageSize)

	bodyBytes, err := api.get("everything", params)
	if err != nil {
		return nil, errors.New(prefix + " (Api get): " + err.Error())
	}

	var articleResponse ApiArticlesResponse
	if err := json.Unmarshal(bodyBytes, &articleResponse); err != nil {
		return nil, errors.New(prefix + " (json unmarshal): " + err.Error())
	}

	articles := make([]*Article, len(articleResponse.Articles))
	for index, apiArticle := range articleResponse.Articles {
		articles[index] = api.convertArticle(&apiArticle)
	}

	return articles, nil
}

// TODO: json anyway to copy from one structure to another only some values
// basicaly any way to reduce the code here.
func (api *NewsApi) convertSource(apiSource *ApiSource) *Source {
	source := Source{
		ApiSourceName: api.sourceName,
		SourceId:      apiSource.Id,
		Name:          apiSource.Name,
		Description:   apiSource.Description,
		Url:           apiSource.URL,
		Category:      apiSource.Category,
		Language:      apiSource.Language,
		Country:       apiSource.Country,
	}

	return &source
}

func (api *NewsApi) convertArticle(apiArticle *ApiArticle) *Article {
	article := Article{
		ApiSourceName: api.sourceName,
		Author:        apiArticle.Author,
		Title:         apiArticle.Title,
		Description:   apiArticle.Description,
		Url:           apiArticle.URL,
		UrlToImage:    apiArticle.URLToImage,
		PublishedAt:   apiArticle.PublishedAt,
	}

	// TODO: check if there is no source. will that give an invalid json, so a error in previous, or here?
	// can't compare apiArticle.Source with nil
	article.SourceId = apiArticle.Source.Id
	article.SourceName = apiArticle.Source.Name

	return &article
}

// TODO: handle errors better
// if NOT 200 OK, then return a ApiError
func (api *NewsApi) get(endpoint string, params map[string]string) ([]byte, error) {
	prefix := "main.NewsApi.get"
	req, err := http.NewRequest("GET", api.url+"/"+endpoint, nil)
	if err != nil {
		return nil, errors.New(prefix + " (http request): " + err.Error())
	}

	req.Header.Add("X-Api-Key", api.key)

	if params != nil {
		q := req.URL.Query()
		for key, value := range params {
			q.Add(key, value)
		}
		req.URL.RawQuery = q.Encode()
	}

	resp, err := api.client.Do(req)
	if err != nil {
		return nil, errors.New(prefix + " (client do): " + err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(prefix + " (HTTP Status Error): " + resp.Status)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New(prefix + " (read): " + err.Error())
	}

	return bodyBytes, nil
}
