package info

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/sungup/pkg_packer/test"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"
)

func TestFile_Body(t *testing.T) {
	// ignore testcase
}

func TestFile_SetBody(t *testing.T) {
	// ignore testcase
}

func TestFile_load(t *testing.T) {
	a := assert.New(t)

	now := time.Now()
	tested := File{
		MTime: now,
	}

	// temporary change ownership to raise permission denied
	notReadable, _ := test.GetTestFilePath("pkg.info/not-readable")
	_ = os.Chmod(notReadable, 0000)
	defer func() { _ = os.Chmod(notReadable, 0644) }()

	failTc := []string{"not-exist", "dir-file", "not-readable"}
	successTc := []string{"example.sh"}

	for _, tc := range failTc {
		a.Error(tested.load(path.Join("pkg.info", tc)))
		a.Equal(now, tested.MTime)
	}

	for _, tc := range successTc {
		tcFile := path.Join("pkg.info", tc)
		expectedFile, _ := test.GetTestFilePath(tcFile)
		expectedBody, _ := test.LoadTestFile(tcFile)
		expectedStat, _ := os.Stat(expectedFile)

		a.NoError(tested.load(tcFile))
		a.Equal(expectedBody, tested.body)
		a.Equal(expectedStat.ModTime(), tested.MTime)
		a.Equal(uint(expectedStat.Mode()), tested.Mode)
	}
}

func TestFile_UnmarshalYAML(t *testing.T) {
	a := assert.New(t)

	now := time.Now()
	tested := File{
		MTime: now,
	}

	genYml := func(dest, src, body string, mode uint) []byte {
		srcYml, bodyYml := "", ""
		if src != "" {
			srcYml = fmt.Sprintf("source: %s\n", src)
		}
		if body != "" {
			bodyYml = fmt.Sprintf("body: %s\n", body)
		}

		return []byte(fmt.Sprintf(
			"%s%sdest: %s\nmode: %d\nowner: test\ngroup: test",
			srcYml, bodyYml, dest, mode,
		))
	}

	expectedDest := "/usr/bin/helloworld"
	expectedSrc := "pkg.info/example.sh"
	invalidSrc := "pkg.info/dir-file"
	testBody := []byte("Text Hello World!")
	testMode := uint(0000)

	// 0. check unexpected yaml format
	a.Error(yaml.Unmarshal([]byte("hello world"), &tested))
	a.Equal(now, tested.MTime)

	// 1. check empty destination
	yml := genYml("", "", "", testMode)
	a.EqualError(yaml.Unmarshal(yml, &tested), errEmptyDestination)
	a.Equal(now, tested.MTime)

	// 2. check load invalid source file if "source" field is not empty
	yml = genYml(expectedDest, invalidSrc, "", testMode)
	a.Error(yaml.Unmarshal(yml, &tested))
	a.Equal(now, tested.MTime)

	// 3. check empty body error
	yml = genYml(expectedDest, "", "", testMode)
	a.EqualError(yaml.Unmarshal(yml, &tested), errBodyDataIsNotSet)
	a.Equal(now, tested.MTime)

	// 4. check file loading
	yml = genYml(expectedDest, expectedSrc, string(testBody), testMode)
	a.NoError(yaml.Unmarshal(yml, &tested))
	a.NotEqual(testBody, tested.Body())
	a.NotEqual(testMode, tested.Mode)
	a.NotEqual(now, tested.MTime)
	a.NotEqual(defaultOwner, tested.Owner)
	a.NotEqual(defaultGroup, tested.Group)

	// 5. check minimal body
	tested = File{
		MTime: now,
	}
	yml = genYml(expectedDest, "", string(testBody), testMode)
	a.NoError(yaml.Unmarshal(yml, &tested))
	a.Equal(testBody, tested.Body())
	a.Equal(testMode, tested.Mode)
	a.NotEqual(now, tested.MTime)
	a.NotEqual(defaultOwner, tested.Owner)
	a.NotEqual(defaultGroup, tested.Group)
}

func TestLoadFile(t *testing.T) {
	// TODO implementing here
	a := assert.New(t)

	// 1. check load fail
	tested, err := LoadFile("pkg.info/not-exist", "dest", "owner", "group")
	a.Nil(tested)
	a.Error(err)

	// 2. check normal load
	tcFile, _ := test.GetTestFilePath("pkg.info/example.sh")
	expectedBody, _ := ioutil.ReadFile(tcFile)
	expectedStat, _ := os.Stat(tcFile)
	tested, err = LoadFile("pkg.info/example.sh", "dest", "owner", "group")
	a.NotNil(tested)
	a.NoError(err)
	a.Equal(expectedBody, tested.body)
	a.Equal(uint(expectedStat.Mode()), tested.Mode)
	a.Equal("dest", tested.Dest)
	a.Equal("owner", tested.Owner)
	a.Equal("group", tested.Group)
}

func TestNewFile(t *testing.T) {
	a := assert.New(t)

	tested, err := NewFile("body", "", "owner", "group", 0000)
	a.Nil(tested)
	a.Error(err)

	tested, err = NewFile("body", "dest", "owner", "group", 0000)
	a.NotNil(tested)
	a.NoError(err)
	a.Equal([]byte("body"), tested.body)
	a.Equal(uint(0000), tested.Mode)
	a.Equal("dest", tested.Dest)
	a.Equal("owner", tested.Owner)
	a.Equal("group", tested.Group)
}

func TestNewDirectory(t *testing.T) {
	a := assert.New(t)

	tested, err := NewDirectory("", "owner", "group", 0000)
	a.Nil(tested)
	a.Error(err)

	tested, err = NewDirectory("dest", "owner", "group", 0000)
	a.NotNil(tested)
	a.NoError(err)
	a.Equal("dest", tested.Dest)
	a.Equal("owner", tested.Owner)
	a.Equal("group", tested.Group)
	a.Equal(uint(0000), tested.Mode)
}
