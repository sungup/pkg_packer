package pack

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type Script string

func (s *Script) load(source string) error {
	source = absSourcePath(source)

	body, err := ioutil.ReadFile(source)
	if err != nil {
		return fmt.Errorf(errLoadingFailed, err)
	}

	// append new line for next appended script
	*s = Script(append(body, '\n'))

	return nil
}

func (s *Script) UnmarshalYAML(value *yaml.Node) error {
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
		*s = Script(tV.Body + "\n")
	}

	return nil
}

func (s *Script) Append(line string) {
	*s += Script(line + "\n")
}

func (s *Script) String() string {
	return string(*s)
}

func (s *Script) Bytes() []byte {
	return []byte(*s)
}
