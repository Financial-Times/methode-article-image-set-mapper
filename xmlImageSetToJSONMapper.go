package main

import (
	"github.com/Financial-Times/uuid-utils-go"
	"strings"
)

type XMLImageSetToJSONMapper interface {
	Map(xmlImageSets []XMLImageSet) ([]JSONImageSet, error)
}

type defaultImageSetToJSONMapper struct{}

func (m defaultImageSetToJSONMapper) Map(xmlImageSets []XMLImageSet) ([]JSONImageSet, error) {
	jsonImageSets := make([]JSONImageSet, 0)
	for _, xmlImageSet := range xmlImageSets {
		members := []JSONMember{
			m.mapMember(xmlImageSet.ImageMedium),
			m.mapMember(xmlImageSet.ImageSmall),
			m.mapMember(xmlImageSet.ImageLarge),
		}
		uuid := uuidutils.NewV3UUID(xmlImageSet.ID)
		jsonImageSet := JSONImageSet{
			UUID: uuid.String(),
			Members: members,
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
