package main

type JsonImageSet struct {
	Uuid string `json:"uuid"`
	Members []JsonMember `json:"members"`
}

type JsonMember struct {
	Uuid string `json:"uuid"`
}
