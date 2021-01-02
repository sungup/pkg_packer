package info

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

const (
	GenericFile fileType = iota
	ConfigFile
	DocumentFile
	NotUseFile
	MissingOkFile
	NoReplaceFile
	SpecFile
	GhostFile
	LicenseFile
	ReadMeFile
	ExcludeFile
)

const (
	errMsgUnsupportedFileType = "unsupported pkg_packer file type: %v"
)

var (
	fileTypeToName map[fileType]string
	nameToFileType map[string]fileType

	srcHome string
)

func init() {
	// create string => fileType map for Unmarshal
	nameToFileType = map[string]fileType{
		"generic":    GenericFile,
		"config":     ConfigFile,
		"document":   DocumentFile,
		"not_use":    NotUseFile,
		"missing_ok": MissingOkFile,
		"no_replace": NoReplaceFile,
		"spec":       SpecFile,
		"ghost":      GhostFile,
		"license":    LicenseFile,
		"readme":     ReadMeFile,
		"exclude":    ExcludeFile,
	}

	// create fileType => string map
	fileTypeToName = make(map[fileType]string)
	for s, t := range nameToFileType {
		fileTypeToName[t] = s
	}

	if dir, err := os.Getwd(); err == nil {
		UpdateSourceDir(dir)
	} else {
		UpdateSourceDir(".")
	}
}

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

type fileType uint32

func (f *fileType) UnmarshalYAML(value *yaml.Node) error {
	var name string

	if err := value.Decode(&name); err != nil {
		// decoding fail
		return err
	} else if fT, ok := nameToFileType[name]; !ok {
		// unexpected name
		return fmt.Errorf(errMsgUnsupportedFileType, name)
	} else {
		// unmarshal success
		*f = fT
		return nil
	}
}

func (f fileType) MarshalYAML() (interface{}, error) {
	return fileTypeToName[f], nil
}

func (f fileType) String() string {
	return fileTypeToName[f]
}

type Package struct {
	Meta Meta `yaml:"meta"`

	Dirs  []*Directory         `yaml:"directory"`
	Files map[fileType][]*File `yaml:"files"`

	PreIn  script `yaml:"prein"`
	PostIn script `yaml:"postin"`

	PreUn  script `yaml:"preun"`
	PostUn script `yaml:"postun"`

	Dependencies []*Relation `yaml:"dependencies"`
}

func (pkg *Package) AddDirectory(pkgDir *Directory) {
	pkg.Dirs = append(pkg.Dirs, pkgDir)
}

func (pkg *Package) AddFile(fType fileType, pkgFile *File) error {
	pkg.Files[fType] = append(pkg.Files[fType], pkgFile)

	return nil
}

func (pkg *Package) AddDependency(info string) error {
	if dependency, err := NewRelation(info); err != nil {
		return err
	} else {
		pkg.Dependencies = append(pkg.Dependencies, dependency)
		return nil
	}
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

	return pkg, nil
}

func NewPackage(meta Meta) *Package {
	return &Package{
		Meta:         meta,
		Dirs:         make([]*Directory, 0),
		Files:        make(map[fileType][]*File),
		Dependencies: make([]*Relation, 0),
	}
}
