package main

import (
	"net/http"
	"io/ioutil"
	"fmt"
	"github.com/Sirupsen/logrus"
	"encoding/json"
)

type ErrorMessage struct {
	Message string `json:"message"`
}

func NewHttpErrorMessage(msg string) ErrorMessage {
	return ErrorMessage{Message: msg}
}

type HttpMappingHandler struct {
	imageSetMapper ImageSetMapper
}

func newHttpMappingHandler(imageSetMapper ImageSetMapper) HttpMappingHandler {
	return HttpMappingHandler{
		imageSetMapper: imageSetMapper,
	}
}

func (h HttpMappingHandler) handle(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	defer h.closeResponseBody(r)
	if err != nil {
		h.warnAndWriteToHttp500(fmt.Sprintf("Cound't read from request body. %v\n", err), w)
		return
	}

	marshaledJsonImageSets, err := h.imageSetMapper.Map(body)
	if err != nil {
		h.writeToHttp500(fmt.Sprintf("Error mapping the given content. %v\n", err), w)
		return
	}

	w.Header().Add("Content-Type", "application/json;charset=utf-8")
	_, err = w.Write(marshaledJsonImageSets)
	if err != nil {
		h.warnAndWriteToHttp500(fmt.Sprintf("Cound't write response. %v\n", err), w)
	}
}

func (h HttpMappingHandler) warnAndWriteToHttp500(msg string, w http.ResponseWriter) {
	logrus.Warn(msg)
	h.writeToHttp500(msg, w)
}

func (h HttpMappingHandler) writeToHttp500(msg string, w http.ResponseWriter) {
	httpMsg, marshalErr := json.Marshal(NewHttpErrorMessage(msg))
	if marshalErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(httpMsg)
}

func (h HttpMappingHandler) closeResponseBody(r *http.Request) {
	err := r.Body.Close()
	if err != nil {
		logrus.Warnf("Cound't close request body. %v\n", err)
	}
}
