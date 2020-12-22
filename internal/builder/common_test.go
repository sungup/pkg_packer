package builder

import (
	"github.com/sungup/pkg_packer/internal/pack"
	"github.com/sungup/pkg_packer/test"
)

func loadYAMLExpectAndPackage() (*test.YAMLTestData, *pack.Package) {
	expect, _ := test.LoadTestYAMLData(test.ExpectFile1)
	testPath, _ := test.GetTestFilePath(test.ExpectFile1)
	tested, _ := pack.LoadPkgInfo(testPath, test.GetTestFileHome())

	return expect, tested
}
