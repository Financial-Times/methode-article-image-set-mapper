package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestNativeMap_Ok(t *testing.T) {
	m := defaultMessageToNativeMapper{}
	source := []byte(`{"type":"EOM::CompoundStory","value":"PGRvYz48L2RvYz4="}`)
	native, err := m.Map(source)
	assert.NoError(t, err, "Error wasn't expected during mapping")
	assert.Equal(t, compoundStory, native.Type)
	assert.Equal(t, "PGRvYz48L2RvYz4=", native.Value)
}

func TestNativeMap_NOk(t *testing.T) {
	m := defaultMessageToNativeMapper{}
	source := []byte(`{{`)
	_, err := m.Map(source)
	assert.Error(t, err, "An error was expected")
}
