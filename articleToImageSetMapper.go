package main

import (
	"encoding/xml"
	"fmt"
)

type ArticleToImageSetMapper interface {
	Map(source []byte) ([]XMLImageSet, error)
}

type defaultArticleToImageSetMapper struct {}

func (m defaultArticleToImageSetMapper) Map(source []byte) ([]XMLImageSet, error) {
	var article xmlArticle
	err := xml.Unmarshal(source, &article)
	if err != nil {
		return nil, fmt.Errorf("Cound't unmarshall native value as XML doucment. %v", err)
	}
	return article.Body.ImageSets, nil
}
