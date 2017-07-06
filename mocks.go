package main

import (
	"github.com/Financial-Times/message-queue-go-producer/producer"
	"github.com/stretchr/testify/mock"
)

type mockProducer struct {
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

func (m *mockImageSetMapper) Map(source NativeContent, lastModified string, publishReference string) ([]JSONImageSet, error) {
	args := m.Called(source)
	return args.Get(0).([]JSONImageSet), args.Error(1)
}

type mockedArticleToImageSetMapper struct {
	mock.Mock
}

func (m *mockedArticleToImageSetMapper) Map(source []byte) ([]XMLImageSet, error) {
	args := m.Called(source)
	return args.Get(0).([]XMLImageSet), args.Error(1)
}

type mockedXmlImageSetToJSONMapper struct {
	mock.Mock
}

func (m *mockedXmlImageSetToJSONMapper) Map(xmlImageSets []XMLImageSet, articleUuid string, attributes xmlAttributes, lastModified string, publishReference string) ([]JSONImageSet, error) {
	args := m.Called(xmlImageSets)
	return args.Get(0).([]JSONImageSet), args.Error(1)
}

type mockAttributesMapper struct {
	mock.Mock
}

func (m *mockAttributesMapper) Map(source string) (xmlAttributes, error) {
	args := m.Called(source)
	return args.Get(0).(xmlAttributes), args.Error(1)
}
