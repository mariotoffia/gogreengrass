package sdkc

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambda/messages"
)

const (
	msPerSecond = int64(time.Second / time.Millisecond)
	nsPerMS     = int64(time.Millisecond / time.Nanosecond)
)

type ggContext struct {
	RequestID          string                 `json:"aws_request_id"`
	InvokedFunctionARN string                 `json:"invoked_function_arn"`
	FunctionName       string                 `json:"function_name"`
	FunctionVersion    string                 `json:"function_version"`
	Identity           *json.RawMessage       `json:"identity,omitempty"`
	ClientContext      *json.RawMessage       `json:"client_context,omitempty"`
	Headers            map[string]interface{} `json:"headers,omitempty"`
}

func (ggc *ggContext) getHeader(name string) string {

	if val, ok := ggc.Headers[name]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}

	return ""
}

var function *lambda.Function

// Start takes a handler function. See lambda github project for valid signatures
// https://github.com/aws/aws-lambda-go/blob/f24acb29a08c3a45eb95e6cd4ae56fbfabf4f4a5/lambda/entry.go#L39
//
// This will start the dispatcher on same thread and will freeze
func Start(handler interface{}) {
	StartWithOpts(RuntimeOptionSingleThread, handler)
}

// StartWithOpts starts the lambda runtime with the specified option and registers the lambda callback.
//
// Documentation on lambda function layout:
// https://github.com/aws/aws-lambda-go/blob/f24acb29a08c3a45eb95e6cd4ae56fbfabf4f4a5/lambda/entry.go#L39
func StartWithOpts(option RuntimeOption, handler interface{}) {
	function = lambda.NewFunction(lambda.NewHandler(handler))

	GGStartWithOpts(
		option,
		func(lc *LambdaContextSlim) {

			resp, err := invokeJSONRequest(lc.ClientContext, lc.Payload)

			if err != nil {
				lambdaWriteError(err.Error())
			} else {
				lambdaWriteResponse(resp)
			}

		}, true /*payload*/)

}

func invokeJSONRequest(context string, payload []byte) ([]byte, error) {

	req, err := createRequest([]byte(context), payload)

	if err != nil {
		fmt.Printf("%s", err.Error())
	}

	var resp messages.InvokeResponse
	if err := function.Invoke(req, &resp); err != nil {
		return nil, fmt.Errorf("Invoke error: %s", err.Error())
	}

	if resp.Error != nil {
		return nil, fmt.Errorf("%s", err.Error())
	}

	return resp.Payload, nil
}

func createRequest(context []byte, payload []byte) (*messages.InvokeRequest, error) {

	var ctx ggContext
	if err := json.Unmarshal(context, &ctx); nil != err {
		return nil, fmt.Errorf("Failed to unmarshal context: %s, error: %s", string(context), err.Error())
	}

	deadlineEpochMS, err := strconv.ParseInt(ctx.getHeader("Lambda-Runtime-Deadline-Ms"), 10, 64)

	if err != nil {
		deadlineEpochMS = 30000 /* 3s */
	}

	res := &messages.InvokeRequest{
		InvokedFunctionArn: ctx.InvokedFunctionARN,
		XAmznTraceId:       ctx.getHeader("Lambda-Runtime-Trace-Id"),
		RequestId:          ctx.RequestID,
		Deadline: messages.InvokeRequest_Timestamp{
			Seconds: deadlineEpochMS / msPerSecond,
			Nanos:   (deadlineEpochMS % msPerSecond) * nsPerMS,
		},
		Payload: payload,
	}

	if ctx.ClientContext != nil {
		res.ClientContext = []byte((*ctx.ClientContext))
	}

	if ctx.Identity != nil {
		if err := json.Unmarshal([]byte((*ctx.Identity)), res); err != nil {

			return nil, fmt.Errorf(
				"failed to unmarshal cognito identity json: %s, error: %s",
				ctx.Identity, err.Error(),
			)

		}
	}

	return res, nil
}
