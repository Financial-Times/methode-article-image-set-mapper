package main

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
		jsonImageSet := JSONImageSet{
			UUID:    xmlImageSet.ID,
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
