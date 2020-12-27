package main

import (
	"bufio"
	"github.com/sungup/pkg_packer/internal/builder"
	"github.com/sungup/pkg_packer/internal/pack"
	"log"
	"os"
	"path"
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
	pkgInfo := pack.NewPackage(pack.Meta{
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

	dir, _ := pack.NewDirectory("/var/lib/pkg-packer-api-sample", "root", "root", 0755)
	pkgInfo.AddDirectory(dir)

	file, _ := pack.NewFile(
		`[test]\nvalue="Hello api-sample log file!"`,
		"/var/lib/pkg-packer-api-sample/sample1.ini",
		"root", "root", 0644,
	)
	_ = pkgInfo.AddFile("config", file)

	file, _ = pack.NewFile(
		"#!/bin/bash\necho Hello api-sample log file!\n",
		"/var/lib/pkg-packer-api-sample/sample2.sh",
		"root", "root", 0755,
	)
	_ = pkgInfo.AddFile("generic", file)

	rpmBuilder := builder.NewRPMBuilder(pkgInfo)
	writeToDir("temp", rpmBuilder)

	debBuilder := builder.NewDEBBuilder(pkgInfo)
	writeToDir("temp", debBuilder)
}

func yamlSample(yamlPath string) {
	pkgInfo, err := pack.LoadPkgInfo(yamlPath)
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
