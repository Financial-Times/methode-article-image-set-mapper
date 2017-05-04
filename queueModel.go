package main

type publicationMessageBody struct {
	ContentURI   string         `json:"contentUri"`
	Payload      JSONImageSet `json:"payload"`
	LastModified string         `json:"lastModified"`
}
