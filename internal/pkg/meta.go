package pkg

type PackageMeta struct {
	Name    string `yaml:"name,omitempty"`
	Version string `yaml:"version,omitempty"`
	Release string `yaml:"release,omitempty"`
	Arch    string `yaml:"arch,omitempty"`

	Summary     string `yaml:"summary"`
	Description string `yaml:"desc"`

	OS     string `yaml:"os,omitempty"`
	Vendor string `yaml:"vendor,omitempty"`

	URL     string `yaml:"url,omitempty"`
	License string `yaml:"license,omitempty"`

	// TODO add requires
}
