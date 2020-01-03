package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNativeMap_Ok(t *testing.T) {
	m := defaultMessageToNativeMapper{}
	source := []byte(`{"type":"EOM::CompoundStory","value":"PGRvYz48L2RvYz4=","attributes":"\u003c?xml version=\"1.0\" encoding=\"UTF-8\"?\u003e\n"}`)
	native, err := m.Map(source)
	assert.NoError(t, err, "Error wasn't expected during mapping")
	assert.Equal(t, NativeContent{Type: compoundStory, Value: "PGRvYz48L2RvYz4=", Attributes: "\u003c?xml version=\"1.0\" encoding=\"UTF-8\"?\u003e\n"}, native)
}

func TestNativeMap_NOk(t *testing.T) {
	m := defaultMessageToNativeMapper{}
	source := []byte(`{{`)
	_, err := m.Map(source)
	assert.Error(t, err, "An error was expected")
}
