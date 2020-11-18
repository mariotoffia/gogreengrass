package main

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	fmt.Println("Generating embedded resources")

	gen, err := os.Create("generated.go")
	if err != nil {
		fmt.Println("Failed to create generated.go file:", err.Error())
		return
	}

	content, err := ioutil.ReadFile("./internal/files/glue.go")
	if err != nil {
		fmt.Println("Failed to read glue.go file:", err.Error())
		return
	}

	gen.Write([]byte("package main\n\n"))
	gen.Write([]byte("var glueGo = []byte(`"))
	gen.Write(content)
	gen.Write([]byte("`)\n\n"))

	content, err = ioutil.ReadFile("./internal/files/glue.py")
	if err != nil {
		fmt.Println("Failed to read glue.py file:", err.Error())
		return
	}

	gen.Write([]byte("var gluePy = []byte(`"))
	gen.Write(content)
	gen.Write([]byte("`)\n\n"))

	content, err = ioutil.ReadFile("/tmp/gogreengrass/libaws-greengrass-core-sdk-c.so")
	if err != nil {
		fmt.Println("failed to read /tmp/gogreengrass/libaws-greengrass-core-sdk-c.so: ", err.Error())
		return
	}

	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)
	zipFile, err := zipWriter.Create("libaws-greengrass-core-sdk-c.so")
	if err != nil {
		fmt.Println("Failed to create zip of libaws-greengrass-core-sdk-c.so")
		return
	}

	_, err = zipFile.Write(content)
	if err != nil {
		fmt.Println("Failed to create zip of libaws-greengrass-core-sdk-c.so")
		return
	}

	zipWriter.Flush()
	err = zipWriter.Close()
	if err != nil {
		fmt.Println("Failed to create zip of libaws-greengrass-core-sdk-c.so")
		return
	}

	content = []byte(base64.StdEncoding.EncodeToString(buf.Bytes()))

	gen.Write([]byte("var soFile = []byte(`"))
	gen.Write(content)
	gen.Write([]byte("`)\n\n"))

	fmt.Println("Finished generated.go")
}
