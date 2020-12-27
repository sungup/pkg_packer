package pack

import (
	"fmt"
	"github.com/sungup/pkg_packer/pkg/utils"
	"gopkg.in/yaml.v3"
	"reflect"
	"time"
)

var (
	mustFill = []string{"Name", "Version", "Release", "Arch", "Summary", "Maintainer"}
)

const (
	errMetaDataNotFilled  = "metadata field %s is empty"
	errMetaDecodingFailed = "loading metadata configuration failed: %w"
)

type Meta struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
	Release string `yaml:"release"`
	Arch    string `yaml:"arch"`

	Summary     string `yaml:"summary"`
	Description string `yaml:"desc"`

	OS     string `yaml:"os"`
	Vendor string `yaml:"vendor"`

	URL     string `yaml:"url"`
	License string `yaml:"license"`

	Maintainer string `yaml:"maintainer"`

	BuildTime time.Time `yaml:"build_time"`

	// TODO add requires
}

func (meta *Meta) checkMeta() error {
	elements := reflect.ValueOf(meta).Elem()
	for i := 0; i < elements.NumField(); i++ {
		name := elements.Type().Field(i).Name
		if utils.Contains(mustFill, name) && elements.Field(i).String() == "" {
			return fmt.Errorf(errMetaDataNotFilled, name)
		}
	}

	return nil
}

func (meta *Meta) UTCBuildTime() time.Time {
	return meta.BuildTime.UTC()
}

func (meta *Meta) UnmarshalYAML(value *yaml.Node) error {
	// set default data
	type tT Meta
	tV := tT{BuildTime: time.Now()}

	if err := value.Decode(&tV); err != nil {
		return fmt.Errorf(errMetaDecodingFailed, err)
	}

	parsed := Meta(tV)
	if err := parsed.checkMeta(); err != nil {
		return err
	}

	*meta = parsed

	return nil
}
