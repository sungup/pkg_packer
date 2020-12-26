package pack

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"
)

var (
	sourceHome string
)

func UpdateSourceDir(dir string) {
	sourceHome = dir
}

func absSourcePath(filepath string) string {
	if filepath != "" && !strings.HasPrefix(filepath, "/") {
		return path.Join(sourceHome, filepath)
	} else {
		return filepath
	}
}

func init() {
	if dir, err := os.Getwd(); err == nil {
		UpdateSourceDir(dir)
	} else {
		UpdateSourceDir(".")
	}
}

type Package struct {
	Meta PackageMeta `yaml:"meta"`

	Dirs  []directory       `yaml:"directory"`
	Files map[string][]file `yaml:"files"`

	PreIn  []string `yaml:"prein"`
	PostIn []string `yaml:"postin"`

	PreUn  []string `yaml:"preun"`
	PostUn []string `yaml:"postun"`

	Dependencies []Relation `yaml:"dependencies"`

	srcHome string
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

func (pkg *Package) AddDirectory(pkgDir directory) {
	pkg.Dirs = append(pkg.Dirs, pkgDir)
}

func (pkg *Package) AddFile(fileType string, pkgFile file) error {
	// TODO check file type
	pkg.Files[fileType] = append(pkg.Files[fileType], pkgFile)

	return nil
}

func LoadPkgInfo(filepath string, srcHome string) (*Package, error) {
	buffer, err := ioutil.ReadFile(filepath)

	if err != nil {
		return nil, err
	}

	pkg := new(Package)

	if err := yaml.Unmarshal(buffer, pkg); err != nil {
		return nil, err
	}

	pkg.srcHome = srcHome

	pkg.init()

	return pkg, nil
}

func NewPackage(meta PackageMeta, srcHome string) *Package {
	pkg := new(Package)

	pkg.Meta = meta

	pkg.Dirs = make([]directory, 0)
	pkg.Files = make(map[string][]file)

	pkg.PreIn = make([]string, 0)
	pkg.PostIn = make([]string, 0)
	pkg.PreUn = make([]string, 0)
	pkg.PostUn = make([]string, 0)

	pkg.Dependencies = make([]Relation, 0)

	pkg.srcHome = srcHome

	pkg.init()

	return pkg
}
