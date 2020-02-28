package builder

import (
	"github.com/stretchr/testify/assert"
	"github.com/sungup/pkg_packer/internal/pkg"
	"github.com/sungup/pkg_packer/test"
	"testing"
)

func LoadYAMLExpectAndPackage() (*test.YAMLTestData, *pkg.Package) {
	expect, _ := test.LoadTestYAMLData(test.ExpectFile)
	testPath, _ := test.GetTestFilePath(test.ExpectFile)
	tested, _ := pkg.LoadPkgInfo(testPath)

	return expect, tested
}

func Test_loadFile(t *testing.T) {
	a := assert.New(t)

	// Generate random string
	defaults := test.RandString(t)
	testFile, _ := test.GetTestFilePath(test.ExpectFile)
	invalidFile, _ := test.GetTestFilePath(test.ExpectFile + ".invalid")
	yamlData, _ := test.LoadTestFile(test.ExpectFile)

	// Check load existed file
	tested := loadFile(testFile, defaults)
	a.NotEmpty(tested)
	a.Equal(yamlData, tested)
	a.NotEqual(defaults, tested)

	// Check load default contents
	tested = loadFile(invalidFile, defaults)
	a.NotEmpty(tested)
	a.NotEqual(yamlData, tested)
	a.Equal(defaults, tested)

	// check empty contents
	tested = loadFile(invalidFile, nil)
	a.Empty(tested)
	a.NotEqual(yamlData, tested)
	a.NotEqual(defaults, tested)
}

func Test_fileStat(t *testing.T) {
	/*
		a := assert.New(t)

		defaultMode := 0777
		defaultMTime := time.Now().Unix()

		testFile,    _ := test.GetTestFilePath(test.ExpectFile)
		invalidFile, _ := test.GetTestFilePath(test.ExpectFile + ".invalid")

		testedMode, testedTime := fileStat(testFile, defaultMode, defaultMTime)

		testedMode, testedTime := fileStat(invalidFile, defaultMode, defaultMTime)
	*/
}
