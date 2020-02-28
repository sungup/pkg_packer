package builder

import (
	"github.com/sungup/pkg_packer/internal/pkg"
	"io/ioutil"
	"os"
	"time"
)

type PackageBuilder interface {
	Filename(meta *pkg.PackageMeta) (string, error)
	Build(info *pkg.Package) error
}

func loadFile(source string, defaults []byte) []byte {
	if data, err := ioutil.ReadFile(source); err == nil {
		return data
	} else if defaults != nil {
		return defaults
	} else {
		return make([]byte, 0)
	}
}

func fileStat(source string, mode uint, mtime time.Time) (uint, time.Time) {
	if stat, err := os.Stat(source); os.IsExist(err) && !stat.IsDir() {
		return uint(stat.Mode()), stat.ModTime()
	} else {
		return mode, mtime
	}
}
