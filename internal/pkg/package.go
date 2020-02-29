package pkg

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Package struct {
	Meta PackageMeta `yaml:"meta"`

	Dirs  []PackageDir             `yaml:"directory"`
	Files map[string][]PackageFile `yaml:"files"`
}

func LoadPkgInfo(filepath string) (*Package, error) {
	buffer, err := ioutil.ReadFile(filepath)

	if err != nil {
		return nil, err
	}

	pack := new(Package)

	if err := yaml.Unmarshal(buffer, pack); err != nil {
		return nil, err
	}

	return pack, nil
}
