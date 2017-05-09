package main

import (
	"encoding/base64"
	"fmt"
	"github.com/Sirupsen/logrus"
)

const compoundStory = "EOM::CompoundStory"

type ImageSetMapper interface {
	Map(source NativeContent) ([]JSONImageSet, error)
}

type defaultImageSetMapper struct {
	articleToImageSetMapper ArticleToImageSetMapper
	xmlImageSetToJSONMapper XMLImageSetToJSONMapper
}

func newImageSetMapper(articleToImageSetMApper ArticleToImageSetMapper, xmlImageSetToJSONMapper XMLImageSetToJSONMapper) ImageSetMapper {
	return defaultImageSetMapper{
		articleToImageSetMapper: articleToImageSetMApper,
		xmlImageSetToJSONMapper: xmlImageSetToJSONMapper,
	}
}

func (m defaultImageSetMapper) Map(source NativeContent) ([]JSONImageSet, error) {
	xmlDocument, err := base64.StdEncoding.DecodeString(source.Value)
	if err != nil {
		msg := fmt.Errorf("Cound't decode string as base64. %v\n", err)
		logrus.Warn(msg)
		return nil, msg
	}

	xmlImageSets, err := m.articleToImageSetMapper.Map(xmlDocument)
	if err != nil {
		msg := fmt.Errorf("Couldn't parse XML document. %v\n", err)
		logrus.Warn(msg)
		return nil, msg
	}

	jsonImageSets, err := m.xmlImageSetToJSONMapper.Map(xmlImageSets)
	if err != nil {
		msg := fmt.Errorf("Couldn't map ImageSets from model soruced from XML to model targeted for JSON. %v\n", err)
		logrus.Warn(msg)
		return nil, msg
	}
	return jsonImageSets, nil
}
