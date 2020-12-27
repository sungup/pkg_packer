package pack

import (
	"github.com/stretchr/testify/assert"
	"github.com/sungup/pkg_packer/test"
	"os"
	"path"
	"testing"
)

func TestUpdateSourceDir(t *testing.T) {
	a := assert.New(t)

	// initialized home directory is GetTestFileHome() which has been updated at init_test.go
	expectedHome := test.GetTestFileHome()
	expectedMoved, _ := test.GetTestFilePath("internal.pack/dir-file")

	a.Equal(expectedHome, UpdateSourceDir(expectedMoved))
	a.Equal(expectedMoved, UpdateSourceDir(expectedHome))
}

func TestAbsSourcePath(t *testing.T) {
	a := assert.New(t)

	cwd, _ := os.Getwd()
	osd := UpdateSourceDir(cwd)
	defer func(path string) { _ = UpdateSourceDir(path) }(osd)

	testCase := map[string]string{
		"":              "",
		"example.path":  path.Join(cwd, "example.path"),
		"/example.root": "/example.root",
	}

	for tc, expectedPath := range testCase {
		a.Equal(expectedPath, absSourcePath(tc))
	}
}
