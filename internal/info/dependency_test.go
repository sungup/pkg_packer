package info

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestDependency_UnmarshalYAML(t *testing.T) {
	const (
		format = "dep: %s%s%s"
		expectedPkg = "hello"
	)

	a := assert.New(t)

	var testCase = []string {
		"",
		">",
		">=",
		"=",
		"<=",
		"<",
	}

	for _, expectedOp := range testCase {
		expectedVer := ""
		if expectedOp != "" {
			expectedVer = "1.0"
		}

		input := fmt.Sprintf(format, expectedPkg, expectedOp, expectedVer)
		tested := struct {
			Dep Dependency `yaml:"dep"`
		}{}

		a.NoError(yaml.Unmarshal([]byte(input), &tested))
		a.Equal(expectedPkg, tested.Dep.name)
		a.Equal(expectedOp,  tested.Dep.operator)
		a.Equal(expectedVer, tested.Dep.ver)
	}
}