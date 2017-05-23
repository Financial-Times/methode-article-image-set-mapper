package main

import (
	"encoding/json"
	"fmt"
	trans "github.com/Financial-Times/transactionid-utils-go"
	"github.com/Sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
)

type ErrorMessage struct {
	Message string `json:"message"`
}

func newHTTPErrorMessage(msg string) ErrorMessage {
	return ErrorMessage{Message: msg}
}

type HTTPMappingHandler interface {
	handle(w http.ResponseWriter, r *http.Request)
}

type defaultHTTPMappingHandler struct {
	messageToNativeMapper MessageToNativeMapper
	imageSetMapper        ImageSetMapper
}

func newHTTPMappingHandler(messageToNativeMapper MessageToNativeMapper, imageSetMapper ImageSetMapper) HTTPMappingHandler {
	return defaultHTTPMappingHandler{
		messageToNativeMapper: messageToNativeMapper,
		imageSetMapper:        imageSetMapper,
	}
}

func (h defaultHTTPMappingHandler) handle(w http.ResponseWriter, r *http.Request) {
	tid := trans.GetTransactionIDFromRequest(r)
	w.Header().Add("Content-Type", "application/json;charset=utf-8")
	w.Header().Add(trans.TransactionIDHeader, tid)

	body, err := ioutil.ReadAll(r.Body)
	defer h.closeRequestBody(r)
	if err != nil {
		h.warnAndWriteToHTTP500(fmt.Sprintf("Cound't read from request body. %v\n", err), w)
		return
	}

	native, err := h.messageToNativeMapper.Map(body)
	if err != nil {
		h.writeToHTTP500(fmt.Sprintf("Error mapping native message. %v\n", err), w)
		return
	}

	imageSets, err := h.imageSetMapper.Map(native, time.Now().Format(uppDateFormat), tid)
	if err != nil {
		h.writeToHTTP500(fmt.Sprintf("Error mapping the given content. %v\n", err), w)
		return
	}

	marshaledJSONImageSets, err := json.Marshal(imageSets)
	if err != nil {
		h.warnAndWriteToHTTP500(fmt.Sprintf("Couldn't marshall built-up image-sets to JSON. %v\n", err), w)
		return
	}

	_, err = w.Write(marshaledJSONImageSets)
	if err != nil {
		h.warnAndWriteToHTTP500(fmt.Sprintf("Cound't write response. %v\n", err), w)
	}
}

func (h defaultHTTPMappingHandler) warnAndWriteToHTTP500(msg string, w http.ResponseWriter) {
	logrus.Warn(msg)
	h.writeToHTTP500(msg, w)
}

func (h defaultHTTPMappingHandler) writeToHTTP500(msg string, w http.ResponseWriter) {
	httpMsg, marshalErr := json.Marshal(newHTTPErrorMessage(msg))
	if marshalErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusInternalServerError)
	_, err := w.Write(httpMsg)
	if err != nil {
		logrus.Warn("Couldn't write to response. %v\n", err)
	}
}

func (h defaultHTTPMappingHandler) closeRequestBody(r *http.Request) {
	err := r.Body.Close()
	if err != nil {
		logrus.Warnf("Coulnd't close request body. %v\n", err)
	}
}
