package pack

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type FileType string

const (
	GenericFile   = FileType("generic")
	ConfigFile    = FileType("config")
	DocumentFile  = FileType("doc")
	NotUseFile    = FileType("not_use")
	MissingOkFile = FileType("missing_ok")
	NoReplaceFile = FileType("no_replace")
	SpecFile      = FileType("spec")
	GhostFile     = FileType("ghost")
	LicenseFile   = FileType("license")
	ReadMeFile    = FileType("readme")
	ExcludeFile   = FileType("exclude")
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
	Meta `yaml:"meta"`

	Dirs  []*Directory         `yaml:"directory"`
	Files map[FileType][]*File `yaml:"files"`

	PreIn  script `yaml:"prein"`
	PostIn script `yaml:"postin"`

	PreUn  script `yaml:"preun"`
	PostUn script `yaml:"postun"`

	Dependencies []*Relation `yaml:"dependencies"`
}

func (pkg *Package) AddDirectory(pkgDir *Directory) {
	pkg.Dirs = append(pkg.Dirs, pkgDir)
}

func (pkg *Package) AddFile(fileType FileType, pkgFile *File) error {
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
		return nil, fmt.Errorf("loading package recipe file at %s: %w", filepath, err)
	}

	if err = yaml.Unmarshal(buffer, pkg); err != nil {
		return nil, fmt.Errorf("unexpected recipe file format: %w", err)
	}

	// TODO check file type in Files map

	return pkg, nil
}

func NewPackage(meta Meta) *Package {
	pkg := new(Package)

	pkg.Meta = meta

	pkg.Dirs = make([]*Directory, 0)
	pkg.Files = make(map[FileType][]*File)

	pkg.Dependencies = make([]*Relation, 0)

	return pkg
}
