package info

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type script string

func (s *script) load(source string) error {
	source = absSourcePath(source)

	body, err := ioutil.ReadFile(source)
	if err != nil {
		return fmt.Errorf(errLoadingFailed, err)
	}

	// append new line for next appended script
	*s = script(append(body, '\n'))

	return nil
}

func (s *script) UnmarshalYAML(value *yaml.Node) error {
	tV := struct {
		Source string `yaml:"source"`
		Body   string `yaml:"body"`
	}{
		Body: emptyBodyCheckStr,
	}

	if err := value.Decode(&tV); err != nil {
		return fmt.Errorf(errYAMLDecodingFailed, err)
	}

	if tV.Source != "" {
		return s.load(tV.Source)
	} else if tV.Body != emptyBodyCheckStr {
		*s = script(tV.Body + "\n")
	}

	return nil
}

func (s script) MarshalYAML() (interface{}, error) {
	return map[string]string{"body": string(s)}, nil
}

func (s *script) Append(line string) {
	*s += script(line + "\n")
}

func (s *script) String() string {
	return string(*s)
}

func (s *script) Bytes() []byte {
	return []byte(*s)
}
