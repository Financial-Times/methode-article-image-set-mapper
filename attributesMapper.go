package main

import (
	"fmt"
	"encoding/xml"
	"strings"
)

type AttributesMapper interface {
	Map(source string) (xmlAttributes, error)
}

type defaultAttributesMapper struct {}

func (m defaultAttributesMapper) Map(source string) (xmlAttributes, error) {
	strippedNewline := strings.Replace(source, `\n`, "", -1)
	manualUnqouted := strings.Replace(strippedNewline, `\"`, `"`, -1)

	var attributes xmlAttributes
	err := xml.Unmarshal([]byte(manualUnqouted), &attributes)
	if err != nil {
		return xmlAttributes{}, fmt.Errorf("Cound't unmarshall native attributes as XML document. %v", err)
	}
	return attributes, nil
}
