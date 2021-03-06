package sdkc

import (
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambda/messages"
)

var function *lambda.Function

// Start takes a handler function. See lambda github project for valid signatures
// https://github.com/aws/aws-lambda-go/blob/f24acb29a08c3a45eb95e6cd4ae56fbfabf4f4a5/lambda/entry.go#L39
//
// This will start the dispatcher on same thread and will freeze
func Start(handler interface{}) {
	StartWithOpts(RuntimeOptionSingleThread, handler, true /*payload*/)
}

// StartWithOpts starts the lambda runtime with the specified option and registers the lambda callback.
//
// If payload is set to false, the data sent to the lambda may be retrievable by `NewRequestReader()`.
// This is good for lambdas that do not want any payload or if the payload is large and it is possible
// to process this in parts. Hence, memory requirements is much lower.
//
// Documentation on lambda function layout:
// https://github.com/aws/aws-lambda-go/blob/f24acb29a08c3a45eb95e6cd4ae56fbfabf4f4a5/lambda/entry.go#L39
func StartWithOpts(option RuntimeOption, handler interface{}, payload bool) {

	function = lambda.NewFunction(lambda.NewHandler(handler))

	GGStartWithOpts(
		option,
		func(lc *LambdaContextSlim) {

			resp, err := invoke(lc)

			if err != nil {

				Log(LogLevelError, "Got error: %s\n", err.Error())
				lambdaWriteError(err.Error())

			} else {
				lambdaWriteResponse(resp)
			}

		},
		payload,
	)

}

func invoke(lc *LambdaContextSlim) ([]byte, error) {

	req := &messages.InvokeRequest{
		InvokedFunctionArn: lc.FunctionARN,
		XAmznTraceId:       "",
		RequestId:          "",
		Deadline: messages.InvokeRequest_Timestamp{
			Seconds: int64(9223372036854775800),
			Nanos:   0,
		},
		Payload: lc.Payload,
	}

	if lc.ClientContext != "" {
		req.ClientContext = []byte((lc.ClientContext))
	}

	var resp messages.InvokeResponse
	if err := function.Invoke(req, &resp); err != nil {
		return nil, fmt.Errorf("Invoke error: %s", err.Error())
	}

	if resp.Error != nil {
		return nil, fmt.Errorf("%s", resp.Error.Error())
	}

	return resp.Payload, nil
}
