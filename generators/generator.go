package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	fmt.Println("Generating embedded resources")

	content, err := ioutil.ReadFile("main.py")
	if err != nil {
		fmt.Println("Failed to read main.py file:", err.Error())
		return
	}

	gen, err := os.Create("generated.go")
	if err != nil {
		fmt.Println("Failed to create generated.go file:", err.Error())
		return
	}

	gen.Write([]byte("package main\n\n"))
	gen.Write([]byte("var mainPy = []byte(`"))
	gen.Write(content)
	gen.Write([]byte("`)\n"))

	fmt.Println("Finished generated.go")
}
