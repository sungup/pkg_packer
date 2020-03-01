package pkg

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
	"time"
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

func (pkg *Package) init() {
	pkg.Meta.UpdateBuildTime(time.Now().UTC())
}

func LoadPkgInfo(filepath string) (*Package, error) {
	buffer, err := ioutil.ReadFile(filepath)

	if err != nil {
		return nil, err
	}

	pkg := new(Package)

	if err := yaml.Unmarshal(buffer, pkg); err != nil {
		return nil, err
	}

	pkg.init()

	return pkg, nil
}

func NewPackage(meta PackageMeta) *Package {
	pkg := new(Package)

	pkg.Meta = meta

	pkg.Dirs = make([]PackageDir, 0)
	pkg.Files = make(map[string][]PackageFile)

	pkg.PreIn = make([]string, 0)
	pkg.PostIn = make([]string, 0)
	pkg.PreUn = make([]string, 0)
	pkg.PostUn = make([]string, 0)

	pkg.init()

	return pkg
}
