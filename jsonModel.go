package main

type JSONImageSet struct {
	UUID               string           `json:"uuid"`
	Identifiers        []JSONIdentifier `json:"identifiers"`
	Members            []JSONMember     `json:"members"`
	PublishReference   string           `json:"publishReference"`
	LastModified       string           `json:"lastModified"`
	PublishedDate      string           `json:"publishedDate"`
	FirstPublishedDate string           `json:"firstPublishedDate"`
	CanBeDistributed   string           `json:"canBeDistributed"`
	Type               string           `json:"type"`
}

type JSONMember struct {
	UUID string `json:"uuid"`
}

type JSONIdentifier struct {
	Authority       string `json:"authority"`
	IdentifierValue string `json:"identifierValue"`
}
