package main

type xmlAttributes struct {
	ObjectMetadata ObjectMetadata `xml:"ObjectMetadata"`
}

type ObjectMetadata struct {
	OutputChannels OutputChannels `xml:"OutputChannels"`
}

type OutputChannels struct {
	DIFTcom DIFTcom `xml:"DIFTcom"`
}

type DIFTcom struct {
	DIFTcomLastPublication string `xml:"DIFTcomLastPublication"`
	DIFTcomInitialPublication string `xml:"DIFTcomInitialPublication"`
}
