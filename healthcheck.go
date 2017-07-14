package main

import (
	"net/http"

	fthealth "github.com/Financial-Times/go-fthealth/v1_1"
	"github.com/Financial-Times/message-queue-gonsumer/consumer"
	"github.com/Financial-Times/message-queue-go-producer/producer"
	"github.com/Financial-Times/service-status-go/gtg"
)

const healthPath = "/__health"

type HealthCheck struct {
	consumer      consumer.MessageConsumer
	producer      producer.MessageProducer
	appSystemCode string
	appName       string
}

func NewHealthCheck(p producer.MessageProducer, c consumer.MessageConsumer, systemCode, appName string) *HealthCheck {
	return &HealthCheck{
		consumer:      c,
		producer:      p,
		appSystemCode: systemCode,
		appName:       appName,
	}
}

func (h *HealthCheck) Health() func(w http.ResponseWriter, r *http.Request) {
	checks := []fthealth.Check{h.readQueueCheck(), h.writeQueueCheck()}
	hc := fthealth.HealthCheck{
		SystemCode:  h.appSystemCode,
		Name:        h.appName,
		Description: "Maps inline image-sets from bodies of Methode articles.",
		Checks:      checks,
	}
	return fthealth.Handler(hc)
}

func (h *HealthCheck) readQueueCheck() fthealth.Check {
	return fthealth.Check{
		ID:               "read-message-queue-proxy-reachable",
		Name:             "Read Message Queue Proxy Reachable",
		Severity:         1,
		BusinessImpact:   "Publishing or updating image-sets inside article bodies will not be possible, clients will not see images or image-sets in new content.",
		TechnicalSummary: "Read message queue proxy is not reachable/healthy",
		PanicGuide:       "https://dewey.ft.com/methode-article-image-set-mapper.html",
		Checker:          h.consumer.ConnectivityCheck,
	}
}

func (h *HealthCheck) writeQueueCheck() fthealth.Check {
	return fthealth.Check{
		ID:               "write-message-queue-proxy-reachable",
		Name:             "Write Message Queue Proxy Reachable",
		Severity:         1,
		BusinessImpact:   "Publishing or updating image-sets inside article bodies will not be possible, clients will not see images or image-sets in new content.",
		TechnicalSummary: "Write message queue proxy is not reachable/healthy",
		PanicGuide:       "https://dewey.ft.com/methode-article-image-set-mapper.html",
		Checker:          h.producer.ConnectivityCheck,
	}
}

func (h *HealthCheck) GTG() gtg.Status {
	consumerCheck := func() gtg.Status {
		return gtgCheck(h.consumer.ConnectivityCheck)
	}
	producerCheck := func() gtg.Status {
		return gtgCheck(h.producer.ConnectivityCheck)
	}

	return gtg.FailFastParallelCheck([]gtg.StatusChecker{
		consumerCheck,
		producerCheck,
	})()
}

func gtgCheck(handler func() (string, error)) gtg.Status {
	if _, err := handler(); err != nil {
		return gtg.Status{GoodToGo: false, Message: err.Error()}
	}
	return gtg.Status{GoodToGo: true}
}
