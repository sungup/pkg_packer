package builder

import (
	Assert "github.com/stretchr/testify/assert"
	"testing"
)

func TestRPMBuilder_Filename(t *testing.T) {
	assert := Assert.New(t)

	expected, tested := LoadYAMLExpectAndPackage()

	builder := RPMBuilder{}

	var data string
	var err error

	// Normal test
	data, err = builder.Filename(&tested.Meta)
	if assert.NoError(err) {
		assert.Equal(expected.Test["rpmfile"], data)
	}

	// Name check test
	data = tested.Meta.Name
	tested.Meta.Name = ""
	// TODO add error type check
	_, err = builder.Filename(&tested.Meta)
	assert.Error(err)
	tested.Meta.Name = data

	// Version check test
	data = tested.Meta.Version
	tested.Meta.Version = ""
	// TODO add error type check
	_, err = builder.Filename(&tested.Meta)
	assert.Error(err)
	tested.Meta.Version = data

}
