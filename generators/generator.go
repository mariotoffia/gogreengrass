package main

import (
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

	fmt.Println("Finished generated.go")
}
