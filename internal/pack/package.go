package pack

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"
)

var srcHome string

func UpdateSourceDir(dir string) string {
	orig := srcHome

	srcHome = dir

	return orig
}

func absSourcePath(filepath string) string {
	if filepath != "" && !strings.HasPrefix(filepath, "/") {
		return path.Join(srcHome, filepath)
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
	Meta Meta `yaml:"meta"`

	Dirs  []*Directory       `yaml:"directory"`
	Files map[string][]*File `yaml:"files"`

	PreIn  Script `yaml:"prein"`
	PostIn Script `yaml:"postin"`

	PreUn  Script `yaml:"preun"`
	PostUn Script `yaml:"postun"`

	Dependencies []*Relation `yaml:"dependencies"`
}

func (pkg *Package) init() {
	pkg.Meta.UpdateBuildTime(time.Now().UTC())
}

func (pkg *Package) AddDirectory(pkgDir *Directory) {
	pkg.Dirs = append(pkg.Dirs, pkgDir)
}

func (pkg *Package) AddFile(fileType string, pkgFile *File) error {
	// TODO check file type
	pkg.Files[fileType] = append(pkg.Files[fileType], pkgFile)

	return nil
}

func LoadPkgInfo(filepath string) (*Package, error) {
	var (
		buffer []byte
		err    error

		pkg = new(Package)
	)

	if buffer, err = ioutil.ReadFile(filepath); err != nil {
		return nil, err
	}

	if err = yaml.Unmarshal(buffer, pkg); err != nil {
		return nil, err
	}

	pkg.init()

	return pkg, nil
}

func NewPackage(meta Meta) *Package {
	pkg := new(Package)

	pkg.Meta = meta

	pkg.Dirs = make([]*Directory, 0)
	pkg.Files = make(map[string][]*File)

	pkg.Dependencies = make([]*Relation, 0)

	pkg.init()

	return pkg
}
