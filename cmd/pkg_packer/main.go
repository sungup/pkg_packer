package main

import (
	"bufio"
	"flag"
	"github.com/sungup/pkg_packer/internal/builder"
	"github.com/sungup/pkg_packer/internal/info"
	"log"
	"os"
	"path"
)

type PkgPackerArgs struct {
	srcDir   string
	yamlPath string
	yamlHome string
}

func argParse() PkgPackerArgs {
	args := PkgPackerArgs{}

	flag.StringVar(&args.srcDir, "source", ".", "package source directory")
	flag.StringVar(&args.yamlPath, "yaml", "recipe.yml", "packaging recipe file")

	flag.Parse()

	args.yamlHome = path.Dir(args.yamlPath)

	return args
}

func main() {
	args := argParse()

	// 1. load yaml file
	pkgInfo, err := info.LoadPkgInfo(args.yamlPath, args.srcDir)

	if err != nil {
		log.Fatal(err)
	}

	pkgBuilder := builder.NewRPMBuilder(pkgInfo)

	// 2. open rpm file to store
	var rpmFile *os.File
	if rpmPath, err := pkgBuilder.Filename(); err == nil {
		rpmPath = path.Join(args.yamlHome, rpmPath)

		if rpmFile, err = os.Create(rpmPath); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal(err)
	}
	defer func() { _ = rpmFile.Close() }()

	// 3. write and flush
	rpmWriter := bufio.NewWriter(rpmFile)

	if err := pkgBuilder.Build(rpmWriter); err != nil {
		log.Fatal(err)
	}

	if err := rpmWriter.Flush(); err != nil {
		log.Fatal(err)
	}
}
