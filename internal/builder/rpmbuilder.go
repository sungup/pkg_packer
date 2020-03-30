package builder

import (
	"errors"
	"fmt"
	"github.com/google/rpmpack"
	"github.com/sungup/pkg_packer/internal/pkg"
	"io"
)

type RPMBuilder struct {
	PackageBuilder

	pkgInfo *pkg.Package
}

func (rpm *RPMBuilder) rpmMetadata(meta pkg.PackageMeta) rpmpack.RPMMetaData {
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
		BuildTime:   meta.BuildTime(),
		Provides:    nil,
		Obsoletes:   nil,
		Suggests:    nil,
		Recommends:  nil,
		Requires:    nil,
		Conflicts:   nil,
	}
}

func (rpm *RPMBuilder) dirToRPMFile(info pkg.PackageDir) rpmpack.RPMFile {
	// Ignore MTime because the directories' modified time will be changed
	// because of their contents in directory
	return rpmpack.RPMFile{
		Name:  info.Dest,
		Mode:  info.Mode + 040000,
		Owner: info.Owner,
		Group: info.Group,
	}
}

func (rpm *RPMBuilder) fileToRPMFile(typeName string, info pkg.PackageFile) (rpmpack.RPMFile, error) {
	fileType := rpmpack.GenericFile

	// string to type
	switch typeName {
	case "generic":
		fileType = rpmpack.GenericFile
	case "config":
		// file types is a flag variable. so, config file type only keep the original config while removing
		// not updating. thus, if you keep the original config file while updating, add "NoReplaceFile" flag.
		fileType = rpmpack.ConfigFile | rpmpack.NoReplaceFile
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
		return rpmpack.RPMFile{},
			errors.New("unexpected file type: " + typeName)
	}

	return rpmpack.RPMFile{
		Name:  info.Dest,
		Body:  info.FileData(),
		Mode:  info.FileMode(),
		Owner: info.Owner,
		Group: info.Group,
		MTime: uint32(info.FileMTime().Unix()),
		Type:  fileType,
	}, nil
}

func (rpm *RPMBuilder) Filename() (string, error) {
	meta := rpm.pkgInfo.Meta

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

func (rpm *RPMBuilder) Build(writer io.Writer) error {
	var rpmPack *rpmpack.RPM
	var err error

	// 0. make rpm metadata
	if rpmPack, err = rpmpack.NewRPM(rpm.rpmMetadata(rpm.pkgInfo.Meta)); err != nil {
		return err
	}

	// 1. add directory list
	for _, dir := range rpm.pkgInfo.Dirs {
		rpmPack.AddFile(rpm.dirToRPMFile(dir))
	}

	// 2. add files
	for typeName, fList := range rpm.pkgInfo.Files {
		for _, fItem := range fList {
			if rpmFile, err := rpm.fileToRPMFile(typeName, fItem); err == nil {
				rpmPack.AddFile(rpmFile)
			} else {
				return err
			}
		}
	}

	// 3. add prein/postin/preun/postun
	rpmPack.AddPrein(rpm.pkgInfo.PreInScript())
	rpmPack.AddPostin(rpm.pkgInfo.PostInScript())
	rpmPack.AddPreun(rpm.pkgInfo.PreUnScript())
	rpmPack.AddPostun(rpm.pkgInfo.PostUnScript())

	// 4. add dependencies
	for _, dependency := range rpm.pkgInfo.Dependencies {
		if err := rpmPack.Requires.Set(dependency); err != nil {
			return err
		}
	}

	return rpmPack.Write(writer)
}

func NewRPMBuilder(pkgInfo *pkg.Package) *RPMBuilder {
	builder := new(RPMBuilder)

	builder.pkgInfo = pkgInfo

	return builder
}
