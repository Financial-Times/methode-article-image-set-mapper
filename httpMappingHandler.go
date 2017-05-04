package main

import (
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type ErrorMessage struct {
	Message string `json:"message"`
}

func NewHTTPErrorMessage(msg string) ErrorMessage {
	return ErrorMessage{Message: msg}
}

type HTTPMappingHandler struct {
	imageSetMapper ImageSetMapper
}

func newHTTPMappingHandler(imageSetMapper ImageSetMapper) HTTPMappingHandler {
	return HTTPMappingHandler{
		imageSetMapper: imageSetMapper,
	}
}

func (h HTTPMappingHandler) handle(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	defer h.closeResponseBody(r)
	if err != nil {
		h.warnAndWriteToHTTP500(fmt.Sprintf("Cound't read from request body. %v\n", err), w)
		return
	}

	marshaledJSONImageSets, err := h.imageSetMapper.MapToJson(body)
	if err != nil {
		h.writeToHTTP500(fmt.Sprintf("Error mapping the given content. %v\n", err), w)
		return
	}

	w.Header().Add("Content-Type", "application/json;charset=utf-8")
	_, err = w.Write(marshaledJSONImageSets)
	if err != nil {
		h.warnAndWriteToHTTP500(fmt.Sprintf("Cound't write response. %v\n", err), w)
	}
}

func (h HTTPMappingHandler) warnAndWriteToHTTP500(msg string, w http.ResponseWriter) {
	logrus.Warn(msg)
	h.writeToHTTP500(msg, w)
}

func (h HTTPMappingHandler) writeToHTTP500(msg string, w http.ResponseWriter) {
	httpMsg, marshalErr := json.Marshal(NewHTTPErrorMessage(msg))
	if marshalErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	_, err := w.Write(httpMsg)
	if err != nil {
		logrus.Warn("Couldn't write to response. %v\n", err)
	}
}

func (h HTTPMappingHandler) closeResponseBody(r *http.Request) {
	err := r.Body.Close()
	if err != nil {
		logrus.Warnf("Coulnd't close request body. %v\n", err)
	}
}
