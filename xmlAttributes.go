package main

import "encoding/xml"

type xmlAttributes struct {
	XMLName        xml.Name       `xml:"ObjectMetadata"`
	OutputChannels OutputChannels `xml:"OutputChannels"`
}

type OutputChannels struct {
	DIFTcom DIFTcom `xml:"DIFTcom"`
}

type DIFTcom struct {
	DIFTcomLastPublication    string `xml:"DIFTcomLastPublication"`
	DIFTcomInitialPublication string `xml:"DIFTcomInitialPublication"`
}
