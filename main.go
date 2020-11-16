package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/alexflint/go-arg"
)

//go:generate go run generators/generator.go

type args struct {
	Out         string `arg:"-o" help:"The out path to write the generated go and python files. Default is current directory" placeholder:"PATH"`
	Package     string `arg:"-p" help:"an optional package name instead of main" placeholder:"PACKAGE"`
	Binary      string `arg:"-b" help:"an optional name of the binary that the build system produces, default is foldername.o"`
	DownloadSDK bool   `arg:"-d" help:"If set to true, it will download the python sdk (it is needed to be in current folder)"`
}

func (args) Version() string {
	return "gogreengrass v0.0.1"
}

func main() {
	var args args
	arg.MustParse(&args)
	runner(args)
}

func runner(args args) {

	var goFile []byte
	var pyFile []byte

	if "" != args.Package {
		goFile = bytes.Replace(glueGo, []byte("package main"), []byte("package "+args.Package), 1)
	} else {
		goFile = glueGo
	}

	if "" == args.Binary {
		if path, err := os.Getwd(); err == nil {
			args.Binary = filepath.Base(path)
		}
	}

	pyFile = bytes.Replace(gluePy, []byte("./main.so"), []byte(args.Binary+".so"), 1)

	writeFile(args.Out, "glue.go", goFile)
	writeFile(args.Out, "glue.py", pyFile)

	if args.DownloadSDK {
		downloadGreengrassSDK()
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

func downloadGreengrassSDK() {

	sdk := "greengrasssdk-1.6.0"
	fileName := sdk + ".tar.gz"
	URL := "https://files.pythonhosted.org/packages/7f/d8/a17d97ba00275c13f3d0c6c82485aa6aa3ca9c24a61b3e2eae0fadee3d1b/" + fileName

	err := downloadFile(URL, fileName)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Printf("File %s downloaded in current working directory\n", fileName)

	os.RemoveAll(sdk)
	extract(fileName)

	os.Rename(filepath.Join(sdk, "greengrasssdk"), "./greengrasssdk")
	os.RemoveAll(sdk)
	os.Remove(fileName)
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

func downloadFile(URL, fileName string) error {
	//Get the response bytes from the url
	response, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("Received non 200 response code")
	}
	//Create a empty file
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	//Write the bytes to the fiel
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}
