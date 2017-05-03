package main

import "encoding/xml"

type xmlArticle struct {
	XMLName xml.Name `xml:"doc"`
	Body    xmlBody `xml:"story>text>body"`
}

type xmlBody struct {
	ImageSets []XMLImageSet `xml:"image-set"`
}

type XMLImageSet struct {
	ID          string `xml:"id,attr"`
	ImageSmall  XMLImage `xml:"image-small"`
	ImageMedium XMLImage `xml:"image-medium"`
	ImageLarge  XMLImage `xml:"image-large"`
}

type XMLImage struct {
	FileRef string `xml:"fileref,attr"`
}
