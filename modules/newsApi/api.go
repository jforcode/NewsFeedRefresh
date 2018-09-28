package newsApi

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Api struct {
	url    string
	key    string
	client *http.Client
}

func (api *Api) Init(url, key string) error {
	api.url = url
	api.key = key
	api.client = &http.Client{}

	return nil
}

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

func (api *Api) fetchSources() (*ApiSourcesResponse, error) {
	bodyBytes, err := api.get("sources", nil)
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

	return &articleResponse, nil
}
