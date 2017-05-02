package main

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"io"
)

type Mapper struct {}

func (m Mapper) mapArticleImageSets(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	type nativeContent struct {
		Value string `json:"value"`
	}
	var native nativeContent
	err := decoder.Decode(&native)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	w.Write([]byte(native.Value))
}

func nicely(resp *http.Response) {
	_, err := io.Copy(ioutil.Discard, resp.Body)
	if err != nil {
		//log.Warningf("[%v]", err)
	}
}
