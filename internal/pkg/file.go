package pkg

import (
	"io/ioutil"
	"os"
	"time"
)

type PackageFile struct {
	Dest  string    `yaml:"dest,omitempty"`
	Src   string    `yaml:"source"`
	Body  string    `yaml:"body"`
	Mode  uint      `yaml:"mode"`
	Owner string    `yaml:"owner"`
	Group string    `yaml:"group"`
	MTime time.Time `yaml:"mtime"`
}

func (file *PackageFile) FileData() []byte {
	if data, err := ioutil.ReadFile(file.Src); err == nil {
		return data
	} else if file.Body != "" {
		return []byte(file.Body)
	} else {
		return make([]byte, 0)
	}
}

func (file *PackageFile) FileMode() uint {
	if stat, err := os.Stat(file.Src); err == nil && !stat.IsDir() {
		return uint(stat.Mode())
	} else {
		return file.Mode
	}
}

func (file *PackageFile) FileMTime() time.Time {
	if stat, err := os.Stat(file.Src); err == nil && !stat.IsDir() {
		return stat.ModTime()
	} else {
		return file.MTime
	}
}

type PackageDir struct {
	Dest  string `yaml:"dest,omitempty"`
	Mode  uint   `yaml:"mode"`
	Owner string `yaml:"owner"`
	Group string `yaml:"group"`
}
