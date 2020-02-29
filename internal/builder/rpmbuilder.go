package builder

import (
	"errors"
	"fmt"
	"github.com/google/rpmpack"
	"github.com/sungup/pkg_packer/internal/pkg"
	"time"
)

type RPMBuilder struct {
	PackageBuilder
}

func (rpm RPMBuilder) Metadata(meta *pkg.PackageMeta) rpmpack.RPMMetaData {
	return rpmpack.RPMMetaData{
		Name:        meta.Name,
		Summary:     meta.Summary,
		Description: meta.Description,
		Version:     meta.Version,
		Release:     meta.Release,
		Arch:        meta.Arch,
		OS:          meta.OS,
		Vendor:      meta.Vendor,
		URL:         meta.URL,
		Packager:    "",
		Group:       "",
		Licence:     meta.License,
		BuildHost:   "",
		Compressor:  "",
		Epoch:       0,
		BuildTime:   time.Now(),
		Provides:    nil,
		Obsoletes:   nil,
		Suggests:    nil,
		Recommends:  nil,
		Requires:    nil,
		Conflicts:   nil,
	}
}

func (rpm RPMBuilder) File(typeName string, info *pkg.PackageFile) rpmpack.RPMFile {
	fileType := rpmpack.GenericFile

	// string to type
	switch typeName {
	case "config":
		fileType = rpmpack.ConfigFile
	case "doc":
		fileType = rpmpack.DocFile
	case "not_use":
		fileType = rpmpack.DoNotUseFile
	case "missing_ok":
		fileType = rpmpack.MissingOkFile
	case "no_replace":
		fileType = rpmpack.NoReplaceFile
	case "spec":
		fileType = rpmpack.SpecFile
	case "ghost":
		fileType = rpmpack.GhostFile
	case "license":
		fileType = rpmpack.LicenceFile
	case "readme":
		fileType = rpmpack.ReadmeFile
	case "exclude":
		fileType = rpmpack.ExcludeFile
	default:
		fileType = rpmpack.GenericFile
	}

	return rpmpack.RPMFile{
		Name:  info.Dest,
		Body:  info.FileData(),
		Mode:  info.FileMode(),
		Owner: info.Owner,
		Group: info.Group,
		MTime: uint32(info.FileMTime().Unix()),
		Type:  fileType,
	}
}

func (rpm RPMBuilder) Filename(meta *pkg.PackageMeta) (string, error) {
	if meta.Name == "" {
		return "", errors.New("undefined package name")
	} else if meta.Version == "" {
		return "", errors.New("undefined package version")
	}

	release := ""
	if meta.Release != "" {
		release = "." + meta.Release
	}

	arch := "noarch"
	if meta.Arch != "" {
		arch = meta.Arch
	}

	return fmt.Sprintf(
		"%s-%s%s.%s.rpm",
		meta.Name, meta.Version, release, arch,
	), nil
}

func (rpm RPMBuilder) Build(_ *pkg.Package) error {
	// TODO Implementing here
	return errors.New("build function not yet implemented")
}
