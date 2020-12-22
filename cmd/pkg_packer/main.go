package main

import (
	"bufio"
	"flag"
	"github.com/sungup/pkg_packer/internal/builder"
	"github.com/sungup/pkg_packer/internal/pack"
	"log"
	"os"
	"path"
)

type PkgPackerArgs struct {
	srcDir   string
	yamlPath string
	yamlHome string
	pkgType  string
}

func argParse() PkgPackerArgs {
	args := PkgPackerArgs{}

	flag.StringVar(&args.srcDir, "source", ".", "package source directory")
	flag.StringVar(&args.yamlPath, "yaml", "recipe.yml", "packaging recipe file")
	flag.StringVar(&args.pkgType, "pkg-type", "all", "package type [all|rpm|deb]")

	flag.Parse()

	args.yamlHome = path.Dir(args.yamlPath)

	return args
}

type builderFunc func(*pack.Package) builder.PackageBuilder

func (args *PkgPackerArgs) builders() []builderFunc {
	switch args.pkgType {
	case "rpm":
		return []builderFunc{builder.NewRPMBuilder}
	case "deb":
		return []builderFunc{builder.NewDEBBuilder}
	case "all":
		return []builderFunc{builder.NewRPMBuilder, builder.NewDEBBuilder}
	}

	log.Fatalln("unsupported builder type")

	return nil
}

func buildPackage(builder builder.PackageBuilder, directory string) {
	var pkgFile *os.File
	if pkgPath, err := builder.Filename(); err == nil {
		pkgPath = path.Join(directory, pkgPath)

		if pkgFile, err = os.Create(pkgPath); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal(err)
	}
	defer func() { _ = pkgFile.Close() }()

	pkgWriter := bufio.NewWriter(pkgFile)

	if err := builder.Build(pkgWriter); err != nil {
		log.Fatal(err)
	}

	if err := pkgWriter.Flush(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	args := argParse()

	// 1. load yaml file
	pkgInfo, err := pack.LoadPkgInfo(args.yamlPath, args.srcDir)

	if err != nil {
		log.Fatal(err)
	}

	for _, builderNew := range args.builders() {
		buildPackage(builderNew(pkgInfo), args.yamlHome)
	}
}
