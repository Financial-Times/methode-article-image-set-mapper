package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAISMap_Ok(t *testing.T) {
	m := defaultArticleToImageSetMapper{}
	source := []byte(`
	<doc>
		<story>
			<text>
				<body>
					Somebody
					<image-set id="U22104508221701xCD">
					       <image-small id="U11603507121721yBF" fileref="/FT/Graphics/Online/Z_Undefined/2017/03/timeline-artboards-s.png?uuid=2ae43059-c725-4e6f-95d7-45f04f2e33b6" xtransform=" scale(0.0946667 0.0946667)" tmx="750 890 71 84" />
					       <image-medium id="U11603507121721TXC" fileref="/FT/Graphics/Online/Z_Undefined/2017/03/timeline-artboards-m.png?uuid=78ed71df-457f-41a9-95a2-ef69622ccf13" xtransform=" scale(0.1013514 0.1013514)" tmx="1480 960 150 97" />
					       <image-large id="U11603507121721ey" fileref="/FT/Graphics/Online/Z_Undefined/2017/03/timeline-artboards-l.png?uuid=4a29a412-d94b-46af-a36f-e7be0dfe20f6" xtransform=" scale(0.0775194 0.0775194)" tmx="2580 1124 200 87" />
					</image-set>
					is
					<image-set id="U33104508221999xAA">
					       <image-small id="U11603507221721yBF" fileref="/FT/Graphics/Online/Z_Undefined/2017/03/timeline-artboards-s.png?uuid=83a927a3-69ff-407d-9ae6-ba9d06fbdc89" xtransform=" scale(0.0946667 0.0946667)" tmx="750 890 71 84" />
					       <image-medium id="U11603508121721TXC" fileref="/FT/Graphics/Online/Z_Undefined/2017/03/timeline-artboards-m.png?uuid=0e4116ae-22bb-4eac-8380-26955d5ffc04" xtransform=" scale(0.1013514 0.1013514)" tmx="1480 960 150 97" />
					       <image-large id="U11603507221721ey" fileref="/FT/Graphics/Online/Z_Undefined/2017/03/timeline-artboards-l.png?uuid=0912908c-9f0b-4cc1-be0d-3cce248f4183" xtransform=" scale(0.0775194 0.0775194)" tmx="2580 1124 200 87" />
					</image-set>
					reading.
				</body>
			</text>
		</story>
	</doc>
	`)
	actualImageSets, err := m.Map(source)
	assert.NoError(t, err, "Error wasn't expected during mapping")
	expectedXMLImageSets := []XMLImageSet{
		XMLImageSet {
			ID:          "U22104508221701xCD",
			ImageSmall:  XMLImage{FileRef: "/FT/Graphics/Online/Z_Undefined/2017/03/timeline-artboards-s.png?uuid=2ae43059-c725-4e6f-95d7-45f04f2e33b6"},
			ImageMedium: XMLImage{FileRef: "/FT/Graphics/Online/Z_Undefined/2017/03/timeline-artboards-m.png?uuid=78ed71df-457f-41a9-95a2-ef69622ccf13"},
			ImageLarge:  XMLImage{FileRef: "/FT/Graphics/Online/Z_Undefined/2017/03/timeline-artboards-l.png?uuid=4a29a412-d94b-46af-a36f-e7be0dfe20f6"},
		},
		XMLImageSet {
			ID:          "U33104508221999xAA",
			ImageSmall:  XMLImage{FileRef: "/FT/Graphics/Online/Z_Undefined/2017/03/timeline-artboards-s.png?uuid=83a927a3-69ff-407d-9ae6-ba9d06fbdc89"},
			ImageMedium: XMLImage{FileRef: "/FT/Graphics/Online/Z_Undefined/2017/03/timeline-artboards-m.png?uuid=0e4116ae-22bb-4eac-8380-26955d5ffc04"},
			ImageLarge:  XMLImage{FileRef: "/FT/Graphics/Online/Z_Undefined/2017/03/timeline-artboards-l.png?uuid=0912908c-9f0b-4cc1-be0d-3cce248f4183"},
		},
	}
	assert.Equal(t, expectedXMLImageSets, actualImageSets)
}

func TestAISMap_NOk(t *testing.T) {
	m := defaultArticleToImageSetMapper{}
	_, err := m.Map([]byte("<doc><stor**ERROR**</story></doc>"))
	assert.Error(t, err, "Error was expected during unmarshal")
}
