package main

import (
	"bufio"
	"github.com/google/rpmpack"
	"github.com/sungup/pkg_packer/internal/builder"
	"github.com/sungup/pkg_packer/internal/pkg"
	"log"
	"os"
	"path"
)

func main() {
	meta := pkg.PackageMeta{
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
	}

	rBuild := builder.RPMBuilder{}

	r, err := rpmpack.NewRPM(rBuild.Metadata(&meta))

	if err != nil {
		panic(err)
	}

	r.AddFile(rBuild.File("not_use", &pkg.PackageFile{
		Dest:  "/tmp/rpmpack_test.log",
		Body:  "Hello World\n",
		Mode:  0644,
		Owner: "root",
		Group: "root",
	}))

	r.AddFile(rBuild.File("generic", &pkg.PackageFile{
		Dest:  "/tmp/test.yml",
		Src:   "test/test.yml",
		Owner: "root",
		Group: "root",
	}))

	filename, _ := rBuild.Filename(&meta)

	fRpm, _ := os.Create(path.Join("temp", filename))
	defer func() { _ = fRpm.Close() }()

	rpmWriter := bufio.NewWriter(fRpm)

	if err := r.Write(rpmWriter); err != nil {
		log.Fatalf("write failed: %v", err)
	}

	_ = rpmWriter.Flush()
}
