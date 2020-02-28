package pkg

import (
	"testing"
)

func TestLoadPkgInfo(t *testing.T) {
	// Loading fail will be checked at LoadYAMLExpectPackage
	expect, tested := LoadYAMLExpectAndPackage(t)

	// 1. PackageMeta Test
	CompareMetaData(t, expect.Meta, &tested.Meta)
}
