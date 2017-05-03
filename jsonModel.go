package main

type JSONImageSet struct {
	UUID    string `json:"uuid"`
	Members []JSONMember `json:"members"`
}

type JSONMember struct {
	UUID string `json:"uuid"`
}
