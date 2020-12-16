package debpack

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/blakesmith/ar"
)

type tarGzWriter struct {
	buffer *bytes.Buffer
	gz     *gzip.Writer
	tar    *tar.Writer
}

func (tgz *tarGzWriter) WriteHeader(hdr *tar.Header) error {
	return tgz.tar.WriteHeader(hdr)
}

func (tgz *tarGzWriter) Write(hdr *tar.Header, b []byte) (int, error) {
	if err := tgz.WriteHeader(hdr); err != nil {
		return -1, fmt.Errorf("write header failed: %v", err)
	}

	if sz, err := tgz.tar.Write(b); err != nil {
		return -1, fmt.Errorf("write body failed: %v", err)
	} else {
		return sz, nil
	}
}

func (tgz *tarGzWriter) Close() error {
	if err := tgz.tar.Close(); err != nil {
		return fmt.Errorf("closing tar failed: %v", err)
	}

	if err := tgz.gz.Close(); err != nil {
		return fmt.Errorf("closing gz failed: %v", err)
	}

	return nil
}

func (tgz *tarGzWriter) Bytes() []byte {
	return tgz.buffer.Bytes()
}

func newTarGz() *tarGzWriter {
	buffer := new(bytes.Buffer)
	gzWriter := gzip.NewWriter(buffer)
	tarWriter := tar.NewWriter(gzWriter)

	return &tarGzWriter{
		buffer: buffer,
		gz:     gzWriter,
		tar:    tarWriter,
	}
}

type md5Writer struct {
	buffer *bytes.Buffer
	temp   []byte
}

func (md *md5Writer) Record(b []byte, name string) error {
	digest := md5.New()

	if _, err := digest.Write(b); err != nil {
		return fmt.Errorf("writing digest failed: %v", err)
	}

	if _, err := fmt.Fprintf(md.buffer, "%x %s\n", digest.Sum(md.temp), name); err != nil {
		return fmt.Errorf("storing md5 digest failed: %v", err)
	}

	return nil
}

func (md *md5Writer) MD5Sums() []byte {
	return md.buffer.Bytes()
}

func newMD5() *md5Writer {
	return &md5Writer{
		buffer: new(bytes.Buffer),
		temp:   make([]byte, 0, md5.Size),
	}
}

type DEB struct {
	DEBMetaData

	files []DEBFile

	preIn  string
	postIn string
	preUn  string
	postUn string

	conffiles *bytes.Buffer
}

func (deb *DEB) internalFilePath(file *DEBFile) (string, error) {
	if !strings.HasPrefix(file.Name, "/") {
		return "", fmt.Errorf("input file path is not an absolute path: %s", file.Name)
	}

	return "." + file.Name, nil
}

func (deb *DEB) compressFile(file *DEBFile, data *tarGzWriter, md5sum *md5Writer) error {
	installPath, err := deb.internalFilePath(file)
	if err != nil {
		return err
	}

	h := tar.Header{
		Name:     installPath,
		Size:     int64(len(file.Body)),
		Mode:     int64(file.Mode),
		ModTime:  file.MTime,
		Typeflag: tar.TypeReg,
	}

	if _, err := data.Write(&h, file.Body); err != nil {
		return fmt.Errorf("compressing file failed: %v", err)
	}

	if err := md5sum.Record(file.Body, installPath[2:]); err != nil {
		return fmt.Errorf("generating md5 information for %s failed: %v", installPath, err)
	}

	if file.isConfig() {
		if _, err := fmt.Fprintln(deb.conffiles, file.Name); err != nil {
			return fmt.Errorf("generating conffiles information for %s failed: %v", installPath, err)
		}
	}

	return nil
}

func (deb *DEB) compressDir(file *DEBFile, data *tarGzWriter) error {
	installPath, err := deb.internalFilePath(file)
	if err != nil {
		return err
	}

	if !strings.HasSuffix(installPath, "/") {
		installPath += "/"
	}

	h := tar.Header{
		Name:     installPath,
		Mode:     int64(file.Mode),
		ModTime:  file.MTime,
		Typeflag: tar.TypeDir,
	}

	if err := data.WriteHeader(&h); err != nil {
		return fmt.Errorf("compressing dir failed: %v", err)
	}

	return nil
}

