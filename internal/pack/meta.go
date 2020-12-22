package pack

import "time"

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

	Maintainer string `yaml:"maintainer,omitempty"`

	buildTime time.Time
	// TODO add requires
}

func (meta *PackageMeta) BuildTime() time.Time {
	return meta.buildTime.UTC()
}

func (meta *PackageMeta) UpdateBuildTime(buildTime time.Time) {
	meta.buildTime = buildTime
}
