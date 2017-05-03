package main

import (
	"encoding/xml"
	"fmt"
)

type XmlMapper struct {}

type Article struct {
	XMLName xml.Name `xml:"doc"`
	Body Body `xml:"story>text>body"`
}

type Body struct {
	ImageSets []ImageSet `xml:"image-set"`
}

type ImageSet struct {
	Id string `xml:"id,attr"`
	ImageSmall Image `xml:"image-small"`
	ImageMedium Image `xml:"image-medium"`
	ImageLarge Image `xml:"image-large"`
}

type Image struct {
	FileRef string `xml:"fileref,attr"`
}

func (m XmlMapper) mapXml(source []byte) ([]ImageSet, error) {
	var article Article
	err := xml.Unmarshal(source, &article)
	if err != nil {
		return nil, fmt.Errorf("Cound't unmarshall native value as XML doucment. %v", err)
	}
	return article.Body.ImageSets, nil
}
