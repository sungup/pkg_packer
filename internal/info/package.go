package info

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"path"
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

	Dependencies []Dependency `yaml:"dependencies"`

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

func (pkg *Package) joinedFilePath(filePath string) string {
	if filePath != "" && !strings.HasPrefix(filePath, "/") {
		return path.Join(pkg.srcHome, filePath)
	} else {
		return filePath
	}
}

func (pkg *Package) init() {
	pkg.Meta.UpdateBuildTime(time.Now().UTC())

	// update all pkg directory which have src directory
	for key, fileList := range pkg.Files {
		for idx, file := range fileList {
			pkg.Files[key][idx].Src = pkg.joinedFilePath(file.Src)
		}
	}
}

func (pkg *Package) AddDirectory(pkgDir PackageDir) {
	pkg.Dirs = append(pkg.Dirs, pkgDir)
}

func (pkg *Package) AddFile(fileType string, pkgFile PackageFile) error {
	pkgFile.Src = pkg.joinedFilePath(pkgFile.Src)

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

	pkg.Dirs = make([]PackageDir, 0)
	pkg.Files = make(map[string][]PackageFile)

	pkg.PreIn = make([]string, 0)
	pkg.PostIn = make([]string, 0)
	pkg.PreUn = make([]string, 0)
	pkg.PostUn = make([]string, 0)

	pkg.Dependencies = make([]Dependency, 0)

	pkg.srcHome = srcHome

	pkg.init()

	return pkg
}