func (deb *DEB) compressMeta(filename string, body []byte, meta *tarGzWriter) error {
	h := tar.Header{
		Name:     filename,
		Size:     int64(len(body)),
		Mode:     0644,
		ModTime:  time.Now(),
		Typeflag: tar.TypeReg,
	}

	if _, err := meta.Write(&h, body); err != nil {
		return fmt.Errorf("compressing %s file failed: %v", filename, err)
	}

	return nil
}

func (deb *DEB) compressControl(meta *tarGzWriter) error {
	return deb.compressMeta("control", deb.MakeControl(), meta)
}

func (deb *DEB) compressConfFiles(meta *tarGzWriter) error {
	return deb.compressMeta("conffiles", deb.conffiles.Bytes(), meta)
}

func (deb *DEB) compressMD5(meta *tarGzWriter, md5sum *md5Writer) error {
	return deb.compressMeta("md5sums", md5sum.MD5Sums(), meta)
}

func (deb *DEB) compressScripts(meta *tarGzWriter) error {
	if err := deb.compressMeta("preinst", []byte(deb.preIn), meta); err != nil {
		return err
	}

	if err := deb.compressMeta("postinst", []byte(deb.postIn), meta); err != nil {
		return err
	}

	if err := deb.compressMeta("prerm", []byte(deb.preUn), meta); err != nil {
		return err
	}

	if err := deb.compressMeta("postrm", []byte(deb.postUn), meta); err != nil {
		return err
	}

	return nil
}

func (deb *DEB) arCompress(writer *ar.Writer, filename string, body []byte) error {
	h := ar.Header{
		Name:    filename,
		Size:    int64(len(body)),
		Mode:    0644,
		ModTime: time.Now(),
	}

	if err := writer.WriteHeader(&h); err != nil {
		return fmt.Errorf("cannot write file header: %v", err)
	}

	_, err := writer.Write(body)

	return err
}

func (deb *DEB) Write(w io.Writer) error {
	meta := newTarGz()
	data := newTarGz()
	md5sum := newMD5()

	// 1. compress file information
	for _, file := range deb.files {
		if file.isDir() {
			if err := deb.compressDir(&file, data); err != nil {
				return err
			}
		} else {
			if err := deb.compressFile(&file, data, md5sum); err != nil {
				return err
			}
		}
	}

	if err := deb.compressControl(meta); err != nil {
		return err
	}

	if err := deb.compressMD5(meta, md5sum); err != nil {
		return err
	}

	if err := deb.compressConfFiles(meta); err != nil {
		return err
	}

	if err := deb.compressScripts(meta); err != nil {
		return err
	}

	_ = meta.Close()
	_ = data.Close()

	writer := ar.NewWriter(w)

	if err := writer.WriteGlobalHeader(); err != nil {
		return fmt.Errorf("cannot write ar header to deb file: %v", err)
	}

	if err := deb.arCompress(writer, "debian-binary", []byte("2.0\n")); err != nil {
		return fmt.Errorf("cannot writ ar header to deb file: %v", err)
	}

	if err := deb.arCompress(writer, "control.tar.gz", meta.Bytes()); err != nil {
		return fmt.Errorf("cannot add control.tar.gz to deb: %v", err)
	}

	if err := deb.arCompress(writer, "data.tar.gz", data.Bytes()); err != nil {
		return fmt.Errorf("cannot add data.tar.gz to deb: %v", err)
	}

	// implementing here
	return nil
}

func (deb *DEB) AddFile(file DEBFile) { deb.files = append(deb.files, file) }
func (deb *DEB) AddPrein(s string)    { deb.preIn = s }
func (deb *DEB) AddPostin(s string)   { deb.postIn = s }
func (deb *DEB) AddPreun(s string)    { deb.preUn = s }
func (deb *DEB) AddPostun(s string)   { deb.postUn = s }

// NewDEB is a constructor for the debian package builder.
func NewDEB(meta DEBMetaData) (*DEB, error) {

	return &DEB{
		DEBMetaData: meta,

		conffiles: bytes.NewBufferString(""),
	}, nil
}
