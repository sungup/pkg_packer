package pack

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"time"
)

const (
	defaultOwner = "root"
	defaultGroup = "root"
	defaultMode  = 0644

	emptyBodyCheckStr = "pKGpACKEReMPTYbODYcHECKsTRING"

	errLoadingFailed      = "file loading failed: %w"
	errYAMLDecodingFailed = "loading file configuration failed: %w"
	errEmptyDestination   = "empty destination path"
	errBodyDataIsNotSet   = "package file body not yet described"
	errSourceFileIsDir    = "file loading failed: %s is a directory"
)

type File struct {
	Dest  string    `yaml:"dest"`
	Owner string    `yaml:"owner"`
	Group string    `yaml:"group"`
	Mode  uint      `yaml:"mode"`
	MTime time.Time `yaml:"mtime"`

	body []byte
}

func (f *File) Body() []byte {
	return f.body
}

func (f *File) load(source string) error {
	// to absolute path
	source = absSourcePath(source)

	stat, err := os.Stat(source)
	if err != nil {
		return fmt.Errorf(errLoadingFailed, err)
	} else if stat.IsDir() {
		return fmt.Errorf(errSourceFileIsDir, source)
	}

	body, err := ioutil.ReadFile(source)
	if err != nil {
		return fmt.Errorf(errLoadingFailed, err)
	}

	f.body = body
	f.Mode = uint(stat.Mode())
	f.MTime = stat.ModTime()

	return nil
}

func (f *File) UnmarshalYAML(value *yaml.Node) error {
	type tT File
	tV := struct {
		tT     `yaml:",inline"`
		Source string `yaml:"source"`
		Body   string `yaml:"body"`
	}{
		tT: tT{
			Owner: defaultOwner,
			Group: defaultGroup,
			Mode:  defaultMode,
			MTime: time.Now(),
		},
		Body: emptyBodyCheckStr,
	}

	// 0. yaml decoding
	if err := value.Decode(&tV); err != nil {
		return fmt.Errorf(errYAMLDecodingFailed, err)
	}

	tF := File(tV.tT)

	// 1. check empty body
	if tF.Dest == "" {
		return fmt.Errorf(errEmptyDestination)
	}

	// 2. load or set file data
	if tV.Source != "" {
		if err := tF.load(tV.Source); err != nil {
			return err
		}
	} else {
		if tV.Body == emptyBodyCheckStr {
			return fmt.Errorf(errBodyDataIsNotSet)
		}

		tF.body = []byte(tV.Body)
	}

	*f = tF

	return nil
}

func LoadFile(source string, dest, owner, group string) (*File, error) {
	f := File{
		Dest:  dest,
		Owner: owner,
		Group: group,
	}

	if dest == "" {
		return nil, fmt.Errorf(errEmptyDestination)
	}

	if err := f.load(source); err == nil {
		return &f, nil
	} else {
		return nil, err
	}
}

func NewFile(body string, dest, owner, group string, mode uint) (*File, error) {
	if dest == "" {
		return nil, fmt.Errorf(errEmptyDestination)
	}

	return &File{
		Dest:  dest,
		Mode:  mode,
		Owner: owner,
		Group: group,
		MTime: time.Now(),
		body:  []byte(body),
	}, nil
}

type Directory struct {
	Dest  string `yaml:"dest,omitempty"`
	Owner string `yaml:"owner"`
	Group string `yaml:"group"`
	Mode  uint   `yaml:"mode"`
}

func NewDirectory(dest, owner, group string, mode uint) (*Directory, error) {
	if dest == "" {
		return nil, fmt.Errorf("empty destination path")
	}

	return &Directory{
		Dest:  dest,
		Owner: owner,
		Group: group,
		Mode:  mode,
	}, nil
}
