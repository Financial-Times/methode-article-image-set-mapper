package main

import (
	"errors"
	"strings"
	"testing"

	"github.com/Financial-Times/message-queue-go-producer/producer"
	consumer "github.com/Financial-Times/message-queue-gonsumer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestOnMessage_Ok(t *testing.T) {
	sourceMsg := consumer.Message{
		Headers: map[string]string{
			"X-Request-Id":      "tid_test123",
			"Origin-System-Id":  methodeSystemOrigin,
			"Message-Timestamp": "2017-05-15T15:54:32.166Z",
		},
	}
	nativeContent := NativeContent{
		Type:  compoundStory,
		Value: "",
	}
	mockedMessageToNativeMapper := new(mockMessageToNativeMapper)
	mockedMessageToNativeMapper.On("Map", mock.MatchedBy(func(source []byte) bool { return true })).Return(nativeContent, nil)
	mockedImageSetMapper := new(mockImageSetMapper)
	jsonImageSets := []JSONImageSet{JSONImageSet{UUID: "512c1f3d-e48c-4618-863c-94bc9d913b9b"}, JSONImageSet{UUID: "43dc1ff3-6d6c-41f3-9196-56dcaa554905"}}
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
		Headers: map[string]string{
			"X-Request-Id":     "tid_test123",
			"Origin-System-Id": methodeSystemOrigin,
		},
	}
	nativeContent := NativeContent{
		Type:  compoundStory,
		Value: "",
	}
	mockedMessageToNativeMapper := new(mockMessageToNativeMapper)
	mockedMessageToNativeMapper.On("Map", mock.MatchedBy(func(source []byte) bool { return true })).Return(nativeContent, nil)
	mockedImageSetMapper := new(mockImageSetMapper)
	jsonImageSets := []JSONImageSet{JSONImageSet{UUID: "512c1f3d-e48c-4618-863c-94bc9d913b9b"}, JSONImageSet{UUID: "43dc1ff3-6d6c-41f3-9196-56dcaa554905"}}
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
		Headers: map[string]string{
			"X-Request-Id":     "tid_test123",
			"Origin-System-Id": "some other origin system",
		},
	}
	mockedProducer := new(mockProducer)
	q := newQueue(nil, mockedProducer, nil, nil)
	q.onMessage(sourceMsg)
	mockedProducer.AssertNotCalled(t, "SendMessage", "", mock.MatchedBy(func(msg producer.Message) bool { return true }))
}

func TestOnMessage_SkipsWhenNotRequestId(t *testing.T) {
	sourceMsg := consumer.Message{Headers: map[string]string{}}
	mockedProducer := new(mockProducer)
	q := newQueue(nil, mockedProducer, nil, nil)
	q.onMessage(sourceMsg)
	mockedProducer.AssertNotCalled(t, "SendMessage", "", mock.MatchedBy(func(msg producer.Message) bool { return true }))
}

