package builder

import (
	"github.com/sungup/pkg_packer/internal/pkg"
)

type PackageBuilder interface {
	Filename(meta *pkg.PackageMeta) (string, error)
	Build(info *pkg.Package) error
}
