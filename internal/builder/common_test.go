package builder

import (
	"github.com/sungup/pkg_packer/internal/pkg"
	"github.com/sungup/pkg_packer/test"
)

func loadYAMLExpectAndPackage() (*test.YAMLTestData, *pkg.Package) {
	expect, _ := test.LoadTestYAMLData(test.ExpectFile1)
	testPath, _ := test.GetTestFilePath(test.ExpectFile1)
	tested, _ := pkg.LoadPkgInfo(testPath, test.GetTestFileHome())

	return expect, tested
}
