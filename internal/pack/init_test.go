package pack

import "github.com/sungup/pkg_packer/test"

var (
	loremIpsum = []string{
		"echo Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor",
		"echo incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud",
		"echo exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute",
		"echo irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla",
		"echo pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia",
		"echo deserunt mollit anim id est laborum.",
	}
)

func init() {
	UpdateSourceDir(test.GetTestFileHome())
}
