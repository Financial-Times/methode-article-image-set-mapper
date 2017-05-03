package main

import "encoding/xml"

type xmlArticle struct {
	XMLName xml.Name `xml:"doc"`
	Body    xmlBody `xml:"story>text>body"`
}

type xmlBody struct {
	ImageSets []XmlImageSet `xml:"image-set"`
}

type XmlImageSet struct {
	Id          string `xml:"id,attr"`
	ImageSmall  XmlImage `xml:"image-small"`
	ImageMedium XmlImage `xml:"image-medium"`
	ImageLarge  XmlImage `xml:"image-large"`
}

type XmlImage struct {
	FileRef string `xml:"fileref,attr"`
}
