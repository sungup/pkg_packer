package test

import (
	"crypto/rand"
	"errors"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"testing"
	"time"
)

const (
	ExpectFile     = "test.yml"
	ExpectFileMode = uint(0644)
)

var (
	ExpectFileDate = time.Date(
		2020, 02, 29,
		01, 19, 47, 843167765,
		time.UTC,
	)
)

type YAMLTestData struct {
	Meta map[string]string `yaml:"meta"`
	Test map[string]string `yaml:"test"`
}

func RandString(t *testing.T) []byte {
	a := assert.New(t)

	defaults := make([]byte, 16)
	_, err := rand.Read(defaults)
	a.NoError(err)

	return defaults
}

func GetTestFilePath(filepath string) (string, error) {
	_, thisFilePath, _, _ := runtime.Caller(0)

	testPath := path.Join(path.Dir(thisFilePath), filepath)

	if _, err := os.Stat(testPath); os.IsNotExist(err) {
		return "", errors.New("test file doesn't exist")
	} else {
		return testPath, nil
	}
}

func LoadTestFile(filepath string) ([]byte, error) {
	var err error
	if filepath, err = GetTestFilePath(filepath); err != nil {
		return nil, err
	}

	return ioutil.ReadFile(filepath)
}

func LoadTestYAMLData(filepath string) (*YAMLTestData, error) {
	data := new(YAMLTestData)

	if buffer, err := LoadTestFile(filepath); err != nil {
		return nil, err
	} else {
		if err := yaml.Unmarshal(buffer, data); err != nil {
			return nil, err
		} else {
			return data, nil
		}
	}
}
