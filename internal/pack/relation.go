package pack

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"regexp"
	"strings"
)

const (
	opAny = operator("")
	opLT  = operator("<")
	opLE  = operator("<=")
	opEQ  = operator("=")
	opGT  = operator(">")
	opGE  = operator(">=")
)

var (
	opRegExp = regexp.MustCompile(`([^=<>\s]*)\s*([<>=]*)\s*(.*)?`)

	opMaps = map[string]operator{
		"":   opAny,
		"<":  opLT,
		">":  opGT,
		"=":  opEQ,
		"<=": opLE,
		">=": opGE,
	}
)

type operator string

type Relation struct {
	name string
	ver  string
	op   operator
}

func (r *Relation) parse(rel string) error {
	var (
		fields = opRegExp.FindStringSubmatch(rel)

		op operator
		ok bool
	)

	if fields[1] == "" {
		return fmt.Errorf("relation package name should not be empty: %s", rel)
	} else if op, ok = opMaps[fields[2]]; !ok {
		return fmt.Errorf("unexpected operator value: %s", fields[2])
	} else if op == opAny && fields[3] != "" {
		return fmt.Errorf("version should be empty for the empty operator: %s", rel)
	} else if op != opAny && fields[3] == "" {
		return fmt.Errorf("version should be exist if operator specified: %s", op)
	}

	r.name, r.op, r.ver = fields[1], op, fields[3]

	return nil
}

func (r *Relation) UnmarshalYAML(value *yaml.Node) error {
	var verString string

	if err := value.Decode(&verString); err != nil {
		return fmt.Errorf("unexpected version type string")
	}

	return r.parse(verString)
}

func (r *Relation) RpmFormat() string {
	return fmt.Sprintf("%s%s%s", r.name, r.op, r.ver)
}

func (r *Relation) DebFormat() string {
	var ver string

	if r.op != "" {
		// debian deprecated single < and >. they recommend using << and >>.
		if r.op == opLT || r.op == opGT {
			ver = fmt.Sprintf(" (%s %s)", strings.Repeat(string(r.op), 2), r.ver)
		} else {
			ver = fmt.Sprintf(" (%s %s)", r.op, r.ver)
		}
	}

	return fmt.Sprintf("%s%s", r.name, ver)
}

func NewRelation(info string) (*Relation, error) {
	var r Relation

	if err := r.parse(info); err != nil {
		return nil, err
	} else {
		return &r, nil
	}
}
