package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/mariotoffia/gogreengrass/sdk"
)

// once is invoked once when ggc lambda startup
func once() {

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
func main() {
	// Standard AWS cloud lambda initialization code
}
