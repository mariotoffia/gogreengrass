package main

import (
	"C"
	"context"
	"fmt"

	"github.com/mariotoffia/gogreengrass/sdk"
)
import "github.com/aws/aws-lambda-go/lambdacontext"

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
		return MyResponse{Age: 19, Topic: "feed/myfunc"}, nil
	})

}

func main() {}
