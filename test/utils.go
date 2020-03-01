package test

import (
	crand "crypto/rand"
	"errors"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	mrand "math/rand"
	"os"
	"path"
	"runtime"
	"testing"
	"time"
)

const (
	ExpectFile     = "test.yml"
	ExpectFileMode = uint(0644)

	randStrMin = 5
	randStrMax = 10
)

var (
	ExpectFileDate = time.Date(
		2020, 03, 01,
		05, 12, 36, 926074319,
		time.UTC,
	)

	randStringMap = []byte("abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ-_.")
)

type YAMLTestData struct {
	Meta map[string]string `yaml:"meta"`
	Test map[string]string `yaml:"test"`
}

func RandBytes(t *testing.T) []byte {
	a := assert.New(t)

	defaults := make([]byte, 16)
	_, err := crand.Read(defaults)
	a.NoError(err)

	for i := 0; i < len(defaults); i++ {
		defaults[i] = randStringMap[int(defaults[i])%len(randStringMap)]
	}

	return defaults
}

func RandString(t *testing.T) string {
	return string(RandBytes(t))
}

func RandStrings(t *testing.T) []string {
	buffer := make([]string, 0)

	seed := mrand.NewSource(time.Now().UnixNano())
	items := int(randStrMin + mrand.New(seed).Int31n(randStrMax-randStrMin))

	for i := 0; i < items; i++ {
		buffer = append(buffer, RandString(t))
	}

	return buffer
}

func GetTestFilePath(filepath string) (string, error) {
	_, thisFilePath, _, _ := runtime.Caller(0)

	testPath := path.Join(path.Dir(thisFilePath), filepath)

	if _, err := os.Stat(testPath); os.IsNotExist(err) {
		return testPath, errors.New("test file doesn't exist")
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
