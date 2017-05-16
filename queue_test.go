package main

import (
	"testing"
	"github.com/Financial-Times/message-queue-gonsumer/consumer"
	"github.com/Financial-Times/message-queue-go-producer/producer"
	"github.com/stretchr/testify/mock"
	"strings"
	"github.com/pkg/errors"
)

func TestOnMessage_Ok(t *testing.T) {
	sourceMsg := consumer.Message{
		Headers: map[string]string {
			"X-Request-Id": "tid_test123",
			"Origin-System-Id": methodeSystemOrigin,
			"Message-Timestamp": "2017-05-15T15:54:32.166Z",
		},
	}
	nativeContent := NativeContent{
		Type: compoundStory,
		Value: "",
	}
	mockedMessageToNativeMapper := new(mockMessageToNativeMapper)
	mockedMessageToNativeMapper.On("Map", mock.MatchedBy(func(source []byte) bool { return true })).Return(nativeContent, nil)
	mockedImageSetMapper := new(mockImageSetMapper)
	jsonImageSets := []JSONImageSet{JSONImageSet{UUID: "512c1f3d-e48c-4618-863c-94bc9d913b9b"}, JSONImageSet{ UUID: "43dc1ff3-6d6c-41f3-9196-56dcaa554905"}}
	mockedImageSetMapper.On("Map", mock.MatchedBy(func(source NativeContent) bool { return true })).Return(jsonImageSets, nil)
	mockedProducer := new(mockProducer)
	mockedProducer.On("SendMessage", "", mock.MatchedBy(func(msg producer.Message) bool { return true })).Return(nil)
	q := newQueue(nil, mockedProducer, mockedMessageToNativeMapper, mockedImageSetMapper)
	q.onMessage(sourceMsg)
	mockedProducer.AssertCalled(t, "SendMessage", "",
		mock.MatchedBy(func(msg producer.Message) bool {
			return strings.Contains(msg.Body, "512c1f3d-e48c-4618-863c-94bc9d913b9b") && strings.Contains(msg.Body, "2017-05-15T15:54:32.166Z")
		}))
	mockedProducer.AssertCalled(t, "SendMessage", "",
		mock.MatchedBy(func(msg producer.Message) bool {
			return strings.Contains(msg.Body, "43dc1ff3-6d6c-41f3-9196-56dcaa554905") && strings.Contains(msg.Body, "2017-05-15T15:54:32.166Z")
		}))
	mockedProducer.AssertCalled(t, "SendMessage", "", mock.MatchedBy(func(msg producer.Message) bool { return strings.Contains(msg.Body, "2017-05-15T15:54:32.166Z") }))
	mockedProducer.AssertNumberOfCalls(t, "SendMessage", 2)
}

func TestOnMessage_OkWhenNoTimestamp(t *testing.T) {
	sourceMsg := consumer.Message{
		Headers: map[string]string {
			"X-Request-Id": "tid_test123",
			"Origin-System-Id": methodeSystemOrigin,
		},
	}
	nativeContent := NativeContent{
		Type: compoundStory,
		Value: "",
	}
	mockedMessageToNativeMapper := new(mockMessageToNativeMapper)
	mockedMessageToNativeMapper.On("Map", mock.MatchedBy(func(source []byte) bool { return true })).Return(nativeContent, nil)
	mockedImageSetMapper := new(mockImageSetMapper)
	jsonImageSets := []JSONImageSet{JSONImageSet{UUID: "512c1f3d-e48c-4618-863c-94bc9d913b9b"}, JSONImageSet{ UUID: "43dc1ff3-6d6c-41f3-9196-56dcaa554905"}}
	mockedImageSetMapper.On("Map", mock.MatchedBy(func(source NativeContent) bool { return true })).Return(jsonImageSets, nil)
	mockedProducer := new(mockProducer)
	mockedProducer.On("SendMessage", "", mock.MatchedBy(func(msg producer.Message) bool { return true })).Return(nil)
	q := newQueue(nil, mockedProducer, mockedMessageToNativeMapper, mockedImageSetMapper)
	q.onMessage(sourceMsg)
	mockedProducer.AssertCalled(t, "SendMessage", "",
		mock.MatchedBy(func(msg producer.Message) bool {
			return strings.Contains(msg.Body, "512c1f3d-e48c-4618-863c-94bc9d913b9b") && strings.Contains(msg.Body, "lastModified")
		}))
	mockedProducer.AssertCalled(t, "SendMessage", "",
		mock.MatchedBy(func(msg producer.Message) bool {
			return strings.Contains(msg.Body, "43dc1ff3-6d6c-41f3-9196-56dcaa554905") && strings.Contains(msg.Body, "lastModified")
		}))
	mockedProducer.AssertNumberOfCalls(t, "SendMessage", 2)
}

