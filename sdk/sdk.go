package sdk

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"os"

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

// GGEnvironment is a set of well known GreenGrass
// environment parameters.
type GGEnvironment struct {
	AuthToken                  string `env:"AWS_CONTAINER_AUTHORIZATION_TOKEN"`
	FunctionArn                string `env:"MY_FUNCTION_ARN"`
	EncodingType               string `env:"ENCODING_TYPE"`
	ShadowFunctionArn          string `env:"SHADOW_FUNCTION_ARN"`
	RouterFunctionArn          string `env:"ROUTER_FUNCTION_ARN"`
	GGCMaxInterfaceVersion     string `env:"GGC_MAX_INTERFACE_VERSION"`
	SecretesManagerFunctionArn string `env:"SECRETS_MANAGER_FUNCTION_ARN"`
	GGDeamonPort               string `env:"GG_DAEMON_PORT"`
}

// GetGGEnvironment fetches a new instance of `GetGGEnvironment` where it has
// be populated from the environment variables.
func GetGGEnvironment() GGEnvironment {

	return GGEnvironment{
		AuthToken:                  os.Getenv("AWS_CONTAINER_AUTHORIZATION_TOKEN"),
		FunctionArn:                os.Getenv("MY_FUNCTION_ARN"),
		EncodingType:               os.Getenv("ENCODING_TYPE"),
		ShadowFunctionArn:          os.Getenv("SHADOW_FUNCTION_ARN"),
		RouterFunctionArn:          os.Getenv("ROUTER_FUNCTION_ARN"),
		GGCMaxInterfaceVersion:     os.Getenv("GGC_MAX_INTERFACE_VERSION"),
		SecretesManagerFunctionArn: os.Getenv("SECRETS_MANAGER_FUNCTION_ARN"),
		GGDeamonPort:               os.Getenv("GG_DAEMON_PORT"),
	}

}

func (ggc *ggContext) getHeader(name string) string {

	if val, ok := ggc.Headers[name]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}

	return ""
}

var mtx sync.Mutex
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

/*

 client.publish(
                topic="feed/ggtest",
                queueFullPolicy="AllOrException",
                payload=json.dumps(
                    {
                        "message": "Hello world! Sent from Greengrass Core running on platform: {}.".format(my_platform)
                        + "  Invocation Count: {}".format(my_counter)
                    }
                ),
			)

			-->

Publishing message on topic "feed/ggtest" with Payload "
{
    "message": "Hello world! Sent from Greengrass Core running on platform: Linux-4.19.104-microsoft-standard-x86_64-with-glibc2.2.5.  Invocation Count: 2"
}*/
