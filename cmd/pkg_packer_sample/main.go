package main

import (
	"bufio"
	"github.com/sungup/pkg_packer/internal/builder"
	"github.com/sungup/pkg_packer/internal/pkg"
	"log"
	"os"
	"path"
	"time"
)

func writeToDir(dir string, pkgBuilder *builder.RPMBuilder) {
	// 1. open file
	var rpmFile *os.File
	if rpmPath, err := pkgBuilder.Filename(); err == nil {
		rpmPath = path.Join(dir, rpmPath)

		if rpmFile, err = os.Create(rpmPath); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal(err)
	}
	defer func() { _ = rpmFile.Close() }()

	// 2. write and flush
	rpmWriter := bufio.NewWriter(rpmFile)

	if err := pkgBuilder.Build(rpmWriter); err != nil {
		log.Fatal(err)
	}

	if err := rpmWriter.Flush(); err != nil {
		log.Fatal(err)
	}
}

func apiSample() {
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
	}, ".")

	pkgInfo.AddDirectory(pkg.PackageDir{
		Dest:  "/var/lib/pkg-packer-api-sample",
		Mode:  0755,
		Owner: "root",
		Group: "root",
	})

	_ = pkgInfo.AddFile("config", pkg.PackageFile{
		Dest:  "/var/lib/pkg-packer-api-sample/sample1.ini",
		Body:  `[test]\nvalue="Hello api-sample log file!"`,
		Mode:  0644,
		Owner: "root",
		Group: "root",
		MTime: time.Now(),
	})

	_ = pkgInfo.AddFile("generic", pkg.PackageFile{
		Dest:  "/var/lib/pkg-packer-api-sample/sample2.sh",
		Body:  "#!/bin/bash\necho Hello api-sample log file!\n",
		Mode:  0755,
		Owner: "root",
		Group: "root",
		MTime: time.Now(),
	})

	pkgBuilder := builder.NewRPMBuilder(pkgInfo)

	writeToDir("temp", pkgBuilder)
}

func yamlSample(yamlPath string) {
	pkgInfo, err := pkg.LoadPkgInfo(yamlPath, ".")
	if err != nil {
		log.Fatal(err)
	}

	pkgBuilder := builder.NewRPMBuilder(pkgInfo)

	writeToDir("temp", pkgBuilder)
}

func main() {
	// using api example
	apiSample()

	// using yaml example
	yamlSample("test/test.yml")
}
