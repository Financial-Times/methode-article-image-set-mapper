package main

import (
	"encoding/base64"
	"fmt"

	uuidutils "github.com/Financial-Times/uuid-utils-go"
	"github.com/sirupsen/logrus"
)

const compoundStory = "EOM::CompoundStory"

type ImageSetMapper interface {
	Map(source NativeContent, lastModified string, publishReference string) ([]JSONImageSet, error)
}

type defaultImageSetMapper struct {
	articleToImageSetMapper ArticleToImageSetMapper
	attributesMapper        AttributesMapper
	xmlImageSetToJSONMapper XMLImageSetToJSONMapper
}

func newImageSetMapper(articleToImageSetMApper ArticleToImageSetMapper, attributesMapper AttributesMapper,
	xmlImageSetToJSONMapper XMLImageSetToJSONMapper) ImageSetMapper {
	return defaultImageSetMapper{
		articleToImageSetMapper: articleToImageSetMApper,
		attributesMapper:        attributesMapper,
		xmlImageSetToJSONMapper: xmlImageSetToJSONMapper,
	}
}

func (m defaultImageSetMapper) Map(source NativeContent, lastModified string, publishReference string) ([]JSONImageSet, error) {
	articleUUID := source.UUID
	err := uuidutils.ValidateUUID(articleUUID)
	if err != nil {
		msg := fmt.Errorf("no valid UUID found in article %v", err)
		logrus.Warn(msg)
		return nil, msg
	}

	valueXML, err := base64.StdEncoding.DecodeString(source.Value)
	if err != nil {
		msg := fmt.Errorf("couldn't decode string as base64 %v", err)
		logrus.Warn(msg)
		return nil, msg
	}

	XMLImageSets, err := m.articleToImageSetMapper.Map(valueXML)
	if err != nil {
		msg := fmt.Errorf("couldn't parse XML document %v", err)
		logrus.Warn(msg)
		return nil, msg
	}

	attributes, err := m.attributesMapper.Map(source.Attributes)
	if err != nil {
		msg := fmt.Errorf("couldn't parse attributes XML %v", err)
		logrus.Warn(msg)
		return nil, msg
	}

	jsonImageSets, err := m.xmlImageSetToJSONMapper.Map(XMLImageSets, articleUUID, attributes, lastModified, publishReference)
	if err != nil {
		msg := fmt.Errorf("couldn't map ImageSets from model sourced from XML to model targeted for JSON. %v", err)
		logrus.Error(msg)
		return nil, msg
	}
	return jsonImageSets, nil
}
