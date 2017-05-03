package main

type JsonImageSet struct {
	Uuid string `json:"uuid"`
	Members []Member `json:"members"`
}

type Member struct {
	Uuid string `json:"uuid"`
}

type XmlModelToJsonModelMapper struct {
}

func (m XmlModelToJsonModelMapper) mapp(xmlImageSets []ImageSet) ([]JsonImageSet, error) {
	jsonImageSets := make([]JsonImageSet, 0)
	for _, xmlImageSet := range xmlImageSets {
		members := []Member {
			m.mapMember(xmlImageSet.ImageLarge),
			m.mapMember(xmlImageSet.ImageMedium),
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

func (m XmlModelToJsonModelMapper) mapMember(xmlImage Image) Member {
	return Member{
		Uuid: xmlImage.FileRef,
	}
}
