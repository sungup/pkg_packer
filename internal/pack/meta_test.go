package pack

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestMeta_checkMeta(t *testing.T) {
	a := assert.New(t)

	tested := Meta{
		Name:       "pkg_packer",
		Version:    "0.2.0",
		Release:    "universal",
		Arch:       "x86_64",
		Summary:    "test summary",
		Maintainer: "Sungup Moon <sungup@me.com>",
	}

	a.NoError(tested.checkMeta())

	elements := reflect.ValueOf(&tested).Elem()

	for _, name := range mustFill {
		field := elements.FieldByName(name)
		orig := field.String()
		field.SetString("")

		a.EqualError(tested.checkMeta(), fmt.Sprintf(errMetaDataNotFilled, name))

		field.SetString(orig)
	}
}

func TestMeta_UTCBuildTime(t *testing.T) {
	a := assert.New(t)

	loc, _ := time.LoadLocation("Asia/Seoul")
	expectedTime := time.Now().In(loc)

	meta := Meta{BuildTime: expectedTime}

	a.NotEqual(expectedTime, meta.UTCBuildTime())
	a.Equal(expectedTime.UnixNano(), meta.UTCBuildTime().UnixNano())
}

func TestMeta_UnmarshalYAML(t *testing.T) {
	a := assert.New(t)

	updateBuildTimeString := "2020-12-25T22:00:00+09:00"
	updatedBuildTime, _ := time.Parse(time.RFC3339, updateBuildTimeString)
	defaultBuildTime := time.Now().Add(-1 * time.Hour)
	defaultDesc := "default description"
	updatedDesc := "updated description"
	tested := Meta{Description: defaultDesc, BuildTime: defaultBuildTime}

	buffer := map[string]string{
		"name":       "pkg_packer",
		"version":    "0.2.0",
		"release":    "universal",
		"arch":       "x86_64",
		"summary":    "test summary",
		"maintainer": "Sungup Moon <sungup@me.com>",
		"desc":       updatedDesc,
		"build_time": updateBuildTimeString,
	}

	// 1. valid parsing with inserted build time
	{
		yamlData, _ := yaml.Marshal(buffer)

		a.NoError(yaml.Unmarshal(yamlData, &tested))
		a.Equal(updatedDesc, tested.Description)
		a.Equal(updatedBuildTime.UTC(), tested.UTCBuildTime())
	}

	{
		delete(buffer, "build_time")

		yamlData, _ := yaml.Marshal(buffer)

		a.NoError(yaml.Unmarshal(yamlData, &tested))
		a.Equal(updatedDesc, tested.Description)
		a.True(updatedBuildTime.Before(tested.UTCBuildTime()))

		buffer["build_time"] = updateBuildTimeString
	}

	// 2. invalid parsing
	for _, name := range mustFill {
		// restore default value
		tested.Description = defaultDesc
		tested.BuildTime = defaultBuildTime

		expectedErr := fmt.Sprintf(errMetaDataNotFilled, name)

		// change name to lower case and remove that field for exception field test
		name = strings.ToLower(name)
		orig := buffer[name]
		delete(buffer, name)
		yamlData, _ := yaml.Marshal(buffer)

		a.EqualError(yaml.Unmarshal(yamlData, &tested), expectedErr)
		a.Equal(defaultDesc, tested.Description)
		a.Equal(defaultBuildTime.UTC(), tested.UTCBuildTime())

		buffer[name] = orig
	}
}
