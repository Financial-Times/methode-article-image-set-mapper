package main

import (
	"net/http"
	"encoding/json"
	"github.com/Sirupsen/logrus"
	"encoding/base64"
)

// Mapper is the main mapper here. I had to comment this line (gometalinter.v1)
type Mapper struct {
	XmlMapper XmlMapper
	XmlToJson XmlModelToJsonModelMapper
}

func newMapper() Mapper {
	return Mapper{
		XmlMapper: XmlMapper{},
		XmlToJson: XmlModelToJsonModelMapper{},
	}
}

func (m Mapper) mapArticleImageSets(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	type nativeContent struct {
		Value string `json:"value"`
	}
	var native nativeContent
	err := decoder.Decode(&native)
	if err != nil {
		logrus.Warnf("Cound't decode native content as JSON doucment. %v\n", err)
	}
	defer closeResponseBody(r)

	xmlDocument, err := base64.StdEncoding.DecodeString(native.Value)
	if err != nil {
		logrus.Warnf("Cound't decode string as base64. %v\n", err)
	}

	xmlImageSets, err := m.XmlMapper.mapXml(xmlDocument)
	if err != nil {
		logrus.Warnf("Couldn't map ImageSets. %v\n", err)
	}

	jsonImageSets, err := m.XmlToJson.mapp(xmlImageSets)

	marshaledJsonImageSets, err := json.Marshal(jsonImageSets)
	if err != nil {
		logrus.Warnf("Couldn't marshall built-up image-sets to JSON. %v\n", err)
	}

	_, err = w.Write(marshaledJsonImageSets)
	if err != nil {
		logrus.Warnf("Cound't write response. %v\n", err)
	}
}

func closeResponseBody(r *http.Request) {
	err := r.Body.Close()
	if err != nil {
		logrus.Warnf("Cound't close request body. %v\n", err)
	}
}
