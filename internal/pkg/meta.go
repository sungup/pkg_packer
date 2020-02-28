package pkg

type PackageMeta struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
	Release string `yaml:"release"`
	Arch    string `yaml:"arch"`

	Summary     string `yaml:"summary"`
	Description string `yaml:"desc"`

	OS     string `yaml:"os,omitempty"`
	Vendor string `yaml:"vendor,omitempty"`

	URL     string `yaml:"url,omitempty"`
	License string `yaml:"license,omitempty"`

	// TODO add requires
}
