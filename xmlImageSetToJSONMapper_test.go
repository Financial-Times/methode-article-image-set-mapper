package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestXMLJSONMap_Ok(t *testing.T) {
	m := defaultImageSetToJSONMapper{}
	uuid := "c17e8abe-1df8-11e7-942c-4a4c42b3072e"
	source := []XMLImageSet{
		XMLImageSet{
			ID: "U11603547146784PeC",
			ImageSmall: XMLImage{
				FileRef: "/FT/Graphics/Online/Z_Undefined/2017/03/timeline-artboards-s.png?uuid=4258f26a-13c5-11e7-9469-afea892e4de3",
			},
			ImageMedium: XMLImage{
				FileRef: "/FT/Graphics/Online/Z_Undefined/2017/03/timeline-artboards-m.png?uuid=41614f4c-13c5-11e7-9469-afea892e4de3",
			},
			ImageLarge: XMLImage{
				FileRef: "/FT/Graphics/Online/Z_Undefined/2017/03/timeline-artboards-l.png?uuid=3ff3b7a8-13c5-11e7-9469-afea892e4de3",
			},
		},
		XMLImageSet{
			ID: "U12345547146784RfD",
			ImageSmall: XMLImage{
				FileRef: "/FT/Graphics/Online/Z_Undefined/2017/03/other-s.png?uuid=404cf8d9-1b88-4883-8afe-580e5174830d",
			},
			ImageMedium: XMLImage{
				FileRef: "/FT/Graphics/Online/Z_Undefined/2017/03/other-m.png?uuid=2fe0b459-a23e-452d-a2aa-2e0503982ed2",
			},
			ImageLarge: XMLImage{
				FileRef: "/FT/Graphics/Online/Z_Undefined/2017/03/other-l.png?uuid=2acf1caa-8014-48ec-b070-a0ffbc45d1d5",
			},
		},
	}
	attributes := xmlAttributes{
		OutputChannels: OutputChannels{
			DIFTcom{
				DIFTcomLastPublication:    "20170518132425",
				DIFTcomInitialPublication: "20170518132400",
			},
		},
	}
	actualImageSets, err := m.Map(source, uuid, attributes, "2017-05-17T13:46:01.100Z", "tid_test")
	if err != nil {
		assert.Error(t, err, "error mapping set")
	}
	expectedImageSets := []JSONImageSet{
		JSONImageSet{
			UUID: "d8367364-c56b-3599-8787-08e1784b02ce",
			Identifiers: []JSONIdentifier{
				JSONIdentifier{
					Authority:       "http://api.ft.com/system/FTCOM-METHODE",
					IdentifierValue: "d8367364-c56b-3599-8787-08e1784b02ce",
				},
			},
			Members: []JSONMember{
				JSONMember{
					UUID: "41614f4c-13c5-11e7-9469-afea892e4de3",
				},
				JSONMember{
					UUID:            "4258f26a-13c5-11e7-9469-afea892e4de3",
					MaxDisplayWidth: "490px",
				},
				JSONMember{
					UUID:            "3ff3b7a8-13c5-11e7-9469-afea892e4de3",
					MinDisplayWidth: "980px",
				},
			},
			PublishReference:   "tid_test",
			LastModified:       "2017-05-17T13:46:01.100Z",
			PublishedDate:      "2017-05-18T13:24:25.000Z",
			FirstPublishedDate: "2017-05-18T13:24:00.000Z",
			CanBeDistributed:   "yes",
			Type:               "ImageSet",
		},
		JSONImageSet{
			UUID: "84be18d3-4622-3bb1-87b6-33786f12902f",
			Identifiers: []JSONIdentifier{
				JSONIdentifier{
					Authority:       "http://api.ft.com/system/FTCOM-METHODE",
					IdentifierValue: "84be18d3-4622-3bb1-87b6-33786f12902f",
				},
			},
			Members: []JSONMember{
				JSONMember{
					UUID: "2fe0b459-a23e-452d-a2aa-2e0503982ed2",
				},
				JSONMember{
					UUID:            "404cf8d9-1b88-4883-8afe-580e5174830d",
					MaxDisplayWidth: "490px",
				},
				JSONMember{
					UUID:            "2acf1caa-8014-48ec-b070-a0ffbc45d1d5",
					MinDisplayWidth: "980px",
				},
			},
			PublishReference:   "tid_test",
			LastModified:       "2017-05-17T13:46:01.100Z",
			PublishedDate:      "2017-05-18T13:24:25.000Z",
			FirstPublishedDate: "2017-05-18T13:24:00.000Z",
			CanBeDistributed:   "yes",
			Type:               "ImageSet",
		},
	}
	assert.Equal(t, expectedImageSets, actualImageSets)
}

