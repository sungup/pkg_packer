package main

import (
	"bufio"
	"github.com/sungup/pkg_packer/internal/builder"
	"github.com/sungup/pkg_packer/internal/pkg"
	"log"
	"os"
	"path"
)

func main() {
	pkgInfo := pkg.NewPackage(pkg.PackageMeta{
		Name:        "rpmpack-test",
		Version:     "0.0.1-1",
		Release:     "el7",
		Arch:        "x86_64",
		Summary:     "",
		Description: "",
		OS:          "linux",
		Vendor:      "",
		URL:         "",
		License:     "",
	})

	pkgInfo.Files["not_use"] = append(
		pkgInfo.Files["not_use"],
		pkg.PackageFile{
			Dest:  "/tmp/rpmpack_test.log",
			Body:  "Hello World\n",
			Mode:  0644,
			Owner: "root",
			Group: "root",
		},
	)

	pkgInfo.Files["generic"] = append(
		pkgInfo.Files["not_use"],
		pkg.PackageFile{
			Dest:  "/tmp/test.yml",
			Src:   "test/test.yml",
			Owner: "root",
			Group: "root",
		},
	)

	rBuild := builder.NewRPMBuilder(pkgInfo)

	var rpmPath string
	var err error
	if rpmPath, err = rBuild.Filename(); err != nil {
		log.Fatal(err)
	}

	rpmPath = path.Join("temp", rpmPath)

	// 4. create file
	var rpmFile *os.File
	if rpmFile, err = os.Create(rpmPath); err != nil {
		log.Fatal(err)
	}
	defer func() { _ = rpmFile.Close() }()

	// 5. Write and flush file
	rpmWriter := bufio.NewWriter(rpmFile)

	if err = rBuild.Build(rpmWriter); err != nil {
		log.Fatal(err)
	}

	if err = rpmWriter.Flush(); err != nil {
		log.Fatal(err)
	}
}
