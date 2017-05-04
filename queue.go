package main

import (
	"github.com/Financial-Times/message-queue-go-producer/producer"
	"github.com/Financial-Times/message-queue-gonsumer/consumer"
	gouuid "github.com/satori/go.uuid"
	"github.com/Sirupsen/logrus"
	"net/http"
	"time"
	"encoding/json"
	"bytes"
	"fmt"
)

const (
	methodeSystemOrigin = "http://cmdb.ft.com/systems/methode-webpub"
	dateFormat          = "2006-01-02T03:04:05.000Z0700"
	contentURIBase = "http://methode-article-images-set-mapper.svc.ft.com/image-set/model/"
)

type queue struct {
	consumerConfig consumer.QueueConfig
	producerConfig producer.MessageProducerConfig
	messageConsumer consumer.MessageConsumer
	messageProducer producer.MessageProducer
	imageSetMapper ImageSetMapper
}

func newQueue(args args, httpClient *http.Client) queue {
	queue := queue{
		consumerConfig: consumer.QueueConfig{
			Addrs:                args.addresses,
			Group:                args.group,
			Topic:                args.readTopic,
			Queue:                args.readQueue,
			ConcurrentProcessing: false,
			AutoCommitEnable:     true,
			AuthorizationKey:     args.authorization,
		},
		producerConfig: producer.MessageProducerConfig{
			Addr:          args.addresses[0],
			Topic:         args.writeTopic,
			Queue:         args.writeQueue,
			Authorization: args.authorization,
		},
	}
	messageConsumer := consumer.NewConsumer(queue.consumerConfig, queue.onMessage, httpClient)
	queue.messageConsumer = messageConsumer
	messageProducer := producer.NewMessageProducerWithHTTPClient(queue.producerConfig, httpClient)
	queue.messageProducer = messageProducer
	return queue
}

func (q queue) onMessage(m consumer.Message) {
	tid := m.Headers["X-Request-Id"]
	if tid == "" {
		logrus.Warnf("X-Request-Id not found in kafka message headers. Skipping message")
		return
	}

	if m.Headers["Origin-System-Id"] != methodeSystemOrigin {
		logrus.Infof("Ignoring message with different originSystemId=%v transactionId=%v ", m.Headers["Origin-System-Id"], tid)
		return
	}

	lastModified := m.Headers["Message-Timestamp"]
	if lastModified == "" {
		lastModified = time.Now().Format(dateFormat)
	}

	imageSets, err := q.imageSetMapper.Map([]byte(m.Body))
	if err != nil {
		logrus.Errorf("Error mapping message to image-sets transactionId=%v %v", tid, err)
		return
	}

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
			return
		}
		logrus.Infof("Mapped and sent for uuid=%v transactionId=%v", uuid, tid)
	}
}

func (q queue) buildMessages(imageSets []JSONImageSet, lastModified string , tid string) (map[string]producer.Message, map[string]error) {
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

func (q queue) buildMessage(imageSet JSONImageSet, lastModified, pubRef string) (producer.Message, error) {
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
	marshaledBody, err := unsafeJSONMarshal(body)
	if err != nil {
		return producer.Message{}, fmt.Errorf("Couldn't marshall message body to JSON skipping message. transactionId=%v %v", pubRef, body)
	}
	return producer.Message{Headers: headers, Body: string(marshaledBody)}, nil
}

func unsafeJSONMarshal(v interface{}) ([]byte, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	b = bytes.Replace(b, []byte("\\u003c"), []byte("<"), -1)
	b = bytes.Replace(b, []byte("\\u003e"), []byte(">"), -1)

	return b, nil
}
