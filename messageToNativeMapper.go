package main

import (
	"encoding/json"
	"fmt"
)

type MessageToNativeMapper interface {
	Map(source []byte) (NativeContent, error)
}

type defaultMessageToNativeMapper struct{}

type NativeContent struct {
	Uuid       string `json:"uuid"`
	Type       string `json:"type"`
	Value      string `json:"value"`
	Attributes string `json:"attributes"`
}

func (m defaultMessageToNativeMapper) Map(source []byte) (NativeContent, error) {
	var native NativeContent
	err := json.Unmarshal(source, &native)
	if err != nil {
		return NativeContent{}, fmt.Errorf("Cound't decode native content as JSON document. %v\n", err)
	}
	return native, nil
}
