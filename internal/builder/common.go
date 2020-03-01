package builder

import "io"

type PackageBuilder interface {
	Filename() (string, error)
	Build(rpmPath io.Writer) error
}
