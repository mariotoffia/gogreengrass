package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-lambda-go/lambda/messages"
)

// LambdaClient is a client to do lambda operations
type LambdaClient interface {
	// Invoke will invoke a lambda function by name or ARN.
	//
	// .Parameters
	// |===
	// |Name |Description
	//
	// |c
	// |Is the context, where the `ClientContext` is extracted from (if any)
	//
	// |functionName
	// a|
	// The Amazon Resource Name (ARN) of the Lambda function to invoke. Name formats:
	// * Qualified ARN - The function ARN with the version suffix. e.g. arn:aws:lambda:aws-region:acct-id:function:helloworld:1
	// * Unqualified ARN - The function ARN without the version suffix. e.g. arn:aws:lambda:aws-region:acct-id:function:helloworld
	//
	// |invocationType
	// a|
	// Two options below are valid:
	// * _RequestResponse_ (default) - Invoke the Lambda synchronously. Block until the function returns a response or times out.
	// * _Event_ - Invoke the Lambda asynchronously. The response only includes empty payload.
	//
	// |qualifier
	// |Optional parameter to specify a Lambda function version if it was not included in the FunctionName field.
	// If you specify a function version, the API uses the qualified function ARN to invoke a specific Lambda function.
	//
	// |payload
	// |Optional input for the Lambda function to invoke.
	// |===
	//
	// If it succeeds it may return the result of the invocation (based on the invocation type). Otherwise
	// an error returned.
	Invoke(c context.Context,
		functionName,
		invocationType,
		qualifier string,
		payload []byte) (io.Reader, error)
}

type invokeResult struct {
	// TODO: ref to something here passed back and forth..
}

// Read - calls the `StreamingBody` on the python side to get more data.
// When no more data _io.EOF_ is returned.
// TODO: Implement me!
func (ir *invokeResult) Read(p []byte) (n int, err error) {
	return 0, nil
}

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

// Register takes a handler function. See lambda github project for valid signatures
// https://github.com/aws/aws-lambda-go/blob/f24acb29a08c3a45eb95e6cd4ae56fbfabf4f4a5/lambda/entry.go#L39
func Register(handler interface{}) {
	function = lambda.NewFunction(lambda.NewHandler(handler))
}

// InvokeJSON invokes the registered handle
func InvokeJSON(context string, payload string) string {
	mtx.Lock()
	defer mtx.Unlock()

	req, err := createRequest([]byte(context), []byte(payload))

	if err != nil {
		fmt.Printf("%s", err.Error())
	}

	var resp messages.InvokeResponse
	if err := function.Invoke(req, &resp); err != nil {

		return fmt.Sprintf(
			`{"errorMessage": "%s", "errorType": "%s"}`, err.Error(), "invoke",
		)

	}

	if resp.Error != nil {

		return fmt.Sprintf(
			`{"errorMessage": "%s", "errorType": "%s"}`, resp.Error.Message, resp.Error.Type,
		)

	}

	payload = string(resp.Payload)
	return fmt.Sprintf(`{"Payload": %s}`, payload)
}

const (
	msPerSecond = int64(time.Second / time.Millisecond)
	nsPerMS     = int64(time.Millisecond / time.Nanosecond)
)

func createRequest(context []byte, payload []byte) (*messages.InvokeRequest, error) {

	var ctx ggContext
	if err := json.Unmarshal(context, &ctx); nil != err {
		return nil, fmt.Errorf("Failed to unmarshal context: %s, error: %s", string(context), err.Error())
	}

	deadlineEpochMS, err := strconv.ParseInt(ctx.getHeader("Lambda-Runtime-Deadline-Ms"), 10, 64)

	if err != nil {

		return nil, fmt.Errorf(
			"Expecting header Lambda-Runtime-Deadline-Ms to be present in a invocation",
		)

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
