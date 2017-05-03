package main

type XmlImageSetToJsonMapper interface {
	Map(xmlImageSets []XmlImageSet) ([]JsonImageSet, error)
}

type defaultImageSetToJsonMapper struct {}

func (m defaultImageSetToJsonMapper) Map(xmlImageSets []XmlImageSet) ([]JsonImageSet, error) {
	jsonImageSets := make([]JsonImageSet, 0)
	for _, xmlImageSet := range xmlImageSets {
		members := []JsonMember{
			m.mapMember(xmlImageSet.ImageMedium),
			m.mapMember(xmlImageSet.ImageLarge),
			m.mapMember(xmlImageSet.ImageSmall),
		}
		jsonImageSet := JsonImageSet {
			Uuid: xmlImageSet.Id,
			Members: members,
		}
		jsonImageSets = append(jsonImageSets, jsonImageSet)
	}
	return jsonImageSets, nil
}

func (m defaultImageSetToJsonMapper) mapMember(xmlImage XmlImage) JsonMember {
	return JsonMember{
		Uuid: xmlImage.FileRef,
	}
}
