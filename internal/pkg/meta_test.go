package pkg

import (
	Assert "github.com/stretchr/testify/assert"
	"testing"
)

func CompareMetaData(t *testing.T, expect map[string]string, tested *PackageMeta) {
	assert := Assert.New(t)

	assert.Equal(expect["name"], tested.Name)
	assert.Equal(expect["version"], tested.Version)
	assert.Equal(expect["release"], tested.Release)
	assert.Equal(expect["arch"], tested.Arch)
	assert.Equal(expect["summary"], tested.Summary)
	assert.Equal(expect["desc"], tested.Description)
}
