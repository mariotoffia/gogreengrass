package sdk

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-lambda-go/lambdacontext"
)

func TestInvokeSimple(t *testing.T) {

	type MyEvent struct {
		Data  int    `json:"data"`
		Hello string `json:"hello"`
	}

	type MyResponse struct {
		Age   int    `json:"age"`
		Topic string `json:"topic"`
	}

	Register(func(c context.Context, data MyEvent) (MyResponse, error) {

		lc, _ := lambdacontext.FromContext(c)

		fmt.Printf(
			"context: %v, topic: %s, data: '%v'\n",
			lc, lc.ClientContext.Custom["subject"], data,
		)

		return MyResponse{Age: 19, Topic: "feed/myfunc"}, nil
	})

	event := `{"data": 22, "hello":"world"}`
	context := `{
		"aws_request_id": "99994f79-2e36-4eb2-6584-b533cb3bc491",
		"client_context": {
			"client": null, 
			"custom": {"subject": "invoke/ggtest"}, 
			"env": null
		},
		"function_name": "ggtest",
		"function_version": "7",
		"identity": null,
		"invoked_function_arn": "arn:aws:lambda:eu-central-1:033549287452:function:ggtest:7",
		"headers": {
			"Lambda-Runtime-Deadline-Ms": "300000"
		}
	}`

	result := InvokeJSON(context, event)

	fmt.Println(result)
}
