package pkg

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
)

type Package struct {
	Meta PackageMeta `yaml:"meta"`

	Dirs  []PackageDir             `yaml:"directory"`
	Files map[string][]PackageFile `yaml:"files"`

	PreIn  []string `yaml:"prein"`
	PostIn []string `yaml:"postin"`

	PreUn  []string `yaml:"preun"`
	PostUn []string `yaml:"postun"`
}

// Setters
func (pkg *Package) AppendPreIn(script string) {
	pkg.PreIn = append(pkg.PreIn, script+";")
}

func (pkg *Package) AppendPostIn(script string) {
	pkg.PostIn = append(pkg.PostIn, script+";")
}

func (pkg *Package) AppendPreUn(script string) {
	pkg.PreUn = append(pkg.PreUn, script+";")
}

func (pkg *Package) AppendPostUn(script string) {
	pkg.PostUn = append(pkg.PostUn, script+";")
}

// Getter
func (pkg *Package) PreInScript() string {
	return strings.Join(pkg.PreIn, "\n")
}

func (pkg *Package) PostInScript() string {
	return strings.Join(pkg.PostIn, "\n")
}

func (pkg *Package) PreUnScript() string {
	return strings.Join(pkg.PreUn, "\n")
}

func (pkg *Package) PostUnScript() string {
	return strings.Join(pkg.PostUn, "\n")
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