func TestOnMessage_SkipsWhenNotOriginSystem(t *testing.T) {
	sourceMsg := consumer.Message{
		Headers: map[string]string {
			"X-Request-Id": "tid_test123",
			"Origin-System-Id": "some other origin system",
		},
	}
	mockedProducer := new(mockProducer)
	q := newQueue(nil, mockedProducer, nil, nil)
	q.onMessage(sourceMsg)
	mockedProducer.AssertNotCalled(t, "SendMessage", "", mock.MatchedBy(func(msg producer.Message) bool { return true }))
}

func TestOnMessage_SkipsWhenNotRequestId(t *testing.T) {
	sourceMsg := consumer.Message{ Headers: map[string]string {} }
	mockedProducer := new(mockProducer)
	q := newQueue(nil, mockedProducer, nil, nil)
	q.onMessage(sourceMsg)
	mockedProducer.AssertNotCalled(t, "SendMessage", "", mock.MatchedBy(func(msg producer.Message) bool { return true }))
}

func TestOnMessage_WarnIfErrorMappingNative(t *testing.T) {
	sourceMsg := consumer.Message{
		Headers: map[string]string {
			"X-Request-Id": "tid_test123",
			"Origin-System-Id": methodeSystemOrigin,
			"Message-Timestamp": "2017-05-15T15:54:32.166Z",
		},
	}
	mockedMessageToNativeMapper := new(mockMessageToNativeMapper)
	mockedMessageToNativeMapper.On("Map", mock.MatchedBy(func(source []byte) bool { return true })).Return(NativeContent{}, errors.New("error mapping to native"))
	mockedProducer := new(mockProducer)
	q := newQueue(nil, mockedProducer, mockedMessageToNativeMapper, nil)
	q.onMessage(sourceMsg)
	mockedProducer.AssertNotCalled(t, "SendMessage", "", mock.MatchedBy(func(msg producer.Message) bool { return true }))
}

func TestOnMessage_SkipOtherTypes(t *testing.T) {
	sourceMsg := consumer.Message{
		Headers: map[string]string {
			"X-Request-Id": "tid_test123",
			"Origin-System-Id": methodeSystemOrigin,
			"Message-Timestamp": "2017-05-15T15:54:32.166Z",
		},
	}
	nativeContent := NativeContent{
		Type: "other type",
		Value: "",
	}
	mockedMessageToNativeMapper := new(mockMessageToNativeMapper)
	mockedMessageToNativeMapper.On("Map", mock.MatchedBy(func(source []byte) bool { return true })).Return(nativeContent, nil)
	mockedProducer := new(mockProducer)
	mockedProducer.On("SendMessage", "", mock.MatchedBy(func(msg producer.Message) bool { return true })).Return(nil)
	q := newQueue(nil, mockedProducer, mockedMessageToNativeMapper, nil)
	q.onMessage(sourceMsg)
	mockedProducer.AssertNotCalled(t, "SendMessage", "", mock.MatchedBy(func(msg producer.Message) bool { return true }))
}

func TestOnMessage_WarnIfErrorInImageSetMapper(t *testing.T) {
	sourceMsg := consumer.Message{
		Headers: map[string]string {
			"X-Request-Id": "tid_test123",
			"Origin-System-Id": methodeSystemOrigin,
			"Message-Timestamp": "2017-05-15T15:54:32.166Z",
		},
	}
	nativeContent := NativeContent{
		Type: compoundStory,
		Value: "",
	}
	mockedMessageToNativeMapper := new(mockMessageToNativeMapper)
	mockedMessageToNativeMapper.On("Map", mock.MatchedBy(func(source []byte) bool { return true })).Return(nativeContent, nil)
	mockedImageSetMapper := new(mockImageSetMapper)
	mockedImageSetMapper.On("Map", mock.MatchedBy(func(source NativeContent) bool { return true })).Return([]JSONImageSet{}, errors.New("error mapping to image sets"))
	mockedProducer := new(mockProducer)
	mockedProducer.On("SendMessage", "", mock.MatchedBy(func(msg producer.Message) bool { return true })).Return(nil)
	q := newQueue(nil, mockedProducer, mockedMessageToNativeMapper, mockedImageSetMapper)
	q.onMessage(sourceMsg)
	mockedProducer.AssertNotCalled(t, "SendMessage", "", mock.MatchedBy(func(msg producer.Message) bool { return true }))
}

type mockProducer struct{
	mock.Mock
}

func (p *mockProducer) SendMessage(key string, msg producer.Message) error {
	args := p.Called(key, msg)
	return args.Error(0)
}

func (p *mockProducer) ConnectivityCheck() (string, error) {
	args := p.Called()
	return args.String(0), args.Error(1)
}

type mockMessageToNativeMapper struct {
	mock.Mock
}

func (m *mockMessageToNativeMapper) Map(source []byte) (NativeContent, error) {
	args := m.Called(source)
	return args.Get(0).(NativeContent), args.Error(1)
}

type mockImageSetMapper struct {
	mock.Mock
}

func (m *mockImageSetMapper) Map(source NativeContent) ([]JSONImageSet, error) {
	args := m.Called(source)
	return args.Get(0).([]JSONImageSet), args.Error(1)
}
