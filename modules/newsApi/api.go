package newsApi

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/jforcode/NewsFeedRefresh/modules/common"
)

type Api struct {
	url        string
	key        string
	client     *http.Client
	sourceName string
}

func (api *Api) Init(url, key string) error {
	api.url = url
	api.key = key
	api.client = &http.Client{}
	api.sourceName = "news_api"

	return nil
}

func (api *Api) FetchSources() ([](*common.Source), error) {
	bodyBytes, err := api.get("sources", nil)
	if err != nil {
		return nil, err
	}

	var sourceResponse ApiSourcesResponse
	if err := json.Unmarshal(bodyBytes, &sourceResponse); err != nil {
		return nil, err
	}

	sources := make([](*common.Source), len(sourceResponse.Sources))
	for index, apiSource := range sourceResponse.Sources {
		sources[index] = api.convertSource(&apiSource)
	}

	return sources, nil
}

func (api *Api) FetchArticles(sourceIds string, pageNum, pageSize int) ([](*common.Article), error) {
	params := make(map[string]string)
	params["sources"] = sourceIds
	params["page"] = strconv.Itoa(pageNum)
	params["pageSize"] = strconv.Itoa(pageSize)

	bodyBytes, err := api.get("everything", params)
	if err != nil {
		return nil, err
	}

	var articleResponse ApiArticlesResponse
	if err := json.Unmarshal(bodyBytes, &articleResponse); err != nil {
		return nil, err
	}

	articles := make([](*common.Article), len(articleResponse.Articles))
	for index, apiArticle := range articleResponse.Articles {
		articles[index] = api.convertArticle(&apiArticle)
	}

	return articles, nil
}

// TODO: json anyway to copy from one structure to another only some values
// basicaly any way to reduce the code here.
func (api *Api) convertSource(apiSource *ApiSource) *common.Source {
	source := common.Source{
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

func (api *Api) convertArticle(apiArticle *ApiArticle) *common.Article {
	article := common.Article{
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
func (api *Api) get(endpoint string, params map[string]string) ([]byte, error) {
	req, err := http.NewRequest("GET", api.url+"/"+endpoint, nil)
	if err != nil {
		return nil, err
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
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("HTTP Status Error: " + resp.Status)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return bodyBytes, nil
}
