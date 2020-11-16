package main

import "C"
import (
	"fmt"
)

// Setup is the testrunner for this
func once() {
	secret, _ := GGSecretsManager.GetSecret("testId", "v1.0.0", "AWSCURRENT")
	fmt.Println(secret)
	fmt.Println("bin: " + string(secret.SecretBinary))
	fmt.Println("sec: " + secret.SecretString)
}

func main() {}
