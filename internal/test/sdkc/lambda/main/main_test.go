package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/mariotoffia/gogreengrass/sdkc"
)

func TestGreenGrassLambda(t *testing.T) {

	sdkc.GGStart(func(lc *sdkc.LambdaContextSlim) {

		fmt.Printf("%s, %s\n", lc.ClientContext, lc.FunctionARN)
		fmt.Println("Payload: ", string(lc.Payload))

	})
}

func TestStandardLambda(t *testing.T) {

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

		return MyResponse{Age: 19, Topic: "feed/myfunc"}, nil
	})

}
