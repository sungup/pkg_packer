package main

import (
	"bufio"
	"github.com/sungup/pkg_packer/internal/builder"
	"github.com/sungup/pkg_packer/internal/info"
	"log"
	"os"
	"path"
	"time"
)

func writeToDir(dir string, pkgBuilder builder.PackageBuilder) {
	// 1. open file
	var pkgFile *os.File
	if pkgPath, err := pkgBuilder.Filename(); err == nil {
		pkgPath = path.Join(dir, pkgPath)

		if pkgFile, err = os.Create(pkgPath); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal(err)
	}
	defer func() { _ = pkgFile.Close() }()

	// 2. write and flush
	pkgWriter := bufio.NewWriter(pkgFile)

	if err := pkgBuilder.Build(pkgWriter); err != nil {
		log.Fatal(err)
	}

	if err := pkgWriter.Flush(); err != nil {
		log.Fatal(err)
	}
}

func apiSample() {
	pkgInfo := info.NewPackage(info.PackageMeta{
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

	pkgInfo.AddDirectory(info.PackageDir{
		Dest:  "/var/lib/pkg-packer-api-sample",
		Mode:  0755,
		Owner: "root",
		Group: "root",
	})

	_ = pkgInfo.AddFile("config", info.PackageFile{
		Dest:  "/var/lib/pkg-packer-api-sample/sample1.ini",
		Body:  `[test]\nvalue="Hello api-sample log file!"`,
		Mode:  0644,
		Owner: "root",
		Group: "root",
		MTime: time.Now(),
	})

	_ = pkgInfo.AddFile("generic", info.PackageFile{
		Dest:  "/var/lib/pkg-packer-api-sample/sample2.sh",
		Body:  "#!/bin/bash\necho Hello api-sample log file!\n",
		Mode:  0755,
		Owner: "root",
		Group: "root",
		MTime: time.Now(),
	})

	rpmBuilder := builder.NewRPMBuilder(pkgInfo)
	writeToDir("temp", rpmBuilder)

	debBuilder := builder.NewDEBBuilder(pkgInfo)
	writeToDir("temp", debBuilder)
}

func yamlSample(yamlPath string) {
	pkgInfo, err := info.LoadPkgInfo(yamlPath, ".")
	if err != nil {
		log.Fatal(err)
	}

	rpmBuilder := builder.NewRPMBuilder(pkgInfo)
	writeToDir("temp", rpmBuilder)

	debBuilder := builder.NewDEBBuilder(pkgInfo)
	writeToDir("temp", debBuilder)
}

func main() {
	// using api example
	apiSample()

	// using yaml example
	yamlSample("test/test.yml")
}
