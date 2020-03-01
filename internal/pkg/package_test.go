package pkg

import (
	"github.com/stretchr/testify/assert"
	"github.com/sungup/pkg_packer/test"
	"strings"
	"testing"
)

func createTestScripts(t *testing.T) ([]string, string) {
	scripts := test.RandStrings(t)

	body := make([]string, 0)

	for _, item := range scripts {
		body = append(body, item+";")
	}

	return scripts, strings.Join(body, "\n")
}

func TestLoadPkgInfo(t *testing.T) {
	// Prepare expect data structure
	expect, err := test.LoadTestYAMLData(test.ExpectFile)

	assert.NoErrorf(t, err, "reference set loading fail")

	// 1. LoadPkgInfo test
	expectPath, _ := test.GetTestFilePath(test.ExpectFile)
	tested, err := LoadPkgInfo(expectPath)

	assert.NoErrorf(t, err, "unexpected load error")

	// 2. PackageMeta Test
	compareMetaData(t, expect.Meta, &tested.Meta)
}

func TestPackage_PreInScript(t *testing.T) {
	a := assert.New(t)
	scripts, expect := createTestScripts(t)

	pkg := Package{PreIn: make([]string, 0)}

	for _, item := range scripts {
		pkg.AppendPreIn(item)
	}

	tested := pkg.PreInScript()

	a.Equal(expect, tested)
}

func TestPackage_PostInScript(t *testing.T) {
	a := assert.New(t)
	scripts, expect := createTestScripts(t)

	pkg := Package{PostIn: make([]string, 0)}

	for _, item := range scripts {
		pkg.AppendPostIn(item)
	}

	tested := pkg.PostInScript()

	a.Equal(expect, tested)
}

func TestPackage_PreUnScript(t *testing.T) {
	a := assert.New(t)
	scripts, expect := createTestScripts(t)

	pkg := Package{PreUn: make([]string, 0)}

	for _, item := range scripts {
		pkg.AppendPreUn(item)
	}

	tested := pkg.PreUnScript()

	a.Equal(expect, tested)
}

func TestPackage_PostUnScript(t *testing.T) {
	a := assert.New(t)
	scripts, expect := createTestScripts(t)

	pkg := Package{PostUn: make([]string, 0)}

	for _, item := range scripts {
		pkg.AppendPostUn(item)
	}

	tested := pkg.PostUnScript()

	a.Equal(expect, tested)
}
