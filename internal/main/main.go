package main

/*
#include <stdlib.h>

typedef void (*publish) (char* topic, char *policy, char *payload);

static inline void call_c_func(publish ptr,char* topic, char *policy, char *payload) {
	(ptr)(topic, policy, payload);
}
*/
import "C"

import (
	"context"
	"fmt"
	"unsafe"

	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/mariotoffia/gogreengrass/sdk"
)

var fnPublish C.publish

//export initcb
func initcb(fn C.publish) {
	fnPublish = fn
}

// Publish will publish the data onto specified topic.
func Publish(topic string, policy sdk.QueueFullPolicy, data string) {

	if nil == fnPublish {
		return
	}

	ct := C.CString(topic)
	cp := C.CString(string(policy))
	cd := C.CString(data)

	defer func() {
		C.free(unsafe.Pointer(ct))
		C.free(unsafe.Pointer(cp))
		C.free(unsafe.Pointer(cd))
	}()

	C.call_c_func(fnPublish, ct, cp, cd)
}

//export invokeJSON
func invokeJSON(context string, payload string) *C.char {
	return C.CString(sdk.InvokeJSON(context, payload))
}

//export setup
func setup() {

	type MyEvent struct {
		Data  int    `json:"data"`
		Hello string `json:"hello"`
	}

	type MyResponse struct {
		Age   int    `json:"age"`
		Topic string `json:"topic"`
	}

	sdk.Register(func(c context.Context, data MyEvent) (MyResponse, error) {

		lc, _ := lambdacontext.FromContext(c)

		fmt.Printf(
			"context: %v, topic: %s, data: '%v'\n",
			lc, lc.ClientContext.Custom["subject"], data,
		)

		Publish("hello/world", sdk.QueueFullPolicyAllOrException, `{"my":"prop"}`)
		return MyResponse{Age: 19, Topic: "feed/myfunc"}, nil
	})

}

func main() {}