func TestXMLJSONMap_LessThan3(t *testing.T) {
	m := defaultImageSetToJSONMapper{}
	uuid := "c17e8abe-1df8-11e7-942c-4a4c42b3072e"
	source := []XMLImageSet{
		XMLImageSet{
			ID: "U11603547146784PeC",
			ImageSmall: XMLImage{
				FileRef: "/FT/Graphics/Online/Z_Undefined/2017/03/timeline-artboards-s.png?uuid=4258f26a-13c5-11e7-9469-afea892e4de3",
			},
		},
	}
	attributes := xmlAttributes{
		OutputChannels: OutputChannels{
			DIFTcom{
				DIFTcomLastPublication:    "20170518022425",
				DIFTcomInitialPublication: "20170518022400",
			},
		},
	}
	actualImageSets, err := m.Map(source, uuid, attributes, "2017-05-17T13:46:01.100Z", "tid_test")
	if err != nil {
		assert.Error(t, err, "error mapping set")
	}
	expectedImageSets := []JSONImageSet{
		JSONImageSet{
			UUID: "d8367364-c56b-3599-8787-08e1784b02ce",
			Identifiers: []JSONIdentifier{
				JSONIdentifier{
					Authority:       "http://api.ft.com/system/FTCOM-METHODE",
					IdentifierValue: "d8367364-c56b-3599-8787-08e1784b02ce",
				},
			},
			Members: []JSONMember{
				JSONMember{
					UUID:            "4258f26a-13c5-11e7-9469-afea892e4de3",
					MaxDisplayWidth: "490px",
				},
			},
			PublishReference:   "tid_test",
			LastModified:       "2017-05-17T13:46:01.100Z",
			PublishedDate:      "2017-05-18T02:24:25.000Z",
			FirstPublishedDate: "2017-05-18T02:24:00.000Z",
			CanBeDistributed:   "yes",
			Type:               "ImageSet",
		},
	}
	assert.Equal(t, expectedImageSets, actualImageSets)
}

func TestAppendIfPresent_Present(t *testing.T) {
	mapper := defaultImageSetToJSONMapper{}
	members := make([]JSONMember, 0)
	mapper.appendIfPresent(&members, XMLImage{FileRef: "/FT/Graphics/Online/Z_Undefined/2017/03/timeline-artboards-s.png?uuid=4258f26a-13c5-11e7-9469-afea892e4de3"}, "any", "", "980px")
	assert.Contains(t, members, JSONMember{UUID: "4258f26a-13c5-11e7-9469-afea892e4de3", MinDisplayWidth: "980px"})
}

func TestAppendIfPresent_NoUuid(t *testing.T) {
	mapper := defaultImageSetToJSONMapper{}
	members := make([]JSONMember, 0)
	mapper.appendIfPresent(&members, XMLImage{FileRef: "/FT/Graphics/Online/Z_Undefined/2017/03/timeline-artboards-s.png"}, "any", "", "980px")
	assert.Equal(t, len(members), 0)
}

func TestAppendIfPresent_NoEntry(t *testing.T) {
	mapper := defaultImageSetToJSONMapper{}
	members := make([]JSONMember, 0)
	mapper.appendIfPresent(&members, XMLImage{}, "any", "", "980px")
	assert.Equal(t, len(members), 0)
}
