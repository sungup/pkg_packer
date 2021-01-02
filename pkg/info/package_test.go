package info

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/sungup/pkg_packer/test"
	"gopkg.in/yaml.v3"
	"os"
	"path"
	"testing"
)

func TestUpdateSourceDir(t *testing.T) {
	a := assert.New(t)

	// initialized home directory is GetTestFileHome() which has been updated at init_test.go
	expectedHome := test.GetTestFileHome()
	expectedMoved, _ := test.GetTestFilePath("pkg.info/dir-file")

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

func TestFileType_UnmarshalYAML(t *testing.T) {
	a := assert.New(t)

	tested := struct {
		Type fileType `yaml:"test"`
	}{}

	testYmlForm := "test: %s"

	// 1. decoding error
	a.Error(yaml.Unmarshal([]byte("test"), &tested))

	// 2. test not existed file type
	testType := "not_exist"
	expectedErr := fmt.Sprintf(errMsgUnsupportedFileType, testType)
	ymlBody := fmt.Sprintf(testYmlForm, testType)
	a.EqualError(yaml.Unmarshal([]byte(ymlBody), &tested), expectedErr)

	// 3. check support file type
	for expectedType, testType := range fileTypeToName {
		ymlBody := fmt.Sprintf(testYmlForm, testType)
		a.NoError(yaml.Unmarshal([]byte(ymlBody), &tested))
		a.Equal(expectedType, tested.Type)
	}
}

func TestFileType_MarshalYAML(t *testing.T) {
	a := assert.New(t)

	fileTypeTest := struct {
		Type fileType `yaml:"test"`
	}{}

	testYmlForm := "test: %s\n"

	for testType, expectedType := range fileTypeToName {
		expectedBody := fmt.Sprintf(testYmlForm, expectedType)
		fileTypeTest.Type = testType

		tested, err := yaml.Marshal(fileTypeTest)
		a.NoError(err)
		a.Equal(expectedBody, string(tested))
	}
}

func TestFileType_String(t *testing.T) {
	a := assert.New(t)

	for tested, expected := range fileTypeToName {
		a.Equal(expected, tested.String())
	}
}

func TestPackage_AddDirectory(t *testing.T) {
	// skip function test because Package.AddDirectory is simple setter function.
}

func TestPackage_AddFile(t *testing.T) {
	// skip function test because Package.AddFile is simple setter function.
}

func TestPackage_AddDependency(t *testing.T) {
	a := assert.New(t)

	{
		tested := Package{}
		a.NoError(tested.AddDependency(expectedPkg))
		a.NotEmpty(tested.Dependencies)
	}

	for _, tc := range referenceInvalidOps {
		tested := Package{}
		a.Error(tested.AddDependency(fmt.Sprintf("%s %s %s", expectedPkg, tc, expectedVer)))
		a.Empty(tested.Dependencies)
	}

	for _, tc := range referenceValidOps {
		tested := Package{}
		a.NoError(tested.AddDependency(fmt.Sprintf("%s %s %s", expectedPkg, tc, expectedVer)))
		a.NotEmpty(tested.Dependencies)
	}
}

func TestLoadPkgInfo(t *testing.T) {
	a := assert.New(t)
	buffer := make(map[string]interface{})
	body, _ := yaml.Marshal(&referencePackage)
	_ = yaml.Unmarshal(body, buffer)

	// 1. Not exist test
	{
		testFilePath, _ := test.GetTestFilePath(path.Join("temp", "not-exist.yml"))
		pkg, err := LoadPkgInfo(testFilePath)
		a.Nil(pkg)
		a.Error(err)
	}

	// 2. Not readable test
	{
		testFilePath, _ := test.GetTestFilePath(path.Join("pkg.info", "not-readable"))
		_ = os.Chmod(testFilePath, 0000)

		pkg, err := LoadPkgInfo(testFilePath)
		a.Nil(pkg)
		a.Error(err)

		_ = os.Chmod(testFilePath, 0644)
	}

	for name := range nameToFileType {
		// 1. create invalid yaml file
		buffer["files"] = map[string][]map[string]interface{}{name + "_unexpected": referenceFileList}
		testFilePath, _ := test.GetTestFilePath(path.Join("temp", name+"_unexpected.yml"))
		testFilePtr, _ := os.Create(testFilePath)
		testEncoder := yaml.NewEncoder(testFilePtr)
		_ = testEncoder.Encode(buffer)
		_ = testEncoder.Close()
		_ = testFilePtr.Close()

		tested, err := LoadPkgInfo(testFilePath)
		a.Nil(tested)
		a.Error(err)

		_ = os.Remove(testFilePath)

		a.True(true)
		// 2. create valid yaml file
		buffer["files"] = map[string][]map[string]interface{}{name: referenceFileList}
		testFilePath, _ = test.GetTestFilePath(path.Join("temp", name+".yml"))
		testFilePtr, _ = os.Create(testFilePath)
		testEncoder = yaml.NewEncoder(testFilePtr)
		_ = testEncoder.Encode(buffer)
		_ = testEncoder.Close()
		_ = testFilePtr.Close()

		tested, err = LoadPkgInfo(testFilePath)
		a.NoError(err)
		a.NotNil(tested)

		a.Equal(referencePackage.Meta, tested.Meta)
		a.Equal(referencePackage.Dirs, tested.Dirs)
		a.Equal(referencePackage.Dependencies, tested.Dependencies)
		a.Len(tested.Files, 1)

		_ = os.Remove(testFilePath)
	}
}

func TestNewPackage(t *testing.T) {
	a := assert.New(t)

	tested := NewPackage(referencePackage.Meta)

	a.Equal(referencePackage.Meta, tested.Meta)
	a.NotNil(tested.Dirs)
	a.NotNil(tested.Files)
	a.NotNil(tested.Dependencies)
	a.Empty(tested.Dirs)
	a.Empty(tested.Files)
	a.Empty(tested.Dependencies)
}
