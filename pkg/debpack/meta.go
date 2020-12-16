package debpack

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

type Relations []string

func (r *Relations) Set(pkgString string) error {
	for _, relation := range *r {
		if relation == pkgString {
			return nil
		}
	}

	*r = append(*r, pkgString)

	return nil
}

type DEBMetaData struct {
	Package,
	Version,
	Maintainer,
	Summary,
	Description,
	Section,
	Priority string
	Essential bool
	Architecture,
	Origin,
	Bugs,
	Homepage string
	Tag    []string
	Source string
	Depends,
	PreDepends,
	Recommends,
	Suggests,
	Breaks,
	Conflicts,
	Replaces,
	Provides Relations
}

const metaHeader = `Package: %s
Version: %s
Architecture: %s
Maintainer: %s
Homepage: %s
Depends: %s
Description: %s
%s`

func (meta *DEBMetaData) depends() string {
	return strings.Join(meta.Depends, ", ")
}

func (meta *DEBMetaData) descriptionForm() string {
	return meta.Description
}

func (meta *DEBMetaData) description() string {
	out := bytes.NewBufferString("")
	in := bufio.NewReader(bytes.NewBufferString(meta.Description))

	for {
		if line, _, err := in.ReadLine(); err != io.EOF {
			_, _ = fmt.Fprintf(out, " %s\n", string(line))
		} else {
			break
		}
	}

	return out.String()
}

func (meta *DEBMetaData) MakeControl() []byte {
	return []byte(fmt.Sprintf(
		metaHeader,
		meta.Package,
		meta.Version,
		meta.Architecture,
		meta.Maintainer,
		meta.Homepage,
		meta.depends(),
		meta.Summary,
		meta.description(),
	))
}
