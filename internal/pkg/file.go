package pkg

import "time"

type PackageFile struct {
	Dest  string `yaml:"dest,omitempty"`
	Src   string `yaml:"source"`
	Body  []byte `yaml:"body"`
	Mode  uint   `yaml:"mode"`
	Owner string `yaml:"owner"`
	Group string `yaml:"group"`
	MTime time.Time
}

type PackageDir struct {
	Dest  string `yaml:"dest,omitempty"`
	Mode  uint   `yaml:"mode"`
	Owner string `yaml:"owner"`
	Group string `yaml:"owner"`
}
