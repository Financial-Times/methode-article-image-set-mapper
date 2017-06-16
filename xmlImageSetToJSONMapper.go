package main

import (
	"fmt"
	"github.com/Financial-Times/uuid-utils-go"
	"github.com/Sirupsen/logrus"
	"strings"
	"time"
)

const (
	methodeAuthority  = "http://api.ft.com/system/FTCOM-METHODE"
	canBeDistributedYes            = "yes"
	methodeDateFormat = "20060102150405"
	uppDateFormat     = "2006-01-02T15:04:05.000Z0700"
	imageSetType      = "ImageSet"
)

type XMLImageSetToJSONMapper interface {
	Map(xmlImageSets []XMLImageSet, attributes xmlAttributes, lastModified string, publishReference string) ([]JSONImageSet, error)
}

type defaultImageSetToJSONMapper struct{}

func (m defaultImageSetToJSONMapper) Map(xmlImageSets []XMLImageSet, attributes xmlAttributes, lastModified string, publishReference string) ([]JSONImageSet, error) {
	jsonImageSets := make([]JSONImageSet, 0)
	for _, xmlImageSet := range xmlImageSets {
		members := make([]JSONMember, 0, 3)
		m.appendIfPresent(&members, xmlImageSet.ImageMedium, "medium", "", "")
		m.appendIfPresent(&members, xmlImageSet.ImageSmall, "small", "490px", "")
		m.appendIfPresent(&members, xmlImageSet.ImageLarge, "large", "", "980px")

		uuid := uuidutils.NewV3UUID(xmlImageSet.ID)
		publishedDate, err := time.Parse(methodeDateFormat, attributes.OutputChannels.DIFTcom.DIFTcomLastPublication)
		if err != nil {
			return nil, fmt.Errorf("Couldn't parse required methode field published date (DIFTcomLastPublication) %v %v", attributes.OutputChannels.DIFTcom.DIFTcomLastPublication, err)
		}
		firstPublishedDate, err := time.Parse(methodeDateFormat, attributes.OutputChannels.DIFTcom.DIFTcomInitialPublication)
		if err != nil {
			return nil, fmt.Errorf("Couldn't parse required methode field initial published date (DIFTcomInitialPublication) %v %v", attributes.OutputChannels.DIFTcom.DIFTcomInitialPublication, err)
		}
		jsonImageSet := JSONImageSet{
			UUID:    uuid.String(),
			Members: members,
			Identifiers: []JSONIdentifier{
				JSONIdentifier{
					Authority:       methodeAuthority,
					IdentifierValue: uuid.String(),
				},
			},
			PublishedDate:      publishedDate.Format(uppDateFormat),
			FirstPublishedDate: firstPublishedDate.Format(uppDateFormat),
			CanBeDistributed:   canBeDistributedYes,
			LastModified:       lastModified,
			PublishReference:   publishReference,
			Type:               imageSetType,
		}
		jsonImageSets = append(jsonImageSets, jsonImageSet)
	}
	return jsonImageSets, nil
}

func (m defaultImageSetToJSONMapper) appendIfPresent(members *[]JSONMember, xmlImage XMLImage, memberName string, maxDisplayWidth string, minDisplayWidth string) {
	jsonMember := m.mapMember(xmlImage, memberName, maxDisplayWidth, minDisplayWidth)
	if jsonMember == nil {
		logrus.Warn("")
	} else {
		*members = append(*members, *jsonMember)
	}
}

func (m defaultImageSetToJSONMapper) mapMember(xmlImage XMLImage, memberName string, maxDisplayWidth string, minDisplayWidth string) *JSONMember {
	if xmlImage.FileRef == "" {
		logrus.Warnf("expected member %v is not present.", memberName)
		return nil
	}
	refs := strings.Split(xmlImage.FileRef, "?uuid=")
	if len(refs) != 2 {
		logrus.Warnf("at member %v fileref attribute doesn't contain uuid fileref=%v", memberName, xmlImage.FileRef)
		return nil
	}
	return &JSONMember{
		UUID: strings.Split(xmlImage.FileRef, "?uuid=")[1],
		MaxDisplayWidth: maxDisplayWidth,
		MinDisplayWidth: minDisplayWidth,
	}
}
