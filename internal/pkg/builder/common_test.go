package builder

import (
	"github.com/sungup/pkg_packer/internal/pkg"
	"github.com/sungup/pkg_packer/test"
)

func LoadYAMLExpectAndPackage() (*test.YAMLTestData, *pkg.Package) {
	expect, _ := test.LoadTestYAMLData(test.ExpectFile)
	testPath, _ := test.GetTestFilePath(test.ExpectFile)
	tested, _ := pkg.LoadPkgInfo(testPath)

	return expect, tested
}
