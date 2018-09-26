package newsApi

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Api struct {
	url string
	key string
}

func (api *Api) Init(url, key string) {
	api.url = url
	api.key = key
}

func (api *Api) fetchSources() (*ApiSourcesResponse, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", api.url+"/sources", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-Api-Key", api.key)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("HTTP Status Error: " + resp.Status)
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var sourceResponse ApiSourcesResponse
	if err := json.Unmarshal(bodyBytes, &sourceResponse); err != nil {
		return nil, err
	}

	return &sourceResponse, nil
}

func (api *Api) fetchArticles(sourceIds string, pageNum, pageSize int) (*ApiArticlesResponse, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", api.url+"/everything", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-Api-Key", api.key)
	q := req.URL.Query()
	q.Add("sources", sourceIds)
	q.Add("page", strconv.Itoa(pageNum))
	q.Add("pageSize", strconv.Itoa(pageSize))
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("HTTP Status Error: " + resp.Status)
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var articleResponse ApiArticlesResponse
	if err := json.Unmarshal(bodyBytes, &articleResponse); err != nil {
		return nil, err
	}

	return &articleResponse, nil
}
