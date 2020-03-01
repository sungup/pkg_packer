package pkg

import (
	"github.com/stretchr/testify/assert"
	"github.com/sungup/pkg_packer/test"
	"testing"
	"time"
)

func TestPackageFile_FileData(t *testing.T) {
	a := assert.New(t)

	testPath, _ := test.GetTestFilePath(test.ExpectFile)
	invalidPath := testPath + ".invalid"

	defaultBody := test.RandBytes(t)
	loadedBody, _ := test.LoadTestFile(test.ExpectFile)

	pkg := PackageFile{}
	var testedBody []byte

	// 1. load data from file
	pkg.Src = testPath
	pkg.Body = string(defaultBody)

	testedBody = pkg.FileData()
	a.NotEmpty(testedBody)
	a.Equal(loadedBody, testedBody)
	a.NotEqual(defaultBody, testedBody)

	// 2. load data from default body
	pkg.Src = invalidPath
	pkg.Body = string(defaultBody)

	testedBody = pkg.FileData()
	a.NotEmpty(testedBody)
	a.NotEqual(loadedBody, testedBody)
	a.Equal(defaultBody, testedBody)

	// 3. load empty data
	pkg.Src = invalidPath
	pkg.Body = ""

	testedBody = pkg.FileData()
	a.Empty(testedBody)
	a.NotEqual(loadedBody, testedBody)
	a.NotEqual(defaultBody, testedBody)
}

func TestPackageFile_FileMode(t *testing.T) {
	a := assert.New(t)

	testPath, _ := test.GetTestFilePath(test.ExpectFile)
	invalidPath := ""

	defaultMode := uint(0111) // --x--x--x permission
	loadedMode := test.ExpectFileMode

	pkg := PackageFile{Mode: defaultMode}
	var testedMode uint

	// 1. retrieve file mode from target file
	pkg.Src = testPath

	testedMode = pkg.FileMode()
	a.Equal(loadedMode, testedMode)
	a.NotEqual(defaultMode, testedMode)

	// 2. retrieve file mode from default mode
	pkg.Src = invalidPath

	testedMode = pkg.FileMode()
	a.NotEqual(loadedMode, testedMode)
	a.Equal(defaultMode, testedMode)
}

func TestPackageFile_FileMTime(t *testing.T) {
	a := assert.New(t)

	testPath, _ := test.GetTestFilePath(test.ExpectFile)
	invalidPath := ""

	defaultMTime := time.Now()
	loadedMTime := test.ExpectFileDate

	pkg := PackageFile{MTime: defaultMTime}
	var testedMTime time.Time

	// 1. retrieve file mode from target file
	pkg.Src = testPath

	testedMTime = pkg.FileMTime()
	// a.Equal(loadedMTime.Unix(), testedMTime.Unix())
	a.LessOrEqual(loadedMTime.Unix(), testedMTime.Unix())
	a.NotEqual(defaultMTime.Unix(), testedMTime.Unix())

	// 2. retrieve file mode from default mode
	pkg.Src = invalidPath

	testedMTime = pkg.FileMTime()
	a.NotEqual(loadedMTime.Unix(), testedMTime.Unix())
	a.Equal(defaultMTime.Unix(), testedMTime.Unix())
}
