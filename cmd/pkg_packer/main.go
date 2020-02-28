package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/google/rpmpack"
	"log"
	"os"
	"time"
)

func toRPMFileName(meta rpmpack.RPMMetaData) (string, error) {
	if meta.Name == "" {
		return "", errors.New("undefined package name")
	} else if meta.Version == "" {
		return "", errors.New("undefined package version")
	}

	buffer := new(bytes.Buffer)

	buffer.WriteString(fmt.Sprintf("%s-%s", meta.Name, meta.Version))

	if meta.Release != "" {
		buffer.WriteString(fmt.Sprintf(".%s", meta.Release))
	}

	if meta.Arch != "" {
		buffer.WriteString(fmt.Sprintf(".%s", meta.Arch))
	} else {
		buffer.WriteString(".noarch")
	}

	buffer.WriteString(".rpm")

	return buffer.String(), nil
}

func main() {
	r, err := rpmpack.NewRPM(
		rpmpack.RPMMetaData{
			Name:        "rpmpack-test",
			Summary:     "",
			Description: "",
			Version:     "0.0.1-1",
			Release:     "el7",
			Arch:        "x86_64",
			OS:          "linux",
			Vendor:      "",
			URL:         "",
			Packager:    "",
			Group:       "",
			Licence:     "",
			BuildHost:   "",
			Compressor:  "",
			Epoch:       0,
			BuildTime:   time.Time{},
			Provides:    nil,
			Obsoletes:   nil,
			Suggests:    nil,
			Recommends:  nil,
			Requires:    nil,
			Conflicts:   nil,
		},
	)

	if err != nil {
		panic(err)
	}

	r.AddFile(rpmpack.RPMFile{
		Name:  "/tmp/rpmpack_test.log",
		Body:  []byte("Hello World\n"),
		Mode:  0644,
		Owner: "root",
		Group: "root",
		Type:  rpmpack.GenericFile,
	})

	filename, _ := toRPMFileName(r.RPMMetaData)

	fRpm, _ := os.Create(filename)
	defer func() { _ = fRpm.Close() }()

	rpmWriter := bufio.NewWriter(fRpm)

	if err := r.Write(rpmWriter); err != nil {
		log.Fatalf("write failed: %v", err)
	}
}
