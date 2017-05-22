package main

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestISMap_Ok(t *testing.T) {
	mockedArticleToImageSetMapper := new(mockedArticleToImageSetMapper)
	mockedArticleToImageSetMapper.On("Map", mock.MatchedBy(func(source []byte) bool { return true })).Return([]XMLImageSet{}, nil)
	mockedAttributesMapper := new(mockAttributesMapper)
	mockedAttributesMapper.On("Map", mock.MatchedBy(func(source string) bool { return true })).Return(xmlAttributes{}, nil)
	mockedXmlImageSetToJSONMapper := new(mockedXmlImageSetToJSONMapper)
	expectedJsonImageSets := []JSONImageSet{
		JSONImageSet{
			UUID: "f02fbe32-9e2f-43fb-adbe-388d75ca23a9",
			Members: []JSONMember{
				JSONMember{
					UUID: "78ed71df-457f-41a9-95a2-ef69622ccf13",
				},
				JSONMember{
					UUID: "2ae43059-c725-4e6f-95d7-45f04f2e33b6",
				},
				JSONMember{
					UUID: "4a29a412-d94b-46af-a36f-e7be0dfe20f6",
				},
			},
			LastModified:     "2017-05-17T13:46:01.100Z",
			PublishReference: "tid_test",
		},
		JSONImageSet{
			UUID: "1ff5b8b1-13b3-4937-92a1-431e92d9b94d",
			Members: []JSONMember{
				JSONMember{
					UUID: "0e4116ae-22bb-4eac-8380-26955d5ffc04",
				},
				JSONMember{
					UUID: "83a927a3-69ff-407d-9ae6-ba9d06fbdc89",
				},
				JSONMember{
					UUID: "0912908c-9f0b-4cc1-be0d-3cce248f4183",
				},
			},
			LastModified:     "2017-05-17T13:46:01.100Z",
			PublishReference: "tid_test",
		},
	}
	mockedXmlImageSetToJSONMapper.On("Map", mock.MatchedBy(func(source []XMLImageSet) bool { return true })).Return(expectedJsonImageSets, nil)
	m := newImageSetMapper(mockedArticleToImageSetMapper, mockedAttributesMapper, mockedXmlImageSetToJSONMapper)
	source := NativeContent{Type: compoundStory, Value: "PGRvYz48L2RvYz4="}
	actualImageSets, err := m.Map(source, "2017-05-17T13:46:01.100Z", "tid_test")
	assert.NoError(t, err, "Error wasn't expected during mapping")
	assert.Equal(t, expectedJsonImageSets, actualImageSets)
}

func TestISMap_ErrorBase64(t *testing.T) {
	m := newImageSetMapper(nil, nil, nil)
	source := NativeContent{Type: compoundStory, Value: "***"}
	_, err := m.Map(source, "2017-05-17T13:46:01.100Z", "tid_test")
	assert.Error(t, err, "Error was expected during base64 decoding")
}

func TestISMap_ErrorXmlMapping(t *testing.T) {
	mockedArticleToImageSetMapper := new(mockedArticleToImageSetMapper)
	mockedArticleToImageSetMapper.On("Map", mock.MatchedBy(func(source []byte) bool { return true })).Return([]XMLImageSet{}, errors.New("error mapping article to xml imageSets"))
	m := newImageSetMapper(mockedArticleToImageSetMapper, nil, nil)
	source := NativeContent{Type: compoundStory, Value: "PGRvYz48L2RvYz4="}
	_, err := m.Map(source, "2017-05-17T13:46:01.100Z", "tid_test")
	assert.Error(t, err, "Error was expected during mapping article to xml imageSets")
}

func TestISMap_ErrorAttributesMapping(t *testing.T) {
	mockedArticleToImageSetMapper := new(mockedArticleToImageSetMapper)
	mockedArticleToImageSetMapper.On("Map", mock.MatchedBy(func(source []byte) bool { return true })).Return([]XMLImageSet{}, nil)
	mockedXmlImageSetToJSONMapper := new(mockedXmlImageSetToJSONMapper)
	mockedXmlImageSetToJSONMapper.On("Map", mock.MatchedBy(func(source []XMLImageSet) bool { return true })).Return([]JSONImageSet{}, nil)
	mockedAttributesMapper := new(mockAttributesMapper)
	mockedAttributesMapper.On("Map", mock.MatchedBy(func(source string) bool { return true })).Return(xmlAttributes{}, errors.New("error mapping attributes"))
	m := newImageSetMapper(mockedArticleToImageSetMapper, mockedAttributesMapper, mockedXmlImageSetToJSONMapper)
	source := NativeContent{Type: compoundStory, Value: "PGRvYz48L2RvYz4="}
	_, err := m.Map(source, "2017-05-17T13:46:01.100Z", "tid_test")
	assert.Error(t, err, "Error was expected during mapping xml attributes")
}

func TestISMap_ErrorJsonMapping(t *testing.T) {
	mockedArticleToImageSetMapper := new(mockedArticleToImageSetMapper)
	mockedArticleToImageSetMapper.On("Map", mock.MatchedBy(func(source []byte) bool { return true })).Return([]XMLImageSet{}, nil)
	mockedXmlImageSetToJSONMapper := new(mockedXmlImageSetToJSONMapper)
	mockedXmlImageSetToJSONMapper.On("Map", mock.MatchedBy(func(source []XMLImageSet) bool { return true })).Return([]JSONImageSet{}, errors.New("error mapping xml image set to json model"))
	mockedAttributesMapper := new(mockAttributesMapper)
	mockedAttributesMapper.On("Map", mock.MatchedBy(func(source string) bool { return true })).Return(xmlAttributes{}, nil)
	m := newImageSetMapper(mockedArticleToImageSetMapper, mockedAttributesMapper, mockedXmlImageSetToJSONMapper)
	source := NativeContent{Type: compoundStory, Value: "PGRvYz48L2RvYz4="}
	_, err := m.Map(source, "2017-05-17T13:46:01.100Z", "tid_test")
	assert.Error(t, err, "Error was expected during mapping xml image set to json model")
}
