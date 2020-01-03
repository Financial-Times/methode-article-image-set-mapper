package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Financial-Times/message-queue-go-producer/producer"
	consumer "github.com/Financial-Times/message-queue-gonsumer"
	trans "github.com/Financial-Times/transactionid-utils-go"
	"github.com/sirupsen/logrus"
	gouuid "github.com/satori/go.uuid"
	"sync"
	"time"
)

const (
	methodeSystemOrigin = "http://cmdb.ft.com/systems/methode-web-pub"
	dateFormat          = "2006-01-02T15:04:05.000Z0700"
	contentURIBase      = "http://methode-article-image-set-mapper.svc.ft.com/image-set/model/"
)

type queue interface {
	onMessage(m consumer.Message)
}

type defaultQueue struct {
	messageConsumer   consumer.MessageConsumer
	messageProducer   producer.MessageProducer
	consumerWaitGroup sync.WaitGroup

	messageToNativeMapper MessageToNativeMapper
	imageSetMapper        ImageSetMapper
}

func newQueue(messageConsumer consumer.MessageConsumer, messageProducer producer.MessageProducer,
	messageToNativeMapper MessageToNativeMapper, imageSetMapper ImageSetMapper) defaultQueue {
	queue := defaultQueue{
		messageConsumer:       messageConsumer,
		messageProducer:       messageProducer,
		messageToNativeMapper: messageToNativeMapper,
		imageSetMapper:        imageSetMapper,
		consumerWaitGroup:     sync.WaitGroup{},
	}
	return queue
}

func (q defaultQueue) onMessage(m consumer.Message) {
	tid := m.Headers[trans.TransactionIDHeader]
	if tid == "" {
		tid = trans.NewTransactionID()
		logrus.Warnf("X-Request-Id not found in kafka message headers. Created now. transactionId=%v", tid)
	}

	if m.Headers["Origin-System-Id"] != methodeSystemOrigin {
		logrus.Infof("Ignoring message with different originSystemId=%v transactionId=%v ", m.Headers["Origin-System-Id"], tid)
		return
	}

	lastModified := m.Headers["Message-Timestamp"]
	if lastModified == "" {
		lastModified = time.Now().Format(dateFormat)
		logrus.Infof("Last modified date was empty on message, created now. transactionId=%v lastModifiedDate=%v", tid, lastModified)
	}

	native, err := q.messageToNativeMapper.Map([]byte(m.Body))
	if err != nil {
		logrus.Errorf("Error mapping native message. transactionId=%v %v", tid, err)
		return
	}
	if native.Type != compoundStory {
		logrus.Infof("Ignoring message of type=%v transactionId=%v", native.Type, tid)
		return
	}

	imageSets, err := q.imageSetMapper.Map(native, lastModified, tid)
	if err != nil {
		logrus.Errorf("Error mapping message to image-sets transactionId=%v %v", tid, err)
		return
	}

	if len(imageSets) == 0 {
		logrus.Infof("No image-sets were found in this article. transactionId=%v", tid)
	} else {
		msgs, errs := q.buildMessages(imageSets, lastModified, tid)
		if len(errs) != 0 {
			for uuid, err := range errs {
				logrus.Errorf("Couldn't build message for image-set transactionId=%v uuid=%v %v", tid, uuid, err)
			}
		}
		for uuid, msg := range msgs {
			err = q.messageProducer.SendMessage("", msg)
			if err != nil {
				logrus.Errorf("Error sending transformed message to queue transactionId=%v uuid=%v %v", tid, uuid, err)
				continue
			}
			logrus.Infof("Mapped and sent for uuid=%v transactionId=%v", uuid, tid)
		}
	}
}

func (q defaultQueue) buildMessages(imageSets []JSONImageSet, lastModified string, tid string) (map[string]producer.Message, map[string]error) {
	errs := make(map[string]error, 0)
	msgs := make(map[string]producer.Message, 0)
	for _, imageSet := range imageSets {
		msg, err := q.buildMessage(imageSet, lastModified, tid)
		if err != nil {
			errs[imageSet.UUID] = err
			continue
		}
		msgs[imageSet.UUID] = msg
	}
	return msgs, errs
}

func (q defaultQueue) buildMessage(imageSet JSONImageSet, lastModified, pubRef string) (producer.Message, error) {
	headers := map[string]string{
		"X-Request-Id":      pubRef,
		"Message-Timestamp": lastModified,
		"Message-Id":        gouuid.NewV4().String(),
		"Message-Type":      "cms-content-published",
		"Content-Type":      "application/json",
		"Origin-System-Id":  methodeSystemOrigin,
	}
	body := publicationMessageBody{
		ContentURI:   contentURIBase + imageSet.UUID,
		Payload:      imageSet,
		LastModified: lastModified,
	}
	marshaledBody, err := q.unsafeJSONMarshal(body)
	if err != nil {
		return producer.Message{}, fmt.Errorf("Couldn't marshall message body to JSON skipping message. transactionId=%v %v", pubRef, body)
	}
	return producer.Message{Headers: headers, Body: string(marshaledBody)}, nil
}

func (q defaultQueue) unsafeJSONMarshal(v interface{}) ([]byte, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	b = bytes.Replace(b, []byte("\\u003c"), []byte("<"), -1)
	b = bytes.Replace(b, []byte("\\u003e"), []byte(">"), -1)
	return b, nil
}

func (q defaultQueue) startConsuming() {
	q.consumerWaitGroup.Add(1)
	go func() {
		q.messageConsumer.Start()
		q.consumerWaitGroup.Done()
	}()
}

func (q defaultQueue) stop() {
	q.messageConsumer.Stop()
	q.consumerWaitGroup.Wait()
}
