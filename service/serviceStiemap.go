package service

import (
	. "zyg/datamodels"
	"zyg/repositories"
)

func SitemapQuery(index int) ([]ArticleStruct, error) {
	max := 3000
	index = (index - 1) * max
	articleSet, err := repositories.QueryBySort(repositories.DB, index, max)
	if err != nil {
		return nil, err
	}
	return articleSet, nil
}
func SitemapCount() (int64, error) {
	Count, err := repositories.CountArticleAll(repositories.DB)
	if err != nil {
		return Count, err
	}
	page := Count / 3000
	return page, nil
}
