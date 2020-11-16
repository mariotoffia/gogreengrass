package main

/*
typedef void (*publish) (char* topic, char *policy, char *payload);
typedef char* (*getThingShadow) (char* thingName);
typedef void (*updateThingShadow) (char *ctx, char* thingName, char *payload);
typedef void (*deleteThingShadow) (char *ctx, char* thingName);
*/
import "C"

import (
	"github.com/mariotoffia/gogreengrass/sdk"
)

// GGDataplane is the functions to for dataplane against GGC API
var GGDataplane sdk.DataplaneClient

//export initcb
func initcb(
	fnPublish C.publish,
	fnGetThingShadow C.getThingShadow,
	fnUpdateThingShadow C.updateThingShadow,
	fnDeleteThingShadow C.deleteThingShadow) {

	GGDataplane = sdk.NewDataplaneClient(
		fnPublish,
		fnGetThingShadow,
		fnUpdateThingShadow,
		fnDeleteThingShadow)

}

//export setup
func setup() {
	Setup()
}

//export set_process_buffer
func set_process_buffer(ctx *C.char, buffer *C.char) {
	sdk.SetProcessBuffer(C.GoString(ctx), C.GoString(buffer))
}

//export invokeJSON
func invokeJSON(context string, payload string) *C.char {
	return C.CString(sdk.InvokeJSON(context, payload))
}
