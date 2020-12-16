package debpack

import (
	"archive/tar"
	"time"
)

// FileType is a type identifier in debian package.
type FileType int32

const (
	Directory FileType = 1 << iota >> 1

	// GenericFile is just a basic file. if value is not specified, this file is a directory
	GenericFile

	// ConfigFile is a configuration file, and an existing file should be saved during a package
	// upgrade operation and not removed during a package removal time. This file will be stored
	// in conffile text based file in deb package archive.
	ConfigFile

	// ExcludeFile is not a part of the package, and should not be installed.
	ExcludeFile
)

type DEBFile struct {
	Name  string
	Body  []byte
	Mode  uint
	Owner string
	Group string
	MTime time.Time
	Type  FileType
}

func (deb *DEBFile) isDir() bool {
	return deb.Type == Directory
}

func (deb *DEBFile) isConfig() bool {
	return (deb.Type & ConfigFile) == ConfigFile
}

func (deb *DEBFile) tarTypeFlag() byte {
	if deb.isDir() {
		return tar.TypeDir
	} else {
		return tar.TypeReg
	}
}
