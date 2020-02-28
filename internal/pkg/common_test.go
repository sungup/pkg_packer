package pkg

import (
	"github.com/stretchr/testify/assert"
	"github.com/sungup/pkg_packer/test"
	"testing"
)

func LoadYAMLExpectAndPackage(t *testing.T) (*test.YAMLTestData, *Package) {
	expect, err := test.LoadTestYAMLData(test.ExpectFile)

	assert.NoErrorf(t, err, "reference set loading fail")

	expectPath, _ := test.GetTestFilePath(test.ExpectFile)
	tested, err := LoadPkgInfo(expectPath)

	assert.NoErrorf(t, err, "unexpected load error")

	return expect, tested
}
