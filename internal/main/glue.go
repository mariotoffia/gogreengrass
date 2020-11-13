package main

/*
typedef void (*publish) (char* topic, char *policy, char *payload);
*/
import "C"

import (
	"github.com/mariotoffia/gogreengrass/sdk"
)

// GGFunctions is the functions to call gg
var GGFunctions sdk.GreenGrassFunctions

//export initcb
func initcb(fnPublish C.publish) {
	GGFunctions = sdk.NewGreenGrassInterface(fnPublish)
}

//export setup
func setup() {
	Setup()
}

//export invokeJSON
func invokeJSON(context string, payload string) *C.char {
	return C.CString(sdk.InvokeJSON(context, payload))
}