func TestOnMessage_WarnIfErrorMappingNative(t *testing.T) {
	sourceMsg := consumer.Message{
		Headers: map[string]string{
			"X-Request-Id":      "tid_test123",
			"Origin-System-Id":  methodeSystemOrigin,
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
		Headers: map[string]string{
			"X-Request-Id":      "tid_test123",
			"Origin-System-Id":  methodeSystemOrigin,
			"Message-Timestamp": "2017-05-15T15:54:32.166Z",
		},
	}
	nativeContent := NativeContent{
		Type:  "other type",
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
		Headers: map[string]string{
			"X-Request-Id":      "tid_test123",
			"Origin-System-Id":  methodeSystemOrigin,
			"Message-Timestamp": "2017-05-15T15:54:32.166Z",
		},
	}
	nativeContent := NativeContent{
		Type:  compoundStory,
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

func TestOnMessage_OneSendFailureShouldNotAffectOther(t *testing.T) {
	sourceMsg := consumer.Message{
		Headers: map[string]string{
			"X-Request-Id":      "tid_test123",
			"Origin-System-Id":  methodeSystemOrigin,
			"Message-Timestamp": "2017-05-15T15:54:32.166Z",
		},
	}
	nativeContent := NativeContent{
		Type:  compoundStory,
		Value: "",
	}
	mockedMessageToNativeMapper := new(mockMessageToNativeMapper)
	mockedMessageToNativeMapper.On("Map", mock.MatchedBy(func(source []byte) bool { return true })).Return(nativeContent, nil)
	mockedImageSetMapper := new(mockImageSetMapper)
	jsonImageSets := []JSONImageSet{JSONImageSet{UUID: "512c1f3d-e48c-4618-863c-94bc9d913b9b"}, JSONImageSet{UUID: "43dc1ff3-6d6c-41f3-9196-56dcaa554905"}}
	mockedImageSetMapper.On("Map", mock.MatchedBy(func(source NativeContent) bool { return true })).Return(jsonImageSets, nil)
	mockedProducer := new(mockProducer)
	mockedProducer.On("SendMessage", "", mock.MatchedBy(func(msg producer.Message) bool {
		return strings.Contains(msg.Body, "512c1f3d-e48c-4618-863c-94bc9d913b9b")
	})).Return(errors.New("error sending first msg"))
	mockedProducer.On("SendMessage", "", mock.MatchedBy(func(msg producer.Message) bool {
		return strings.Contains(msg.Body, "43dc1ff3-6d6c-41f3-9196-56dcaa554905")
	})).Return(nil)
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

func TestBuildMessage_Ok(t *testing.T) {
	q := newQueue(nil, nil, nil, nil)
	actualMsg, err := q.buildMessage(JSONImageSet{
		UUID: "5a8f3f37-3098-48f7-811a-f69d12f2b1be",
		Members: []JSONMember{
			JSONMember{UUID: "8ff1c7f4-a80b-4b8d-8821-b07ff1bfdf87"},
			JSONMember{
				UUID:            "3bea853a-89b8-4831-80b3-8384e962f5dc",
				MaxDisplayWidth: "490px",
			},
			JSONMember{
				UUID:            "c6eeea75-748e-4b1c-a046-6e4c9d81ff25",
				MinDisplayWidth: "980px",
			},
		},
		Identifiers: []JSONIdentifier{
			JSONIdentifier{
				Authority:       methodeAuthority,
				IdentifierValue: "5a8f3f37-3098-48f7-811a-f69d12f2b1be",
			},
		},
		PublishedDate:      "2017-05-18T02:24:25Z",
		FirstPublishedDate: "2017-05-18T02:24:00Z",
		CanBeDistributed:   "yes",
		Type:               "ImageSet",
	}, "2017-05-15T15:54:32.166Z", "tid_test")
	assert.NoError(t, err, "Error wasn't expected during buildMessage()")
	assert.Equal(t, actualMsg.Headers["X-Request-Id"], "tid_test")
	assert.NotEmpty(t, actualMsg.Headers["Message-Id"])
	assert.Equal(t, actualMsg.Headers["Message-Type"], "cms-content-published")
	assert.Equal(t, actualMsg.Headers["Content-Type"], "application/json")
	assert.Equal(t, actualMsg.Headers["Origin-System-Id"], methodeSystemOrigin)
	assert.Equal(t, actualMsg.Body, `{"contentUri":"http://methode-article-image-set-mapper.svc.ft.com/image-set/model/5a8f3f37-3098-48f7-811a-f69d12f2b1be","payload":{"uuid":"5a8f3f37-3098-48f7-811a-f69d12f2b1be","identifiers":[{"authority":"http://api.ft.com/system/FTCOM-METHODE","identifierValue":"5a8f3f37-3098-48f7-811a-f69d12f2b1be"}],"members":[{"uuid":"8ff1c7f4-a80b-4b8d-8821-b07ff1bfdf87"},{"uuid":"3bea853a-89b8-4831-80b3-8384e962f5dc","maxDisplayWidth":"490px"},{"uuid":"c6eeea75-748e-4b1c-a046-6e4c9d81ff25","minDisplayWidth":"980px"}],"publishReference":"","lastModified":"","publishedDate":"2017-05-18T02:24:25Z","firstPublishedDate":"2017-05-18T02:24:00Z","canBeDistributed":"yes","type":"ImageSet"},"lastModified":"2017-05-15T15:54:32.166Z"}`)
}

func TestBuildMessages_Ok(t *testing.T) {
	q := newQueue(nil, nil, nil, nil)
	actualMsgs, errs := q.buildMessages([]JSONImageSet{
		JSONImageSet{
			UUID: "5a8f3f37-3098-48f7-811a-f69d12f2b1be",
			Members: []JSONMember{
				JSONMember{UUID: "8ff1c7f4-a80b-4b8d-8821-b07ff1bfdf87"},
				JSONMember{UUID: "3bea853a-89b8-4831-80b3-8384e962f5dc"},
				JSONMember{UUID: "c6eeea75-748e-4b1c-a046-6e4c9d81ff25"},
			},
			Identifiers: []JSONIdentifier{
				JSONIdentifier{
					Authority:       methodeAuthority,
					IdentifierValue: "5a8f3f37-3098-48f7-811a-f69d12f2b1be",
				},
			},
			PublishedDate:      "2017-05-18T02:24:25Z",
			FirstPublishedDate: "2017-05-18T02:24:00Z",
			CanBeDistributed:   "yes",
			Type:               "ImageSet",
		},
		JSONImageSet{
			UUID: "270c0151-7742-4c1e-b77e-a5557881a042",
			Members: []JSONMember{
				JSONMember{UUID: "667ee7f3-4f58-4080-a6f9-9b16b633dea8"},
				JSONMember{UUID: "a0513a50-08d1-43f6-af2b-7e7dc4d40b31"},
				JSONMember{UUID: "47e5a693-cd39-4ede-a016-244e6413a7fa"},
			},
			Identifiers: []JSONIdentifier{
				JSONIdentifier{
					Authority:       methodeAuthority,
					IdentifierValue: "270c0151-7742-4c1e-b77e-a5557881a042",
				},
			},
			PublishedDate:      "2017-05-18T02:24:25Z",
			FirstPublishedDate: "2017-05-18T02:24:00Z",
			CanBeDistributed:   "yes",
			Type:               "ImageSet",
		},
	}, "2017-05-15T15:54:32.166Z", "tid_test")
	if len(errs) != 0 {
		assert.Fail(t, "errors are not empty")
	}
	assert.Equal(t, actualMsgs["5a8f3f37-3098-48f7-811a-f69d12f2b1be"].Headers["X-Request-Id"], "tid_test")
	assert.NotEmpty(t, actualMsgs["5a8f3f37-3098-48f7-811a-f69d12f2b1be"].Headers["Message-Id"])
	assert.Equal(t, actualMsgs["5a8f3f37-3098-48f7-811a-f69d12f2b1be"].Headers["Message-Type"], "cms-content-published")
	assert.Equal(t, actualMsgs["5a8f3f37-3098-48f7-811a-f69d12f2b1be"].Headers["Content-Type"], "application/json")
	assert.Equal(t, actualMsgs["5a8f3f37-3098-48f7-811a-f69d12f2b1be"].Headers["Origin-System-Id"], methodeSystemOrigin)
	assert.Equal(t, actualMsgs["5a8f3f37-3098-48f7-811a-f69d12f2b1be"].Body, `{"contentUri":"http://methode-article-image-set-mapper.svc.ft.com/image-set/model/5a8f3f37-3098-48f7-811a-f69d12f2b1be","payload":{"uuid":"5a8f3f37-3098-48f7-811a-f69d12f2b1be","identifiers":[{"authority":"http://api.ft.com/system/FTCOM-METHODE","identifierValue":"5a8f3f37-3098-48f7-811a-f69d12f2b1be"}],"members":[{"uuid":"8ff1c7f4-a80b-4b8d-8821-b07ff1bfdf87"},{"uuid":"3bea853a-89b8-4831-80b3-8384e962f5dc"},{"uuid":"c6eeea75-748e-4b1c-a046-6e4c9d81ff25"}],"publishReference":"","lastModified":"","publishedDate":"2017-05-18T02:24:25Z","firstPublishedDate":"2017-05-18T02:24:00Z","canBeDistributed":"yes","type":"ImageSet"},"lastModified":"2017-05-15T15:54:32.166Z"}`)

	assert.Equal(t, actualMsgs["270c0151-7742-4c1e-b77e-a5557881a042"].Headers["X-Request-Id"], "tid_test")
	assert.NotEmpty(t, actualMsgs["270c0151-7742-4c1e-b77e-a5557881a042"].Headers["Message-Id"])
	assert.Equal(t, actualMsgs["270c0151-7742-4c1e-b77e-a5557881a042"].Headers["Message-Type"], "cms-content-published")
	assert.Equal(t, actualMsgs["270c0151-7742-4c1e-b77e-a5557881a042"].Headers["Content-Type"], "application/json")
	assert.Equal(t, actualMsgs["270c0151-7742-4c1e-b77e-a5557881a042"].Headers["Origin-System-Id"], methodeSystemOrigin)
	assert.Equal(t, actualMsgs["270c0151-7742-4c1e-b77e-a5557881a042"].Body, `{"contentUri":"http://methode-article-image-set-mapper.svc.ft.com/image-set/model/270c0151-7742-4c1e-b77e-a5557881a042","payload":{"uuid":"270c0151-7742-4c1e-b77e-a5557881a042","identifiers":[{"authority":"http://api.ft.com/system/FTCOM-METHODE","identifierValue":"270c0151-7742-4c1e-b77e-a5557881a042"}],"members":[{"uuid":"667ee7f3-4f58-4080-a6f9-9b16b633dea8"},{"uuid":"a0513a50-08d1-43f6-af2b-7e7dc4d40b31"},{"uuid":"47e5a693-cd39-4ede-a016-244e6413a7fa"}],"publishReference":"","lastModified":"","publishedDate":"2017-05-18T02:24:25Z","firstPublishedDate":"2017-05-18T02:24:00Z","canBeDistributed":"yes","type":"ImageSet"},"lastModified":"2017-05-15T15:54:32.166Z"}`)

}
