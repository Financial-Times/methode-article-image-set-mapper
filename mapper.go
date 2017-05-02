package main

import (
	"net/http"
	"encoding/json"
)

// Mapper is the main mapper here. I had to comment this line (gometalinter.v1)
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
	defer func() {
		errClose := r.Body.Close()
		if errClose != nil {
			panic(errClose)
		}
	}()
	_, err = w.Write([]byte(native.Value))
	if err != nil {
		panic(err)
	}
}
