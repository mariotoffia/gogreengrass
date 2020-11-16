package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alexflint/go-arg"
)

//go:generate go run generators/generator.go

type args struct {
	Out     string `arg:"-o" help:"The out path to write the generated go and python files. Default is current directory" placeholder:"PATH"`
	Package string `arg:"-p" help:"an optional package name instead of main" placeholder:"PACKAGE"`
	Binary  string `arg:"-b" help:"an optional name of the binary that the build system produces, default is foldername.o"`
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
