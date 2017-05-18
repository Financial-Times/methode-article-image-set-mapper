package main

import (
	"github.com/Financial-Times/uuid-utils-go"
	"strings"
	"time"
	"fmt"
)

const (
	methodeAuthority = "http://api.ft.com/system/FTCOM-METHODE"
	verify = "verify"
	methodeDateFormat = "20060102150405"
)

type XMLImageSetToJSONMapper interface {
	Map(xmlImageSets []XMLImageSet, attributes xmlAttributes) ([]JSONImageSet, error)
}

type defaultImageSetToJSONMapper struct{}

func (m defaultImageSetToJSONMapper) Map(xmlImageSets []XMLImageSet, attributes xmlAttributes) ([]JSONImageSet, error) {
	jsonImageSets := make([]JSONImageSet, 0)
	for _, xmlImageSet := range xmlImageSets {
		members := []JSONMember{
			m.mapMember(xmlImageSet.ImageMedium),
			m.mapMember(xmlImageSet.ImageSmall),
			m.mapMember(xmlImageSet.ImageLarge),
		}
		uuid := uuidutils.NewV3UUID(xmlImageSet.ID)
		publishedDate, err := time.Parse(methodeDateFormat, attributes.OutputChannels.DIFTcom.DIFTcomLastPublication)
		if err != nil {
			return nil, fmt.Errorf("Couldn't parse native published date %v %v", attributes.OutputChannels.DIFTcom.DIFTcomLastPublication, err)
		}
		firstPublishedDate, err := time.Parse(methodeDateFormat, attributes.OutputChannels.DIFTcom.DIFTcomInitialPublication)
		if err != nil {
			return nil, fmt.Errorf("Couldn't parse native initial published date %v %v", attributes.OutputChannels.DIFTcom.DIFTcomInitialPublication, err)
		}
		jsonImageSet := JSONImageSet{
			UUID: uuid.String(),
			Members: members,
			Identifiers: []JSONIdentifier{
				JSONIdentifier{
					Authority: methodeAuthority,
					IdentifierValue: uuid.String(),
				},
			},
			PublishedDate: publishedDate.Format(time.RFC3339Nano),
			FirstPublishedDate: firstPublishedDate.Format(time.RFC3339Nano),
			CanBeDistributed: verify,
		}
		jsonImageSets = append(jsonImageSets, jsonImageSet)
	}
	return jsonImageSets, nil
}

func (m defaultImageSetToJSONMapper) mapMember(xmlImage XMLImage) JSONMember {
	return JSONMember{
		UUID: strings.Split(xmlImage.FileRef, "?uuid=")[1],
	}
}
