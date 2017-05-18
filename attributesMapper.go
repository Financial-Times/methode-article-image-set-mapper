package main

import (
	"strconv"
	"fmt"
	"encoding/xml"
)

type AttributesMapper interface {
	Map(source string) (xmlAttributes, error)
}

type defaultAttributesMapper struct {}

func (m defaultAttributesMapper) Map(source string) (xmlAttributes, error) {
	unquotedSource, err := strconv.Unquote(source)
	if err != nil {
		return  xmlAttributes{}, fmt.Errorf("Couldn't unqoute attributes %v", err)
	}

	var attributes xmlAttributes
	err = xml.Unmarshal([]byte(unquotedSource), &attributes)
	if err != nil {
		return xmlAttributes{}, fmt.Errorf("Cound't unmarshall native attributes as XML doucment. %v", err)
	}
	return attributes, nil
}
