package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/mariotoffia/gogreengrass/sdkc"
)

//go:generate gogreengrass -sdkc

func main() {
	sdkc.Log(sdkc.LogLevelInfo, "Register lambda...\n")

	type MyEvent struct {
		Data  int    `json:"data"`
		Hello string `json:"hello"`
	}

	type MyResponse struct {
		Age   int    `json:"age"`
		Topic string `json:"topic"`
	}

	sdkc.Start(func(c context.Context, data MyEvent) (MyResponse, error) {

		lc, _ := lambdacontext.FromContext(c)

		fmt.Printf(
			"context: %v, topic: %s, data: '%v'\n",
			lc, lc.ClientContext.Custom["subject"], data,
		)

		resp := MyResponse{Age: 19, Topic: "feed/myfunc"}

		sdkc.NewQueueAPI().PublishObject(
			"feed/testlambda", sdkc.QueueFullPolicyOptionAllOrError, &resp,
		)

		return resp, nil
	})

	sdkc.Log(sdkc.LogLevelInfo, "Exit main()\n")
}
