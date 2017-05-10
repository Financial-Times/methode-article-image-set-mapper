package main

import "github.com/Financial-Times/methode-article-image-set-mapper/uuidutils"

type XMLImageSetToJSONMapper interface {
	Map(xmlImageSets []XMLImageSet) ([]JSONImageSet, error)
}

type defaultImageSetToJSONMapper struct{}

func (m defaultImageSetToJSONMapper) Map(xmlImageSets []XMLImageSet) ([]JSONImageSet, error) {
	jsonImageSets := make([]JSONImageSet, 0)
	for _, xmlImageSet := range xmlImageSets {
		members := []JSONMember{
			m.mapMember(xmlImageSet.ImageMedium),
			m.mapMember(xmlImageSet.ImageLarge),
			m.mapMember(xmlImageSet.ImageSmall),
		}
		uuid, err := uuidut123.NewUUIDFromString(xmlImageSet.ID)
		if err != nil {
			return nil, err
		}
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
		UUID: xmlImage.FileRef,
	}
}
