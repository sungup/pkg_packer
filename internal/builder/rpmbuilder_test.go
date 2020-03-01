package builder

import (
	"bufio"
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/sungup/pkg_packer/test"
	"io/ioutil"
	"os"
	"testing"
)

func TestRPMBuilder_Filename(t *testing.T) {
	a := assert.New(t)

	expected, tested := LoadYAMLExpectAndPackage()

	builder := NewRPMBuilder(tested)

	var data string
	var err error

	// Normal test
	data, err = builder.Filename()
	if a.NoError(err) {
		a.Equal(expected.Test["rpmfile"], data)
	}

	// Name check test
	data = tested.Meta.Name
	tested.Meta.Name = ""
	// TODO add error type check
	_, err = builder.Filename()
	a.Error(err)
	tested.Meta.Name = data

	// Version check test
	data = tested.Meta.Version
	tested.Meta.Version = ""
	// TODO add error type check
	_, err = builder.Filename()
	a.Error(err)
	tested.Meta.Version = data
}

func TestRPMBuilder_Build(t *testing.T) {
	var err error

	a := assert.New(t)
	expected, tested := LoadYAMLExpectAndPackage()

	// 0. Sync the rpm file's build time
	tested.Meta.UpdateBuildTime(test.ExpectFileDate.UTC())

	builder := NewRPMBuilder(tested)

	// 1. load expected
	expectedFile, err := test.GetTestFilePath(expected.Test["rpmfile"])
	a.NoError(err)

	expectedData, err := ioutil.ReadFile(expectedFile)
	a.NoError(err)

	// 2. make byte buffer and build
	buffer := new(bytes.Buffer)

	err = builder.Build(buffer)
	a.NoError(err)

	// 3. check output
	a.Equal(expectedData, buffer.Bytes())

	////////
	filePath, _ := builder.Filename()
	filePath, _ = test.GetTestFilePath(filePath)

	rpmFile, _ := os.Create(filePath)
	rpmWriter := bufio.NewWriter(rpmFile)

	_ = builder.Build(rpmWriter)
	_ = rpmWriter.Flush()
}
