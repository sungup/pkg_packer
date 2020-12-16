package builder

import (
	"github.com/sungup/pkg_packer/internal/info"
	"github.com/sungup/pkg_packer/test"
)

func loadYAMLExpectAndPackage() (*test.YAMLTestData, *info.Package) {
	expect, _ := test.LoadTestYAMLData(test.ExpectFile1)
	testPath, _ := test.GetTestFilePath(test.ExpectFile1)
	tested, _ := info.LoadPkgInfo(testPath, test.GetTestFileHome())

	return expect, tested
}
