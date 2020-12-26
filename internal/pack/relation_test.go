package pack

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"testing"
)

var (
	expectedPkg = "pkg-packer"
	expectedVer = "1.2.3"

	validOps   = []string{"<", ">", "=", "<=", ">="}
	invalidOps = []string{"=>", "=<", ">>", "<<", "=="}
)

func Test_OpRegExp(t *testing.T) {
	a := assert.New(t)

	{
		tcString := fmt.Sprintf("%s %s %s", expectedPkg, "", expectedVer)
		tested := opRegExp.FindStringSubmatch(tcString)
		a.Len(tested, 4)
		a.Equal(expectedPkg, tested[1])
		a.Empty(tested[2])
		a.Equal(expectedVer, tested[3])
	}

	for _, tc := range validOps {
		tcString := fmt.Sprintf("%s %s %s", expectedPkg, tc, expectedVer)
		tested := opRegExp.FindStringSubmatch(tcString)
		a.Len(tested, 4)
		a.Equal(expectedPkg, tested[1])
		a.Contains(opMaps, tested[2])
		a.Equal(expectedVer, tested[3])
	}

	for _, tc := range invalidOps {
		tcString := fmt.Sprintf("%s %s %s", expectedPkg, tc, expectedVer)
		tested := opRegExp.FindStringSubmatch(tcString)
		a.Len(tested, 4)
		a.Equal(expectedPkg, tested[1])
		a.NotContains(opMaps, tested[2])
		a.Equal(expectedVer, tested[3])
	}
}

func TestRelation_parse(t *testing.T) {
	a := assert.New(t)

	funcTestInvalid := func(r string) {
		tested := Relation{}
		a.Error(tested.parse(r))
		a.Empty(tested.name)
		a.Empty(tested.op)
		a.Empty(tested.ver)
	}

	// 1. check package name only valid case
	{
		tested := Relation{}
		a.NoError(tested.parse(expectedPkg))
		a.Equal(expectedPkg, tested.name)
		a.Empty(string(opAny))
		a.Empty(tested.ver)
	}

	// 2. check empty operator error
	funcTestInvalid(fmt.Sprintf("%s %s", expectedPkg, expectedVer))

	// 3. check invalid operator error
	for _, tc := range invalidOps {
		funcTestInvalid(fmt.Sprintf("%s %s %s", expectedPkg, tc, expectedVer))
	}

	for _, tc := range validOps {
		// 1. test valid case
		{
			tested := Relation{}
			tcString := fmt.Sprintf("%s %s %s", expectedPkg, tc, expectedVer)
			a.NoError(tested.parse(tcString))
			a.Equal(expectedPkg, tested.name)
			a.Contains(opMaps, string(tested.op))
			a.Equal(expectedVer, tested.ver)
		}

		// 2. check empty package name error
		funcTestInvalid(fmt.Sprintf("%s %s", tc, expectedVer))

		// 3. check empty version error
		funcTestInvalid(fmt.Sprintf("%s %s", expectedPkg, tc))
	}
}

func TestRelation_UnmarshalYAML(t *testing.T) {
	a := assert.New(t)

	yamlForm := "value: %s%s%s"
	tested := struct{ Value Relation }{}

	{
		a.NoError(yaml.Unmarshal([]byte(fmt.Sprintf(yamlForm, expectedPkg, "", "")), &tested))
		a.Equal(expectedPkg, tested.Value.name)
		a.Empty(tested.Value.ver)
		a.Empty(tested.Value.op)
	}

	// format
	for _, tc := range invalidOps {
		a.Error(yaml.Unmarshal([]byte(fmt.Sprintf(yamlForm, expectedPkg, tc, expectedVer)), &tested))
	}

	// format
	for _, tc := range validOps {
		a.NoError(yaml.Unmarshal([]byte(fmt.Sprintf(yamlForm, expectedPkg, tc, expectedVer)), &tested))
		a.Equal(expectedPkg, tested.Value.name)
		a.Equal(expectedVer, tested.Value.ver)
		a.Equal(tc, string(tested.Value.op))
	}
}

func TestRelation_RpmFormat(t *testing.T) {
	a := assert.New(t)
	tested := Relation{}

	_ = tested.parse(expectedPkg)
	a.Equal(expectedPkg, tested.RpmFormat())

	for _, tc := range validOps {
		expectedDep := fmt.Sprintf("%s%s%s", expectedPkg, tc, expectedVer)
		a.NoError(tested.parse(expectedDep))
		a.Equal(expectedDep, tested.RpmFormat())
	}
}

func TestRelation_DebFormat(t *testing.T) {
	a := assert.New(t)
	tested := Relation{}

	_ = tested.parse(expectedPkg)
	a.Equal(expectedPkg, tested.DebFormat())

	for _, tc := range validOps {
		a.NoError(tested.parse(fmt.Sprintf("%s %s %s", expectedPkg, tc, expectedVer)))

		if tc == "<" || tc == ">" {
			tc += tc
		}

		a.Equal(fmt.Sprintf("%s (%s %s)", expectedPkg, tc, expectedVer), tested.DebFormat())
	}
}

func TestNewRelation(t *testing.T) {
	a := assert.New(t)

	for _, tc := range invalidOps {
		tested, err := NewRelation(fmt.Sprintf("%s %s %s", expectedPkg, tc, expectedVer))
		a.Nil(tested)
		a.Error(err)
	}

	for _, tc := range validOps {
		tested, err := NewRelation(fmt.Sprintf("%s %s %s", expectedPkg, tc, expectedVer))
		a.NotNil(tested)
		a.NoError(err)
		a.Equal(expectedPkg, tested.name)
		a.Equal(expectedVer, tested.ver)
		a.Equal(tc, string(tested.op))
	}
}
