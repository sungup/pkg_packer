package info

import (
	"github.com/stretchr/testify/assert"
	"github.com/sungup/pkg_packer/test"
	"gopkg.in/yaml.v3"
	"os"
	"path"
	"strings"
	"testing"
)

func TestScript_load(t *testing.T) {
	a := assert.New(t)

	tested := script("")

	// temporary change ownership to raise permission denied
	notReadable, _ := test.GetTestFilePath("pkg.info/not-readable")
	_ = os.Chmod(notReadable, 0000)
	defer func() { _ = os.Chmod(notReadable, 0644) }()

	failTc := []string{"not-exist", "dir-file", "not-readable"}
	successTc := []string{"example.sh"}

	for _, tc := range failTc {
		a.Error(tested.load(path.Join("pkg.info", tc)))
		a.Empty(tested)
	}

	for _, tc := range successTc {
		tcFile := path.Join("pkg.info", tc)
		expectedBody, _ := test.LoadTestFile(tcFile)
		expectedBody = append(expectedBody, '\n')

		a.NoError(tested.load(tcFile))
		a.Equal(expectedBody, tested.Bytes())
	}
}

func TestScript_UnmarshalYAML(t *testing.T) {
	a := assert.New(t)

	expectedBody := "#!/bin/bash\necho Hello, World! from body section\n"
	exampleStruct := struct {
		Source string `yaml:"source"`
		Body   string `yaml:"body" `
	}{Body: expectedBody}

	// 0. invalid yaml format
	{
		tested := script("")

		a.Error(yaml.Unmarshal([]byte("hello world"), &tested))
		a.Empty(tested)
	}

	// 1. Test normal body loading
	// 1-1. empty yaml
	{
		tested := script("")

		a.NoError(yaml.Unmarshal([]byte(""), &tested))
		a.Empty(tested)
	}

	// 1-2. filled body
	{
		tested := script("")
		yamlBody, _ := yaml.Marshal(exampleStruct)

		a.NoError(yaml.Unmarshal(yamlBody, &tested))
		a.Equal(expectedBody+"\n", tested.String())
	}

	// 2. Test file loading
	// 2-1. invalid path loading
	{
		tested := script("")
		exampleStruct.Source = "pkg.info/dir-file"
		yamlBody, _ := yaml.Marshal(exampleStruct)

		a.Error(yaml.Unmarshal(yamlBody, &tested))
		a.Empty(tested)
	}

	// 2-2. valid path loading
	{
		tested := script("")

		source := "pkg.info/example.sh"
		expectedLoaded, _ := test.LoadTestFile(source)
		expectedLoaded = append(expectedLoaded, '\n')

		exampleStruct.Source = source
		yamlBody, _ := yaml.Marshal(exampleStruct)

		a.NoError(yaml.Unmarshal(yamlBody, &tested))
		a.Equal(expectedLoaded, tested.Bytes())
	}
}

func TestScript_MarshalYAML(t *testing.T) {
	a := assert.New(t)

	expectedBody := "echo hello world"
	expectedYaml, _ := yaml.Marshal(map[string]string{"body": expectedBody})

	testScript := script(expectedBody)
	testedYaml, err := yaml.Marshal(&testScript)

	a.NoError(err)
	a.NotEqual(testScript, testedYaml)
	a.Contains(string(testedYaml), testScript)
	a.Equal(expectedYaml, testedYaml)
}

func TestScript_Append(t *testing.T) {
	a := assert.New(t)

	tcFile := path.Join("pkg.info", "example.sh")

	buffer, _ := test.LoadTestFile(tcFile)
	expectedBody := string(buffer)
	expectedBody = strings.Join(append([]string{expectedBody}, loremIpsum...), "\n") + "\n"

	tested := script("")
	a.NoError(tested.load(tcFile))

	for _, line := range loremIpsum {
		tested.Append(line)
	}

	a.Equal(expectedBody, tested.String())
}

func TestScript_String(t *testing.T) {
	// ignore this case because script.String is a type conversion simple getter
}

func TestScript_Bytes(t *testing.T) {
	// ignore this case because script.Bytes is a type conversion simple getter
}
