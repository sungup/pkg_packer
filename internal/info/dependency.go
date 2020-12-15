package info

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"regexp"
)

type Dependency struct {
	name     string
	ver      string
	operator string
}

func (dep *Dependency) UnmarshalYAML(value *yaml.Node) error {
	var verString string

	if err := value.Decode(&verString); err != nil {
		return fmt.Errorf("unexpected versiont type string")
	}

	re, err := regexp.Compile("[<=>]+")
	if err !=nil {
		return fmt.Errorf("unexpected regex format: %v", err)
	}

	pos := re.FindStringIndex(verString)
	if pos == nil {
		dep.name     = verString
		dep.operator = ""
		dep.ver      = ""
	} else {
		dep.name     = verString[:pos[0]]
		dep.operator = verString[pos[0]:pos[1]]
		dep.ver      = verString[pos[1]:]
	}

	return nil
}

func (dep *Dependency) RpmFormat() string {
	return fmt.Sprintf("%s%s%s", dep.name, dep.operator, dep.ver)
}

func (dep *Dependency) DebFormat() string {
	operator := dep.operator

	// debian deprecated single < and >. they recommend using << and >>.
	if dep.operator == "<" || dep.operator == ">" {
		operator += dep.operator
	}

	return fmt.Sprintf("%s (%s %s)", dep.name, dep.operator, dep.ver)
}
