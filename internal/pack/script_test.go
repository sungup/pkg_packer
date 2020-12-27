package pack

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

	tested := Script("")

	// temporary change ownership to raise permission denied
	notReadable, _ := test.GetTestFilePath("internal.pack/not-readable")
	_ = os.Chmod(notReadable, 0000)
	defer func() { _ = os.Chmod(notReadable, 0644) }()

	failTc := []string{"not-exist", "dir-file", "not-readable"}
	successTc := []string{"example.sh"}

	for _, tc := range failTc {
		a.Error(tested.load(path.Join("internal.pack", tc)))
		a.Empty(tested)
	}

	for _, tc := range successTc {
		tcFile := path.Join("internal.pack", tc)
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
		tested := Script("")

		a.Error(yaml.Unmarshal([]byte("hello world"), &tested))
		a.Empty(tested)
	}

	// 1. Test normal body loading
	// 1-1. empty yaml
	{
		tested := Script("")

		a.NoError(yaml.Unmarshal([]byte(""), &tested))
		a.Empty(tested)
	}

	// 1-2. filled body
	{
		tested := Script("")
		yamlBody, _ := yaml.Marshal(exampleStruct)

		a.NoError(yaml.Unmarshal(yamlBody, &tested))
		a.Equal(expectedBody+"\n", tested.String())
	}

	// 2. Test file loading
	// 2-1. invalid path loading
	{
		tested := Script("")
		exampleStruct.Source = "internal.pack/dir-file"
		yamlBody, _ := yaml.Marshal(exampleStruct)

		a.Error(yaml.Unmarshal(yamlBody, &tested))
		a.Empty(tested)
	}

	// 2-2. valid path loading
	{
		tested := Script("")

		source := "internal.pack/example.sh"
		expectedLoaded, _ := test.LoadTestFile(source)
		expectedLoaded = append(expectedLoaded, '\n')

		exampleStruct.Source = source
		yamlBody, _ := yaml.Marshal(exampleStruct)

		a.NoError(yaml.Unmarshal(yamlBody, &tested))
		a.Equal(expectedLoaded, tested.Bytes())
	}
}

func TestScript_Append(t *testing.T) {
	a := assert.New(t)

	tcFile := path.Join("internal.pack", "example.sh")

	buffer, _ := test.LoadTestFile(tcFile)
	expectedBody := string(buffer)
	expectedBody = strings.Join(append([]string{expectedBody}, loremIpsum...), "\n") + "\n"

	tested := Script("")
	a.NoError(tested.load(tcFile))

	for _, line := range loremIpsum {
		tested.Append(line)
	}

	a.Equal(expectedBody, tested.String())
}

func TestScript_String(t *testing.T) {
	// ignore this case because Script.String is a type conversion simple getter
}

func TestScript_Bytes(t *testing.T) {
	// ignore this case because Script.Bytes is a type conversion simple getter
}
