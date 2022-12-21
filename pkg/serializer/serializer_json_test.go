package serializer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJsonMarshal(t *testing.T) {
	assert := assert.New(t)
	type Test struct {
		Name string
	}
	test := Test{Name: "test"}
	data, err := Marshal(test)
	assert.Nil(err)
	assert.Equal(`{"Name":"test"}`, string(data))
}

func TestJsonUnmarshal(t *testing.T) {
	assert := assert.New(t)
	type Test struct {
		Name string
	}
	test := Test{}
	err := Unmarshal([]byte(`{"Name":"test"}`), &test)
	assert.Nil(err)
	assert.Equal("test", test.Name)
}

func TestNewDecoder(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(NewDecoder(nil))
}

func TestNewEncoder(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(NewEncoder(nil))
}
