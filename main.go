package main

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/alexflint/go-arg"
)

//go:generate go run generators/generator.go

type args struct {
	Out  string `arg:"-o" help:"The out path to write the shared AWS C runtime library mock. Default is /tmp/gogreengrass" placeholder:"PATH"`
	SDKC bool   `arg:"-l" help:"Installs the c runtime shared library in /tmp/gogreengrass (or if -o, some other path)"`
}

// Version will output the current version and quit.
func (args) Version() string {
	return "gogreengrass v0.0.6"
}

func main() {
	var args args
	arg.MustParse(&args)
	runner(args)
}

func runner(args args) {

	if args.Out == "" {
		args.Out = "/tmp/gogreengrass"
	}
	if args.SDKC {
		if err := writeSoFile(args.Out); nil != err {
			fmt.Printf("Failed to write the /tmp/gogreengrass/libaws-greengrass-core-sdk-c.so, error: %s\n", err.Error())
		}
	}
}

func writeFile(path, file string, data []byte) {

	if "" != path {
		os.MkdirAll(path, os.ModePerm)
	}

	fp := filepath.Join(path, file)
	fmt.Println("Writing file: ", fp)

	if f, err := os.Create(fp); err == nil {

		f.Write(data)
		f.Sync()
		f.Close()

	}
}

func extract(filepath string) {
	r, err := os.Open(filepath)
	if err != nil {
		log.Fatalf(err.Error())
	}

	defer r.Close()
	extractFile(r)

}
func extractFile(gzipStream io.Reader) {
	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		log.Fatal("ExtractTarGz: NewReader failed")
	}

	tarReader := tar.NewReader(uncompressedStream)

	for true {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("ExtractTarGz: Next() failed: %s", err.Error())
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.Mkdir(header.Name, 0755); err != nil {
				log.Fatalf("ExtractTarGz: Mkdir() failed: %s", err.Error())
			}
		case tar.TypeReg:
			outFile, err := os.Create(header.Name)
			if err != nil {
				log.Fatalf("ExtractTarGz: Create() failed: %s", err.Error())
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				log.Fatalf("ExtractTarGz: Copy() failed: %s", err.Error())
			}
			outFile.Close()

		default:
			log.Fatalf(
				"ExtractTarGz: unknown type: %v in %s",
				header.Typeflag,
				header.Name)
		}

	}
}

func writeSoFile(path string) error {

	data, err := base64.StdEncoding.DecodeString(string(soFile))
	if err != nil {
		return err
	}

	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return err
	}

	f := reader.File[0]
	r, err := f.Open()
	if err != nil {
		return err
	}

	defer r.Close()

	os.Mkdir(path, os.ModePerm)
	dst, err := os.Create(filepath.Join(path, "libaws-greengrass-core-sdk-c.so"))
	if err != nil {
		return err
	}

	defer func() {
		dst.Sync()
		dst.Close()
	}()

	io.Copy(dst, r)
	return nil
}
