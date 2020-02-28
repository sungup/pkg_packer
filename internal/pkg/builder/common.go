package builder

import "github.com/sungup/pkg_packer/internal/pkg"

type PackageBuilder interface {
	Filename(meta *pkg.PackageMeta) (string, error)
	Build(meta *pkg.PackageMeta) error
}
