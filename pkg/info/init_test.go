package info

import (
	"fmt"
	"github.com/sungup/pkg_packer/test"
	"gopkg.in/yaml.v3"
	"path"
	"time"
)

var (
	loremIpsum = []string{
		"echo Lorem ipsum dolor sit amet, consectetur adipiscing elit,",
		"echo sed do eiusmod tempor incididunt ut labore et dolore magna",
		"echo aliqua. Ut enim ad minim veniam, quis nostrud exercitation",
		"echo ullamco laboris nisi ut aliquip ex ea commodo consequat.",
		"echo Duis aute irure dolor in reprehenderit in voluptate velit",
		"echo esse cillum dolore eu fugiat nulla pariatur. Excepteur sint",
		"echo occaecat cupidatat non proident, sunt in culpa qui officia",
		"echo deserunt mollit anim id est laborum.",
	}

	expectedPkg = "pkg-packer"
	expectedVer = "1.2.3"

	referenceValidOps   = []string{"<", ">", "=", "<=", ">="}
	referenceInvalidOps = []string{"=>", "=<", ">>", "<<", "=="}

	referenceFileList []map[string]interface{}

	referencePackage = Package{
		Meta: Meta{
			Name:        "test_package",
			Version:     "0.0.1b",
			Release:     "all",
			Arch:        "x86_64",
			Summary:     "testcase package",
			Description: "testcase description info",
			OS:          "ubuntu",
			Vendor:      "sungup",
			URL:         "https://github.com/sungup/pkg_packer",
			License:     "Apache",
			Maintainer:  "Sungup Moon <sungup@me.com>",
			BuildTime:   time.Now(),
		},
		Dirs: []*Directory{
			{Dest: "/var/lib/pkg_packer/dir01", Owner: "root", Group: "root", Mode: 0755},
			{Dest: "/var/lib/pkg_packer/dir02", Owner: "root", Group: "root", Mode: 0755},
			{Dest: "/var/lib/pkg_packer/dir03", Owner: "root", Group: "root", Mode: 0755},
		},
		Files:  nil,
		PreIn:  script(loremIpsum[0] + "\n" + loremIpsum[1]),
		PostIn: script(loremIpsum[2] + "\n" + loremIpsum[3]),
		PreUn:  script(loremIpsum[4] + "\n" + loremIpsum[5]),
		PostUn: script(loremIpsum[6] + "\n" + loremIpsum[7]),
		Dependencies: []*Relation{
			{name: "test_sample1", op: opEQ, ver: "0.0.1b"},
			{name: "test_sample2", op: opLE, ver: "0.0.1b"},
			{name: "test_sample3", op: opGT, ver: "0.0.1b"},
		},
	}
)

func init() {
	UpdateSourceDir(test.GetTestFileHome())

	// to clear time.Now's location mismatch
	body, _ := yaml.Marshal(&referencePackage)
	_ = yaml.Unmarshal(body, &referencePackage)

	for lineNo, line := range loremIpsum {
		referenceFileList = append(referenceFileList, map[string]interface{}{
			"dest":  path.Join("/var/lib/pkg_packer", fmt.Sprintf("file%d", lineNo)),
			"owner": "root",
			"group": "root",
			"mode":  0744,
			"mtime": time.Now(),
			"body":  line,
		})
	}
}
