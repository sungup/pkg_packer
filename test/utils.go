package test

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
	"runtime"
)

const (
	ExpectFile = "test.yml"
)

type YAMLTestData struct {
	Meta map[string]string `yaml:"meta"`
	Test map[string]string `yaml:"test"`
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
