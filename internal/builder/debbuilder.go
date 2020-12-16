package builder

import (
	"fmt"
	"github.com/sungup/pkg_packer/internal/info"
	"github.com/sungup/pkg_packer/pkg/debpack"
	"io"
)

type DEBBuilder struct {
	PackageBuilder

	pkgInfo *info.Package

	// conffiles is a contents containing the list of config file should keep.
	conffiles []string
}

// Metadata structure.
//  - control file
//  - conffile file
//  - md5sums file

func (deb *DEBBuilder) metadata(meta info.PackageMeta) debpack.DEBMetaData {
	return debpack.DEBMetaData{
		Package:      meta.Name,
		Version:      meta.Version,
		Maintainer:   meta.Maintainer,
		Summary:      meta.Summary,
		Description:  meta.Description,
		Architecture: deb.arch(),
		Homepage:     meta.URL,
	}
}

func (deb *DEBBuilder) dirToDEBFile(info info.PackageDir) debpack.DEBFile {
	return debpack.DEBFile{
		Name:  info.Dest,
		Mode:  info.Mode + 040000,
		Owner: info.Owner,
		Group: info.Group,
		Type:  debpack.Directory,
	}
}

func (deb *DEBBuilder) fileToDEBFile(typeName string, info info.PackageFile) (debpack.DEBFile, error) {
	fileType := debpack.GenericFile

	// string to type
	switch typeName {
	case "generic":
		fileType = debpack.GenericFile
	case "config":
		// file types is a flag variable. so, config file type only keep the original config while removing
		// not updating. thus, if you keep the original config file while updating, add "NoReplaceFile" flag.
		fileType = debpack.ConfigFile
	case "doc":
		fileType = debpack.GenericFile
	case "not_use":
		fileType = debpack.GenericFile
	case "missing_ok":
		fileType = debpack.GenericFile
	case "no_replace":
		fileType = debpack.ConfigFile

	case "spec":
		fileType = debpack.GenericFile
	case "ghost":
		fileType = debpack.GenericFile
	case "license":
		fileType = debpack.GenericFile
	case "readme":
		fileType = debpack.GenericFile
	case "exclude":
		fileType = debpack.ExcludeFile

	default:
		return debpack.DEBFile{},
			fmt.Errorf("unexpected file type: %s", typeName)
	}

	if (fileType & debpack.ConfigFile) != 0 {
		deb.conffiles = append(deb.conffiles, info.Dest)
	}

	return debpack.DEBFile{
		Name:  info.Dest,
		Body:  info.FileData(),
		Mode:  info.FileMode(),
		Owner: info.Owner,
		Group: info.Group,
		MTime: info.FileMTime(),
		Type:  fileType,
	}, nil
}

func (deb *DEBBuilder) arch() string {
	meta := deb.pkgInfo.Meta
	arch := "all"
	if meta.Arch != "" {
		switch meta.Arch {
		case "386":
			arch = "i386"
		case "x86_64":
			arch = "amd64"
		case "noarch":
			arch = "all"
		default:
			arch = meta.Arch
		}
	}

	return arch
}

// Filename returns the package's filename to store. Debian package's naming rule is defined at the
// https://www.debian.org/doc/manuals/debian-faq/pkg-basics.en.html. Also, pkg-packer will focus on
// the unified build environment using 1 package recipe, so remove the DebianRevisionNumber field.
func (deb *DEBBuilder) Filename() (string, error) {
	meta := deb.pkgInfo.Meta

	if meta.Name == "" {
		return "", fmt.Errorf("undefined package name")
	} else if meta.Version == "" {
		return "", fmt.Errorf("undefined package version")
	}

	return fmt.Sprintf("%s_%s_%s.deb", meta.Name, meta.Version, deb.arch()), nil
}

func (deb *DEBBuilder) Build(writer io.Writer) error {
	var (
		debPack *debpack.DEB
		err     error
	)

	if debPack, err = debpack.NewDEB(deb.metadata(deb.pkgInfo.Meta)); err != nil {
		return err
	}

	for _, dir := range deb.pkgInfo.Dirs {
		debPack.AddFile(deb.dirToDEBFile(dir))
	}

	for typeName, fList := range deb.pkgInfo.Files {
		for _, file := range fList {
			if debFile, err := deb.fileToDEBFile(typeName, file); err == nil {
				debPack.AddFile(debFile)
			} else {
				return err
			}
		}
	}

	debPack.AddPrein(deb.pkgInfo.PreInScript())
	debPack.AddPostin(deb.pkgInfo.PostInScript())
	debPack.AddPreun(deb.pkgInfo.PreUnScript())
	debPack.AddPostun(deb.pkgInfo.PostUnScript())

	for _, dependency := range deb.pkgInfo.Dependencies {
		if err := debPack.Depends.Set(dependency.DebFormat()); err != nil {
			return err
		}
	}

	return debPack.Write(writer)
}

func NewDEBBuilder(pkgInfo *info.Package) PackageBuilder {
	return &DEBBuilder{
		pkgInfo: pkgInfo,
	}
}
