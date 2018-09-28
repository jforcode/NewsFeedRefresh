package converter

import (
	"github.com/jforcode/NewsFeedRefresh/modules/feedSrv"
	"github.com/jforcode/NewsFeedRefresh/modules/newsApi"
)

func convertArticle(*newsApi.ApiArticle) (*feedSrv.Article, error) {
	return nil, nil
}

func convertSource(*newsApi.ApiSource) (*feedSrv.Source, error) {
	return nil, nil
}
