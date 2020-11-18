package sdkc

/*
#cgo LDFLAGS: -L/tmp/gogreengrass -laws-greengrass-core-sdk-c -Wl,-rpath=/tmp/gogreengrass
#include "lib/greengrasssdk.h"

extern void handler(gg_lambda_context *ctx);

typedef void (*Thandler)(gg_lambda_context *ctx);
*/
import "C"
import (
	"fmt"
	"io"
	"io/ioutil"
	"unsafe"
)

// TODO: remove these links when done
// https://github.com/aws/aws-greengrass-core-sdk-c/tree/master/aws-greengrass-core-sdk-c-example
// https://golang.org/cmd/cgo/
// https://docs.aws.amazon.com/greengrass/latest/developerguide/lambda-functions.html#lambda-executables

// LambdaContextSlim slim version of the lambda context
//
// If the `LambdaHandler` is registered as payload set to `false`
// the _Payload_ field will not be populated and the `LambdaHandler`
// need to process the payload itself.
type LambdaContextSlim struct {
	ClientContext string
	FunctionARN   string
	Payload       []byte
}

// LambdaHandler is the handler function that processes incoming requests.
//
// Register the LambdaHandler using either `Start` or `StartWithOpts`
//
// .Example Single Threaded Registration
// [source,go]
// ....
// sdkc.Start(func(lc *sdkc.LambdaContextSlim) { // <1>
//
// 	 fmt.Printf("%s, %s\n", lc.ClientContext, lc.FunctionARN) // <2>
// 	 fmt.Println("Payload: ", string(lc.Payload)) // <3>
//
// })
// ....
// <1>
type LambdaHandler func(lc *LambdaContextSlim)

// RuntimeOption specifies how the runtime shall behave or initialized
type RuntimeOption uint

const (
	// RuntimeOptionSingleThread will start the runtime and register the lambda function and
	// process all request on the caller thread.
	RuntimeOptionSingleThread RuntimeOption = 0
	// RuntimeOptionSeparateThread will start the runtime and register the lambda function in
	// a new thread. When the caller thread / main thread exits this runtime thread also exits.
	RuntimeOptionSeparateThread RuntimeOption = 1
)

// GGStart will start a single threaded runtime and do callback onto the registered `LambdaHandler`
// with full payload.
//
// Since single-threaded this function will freeze the current thread. If you want more control
// over the runtime and how decoding of `LambdaContextSlim` is done use `StartWithOpts` instead.
func GGStart(lh LambdaHandler) {
	GGStartWithOpts(RuntimeOptionSeparateThread, lh, true /*payload*/)
}

// GGStartWithOpts will start the runtime and register the lambda function with options.
//
// When the _payload_ parameter is set to `false` the `LambdaContextSlim.Payload` will
// be empty. In this case the `LambdaHandler` need to read request data itself using
// the `RequestReader` (using `NewRequestReader()` to create one).
func GGStartWithOpts(option RuntimeOption, lh LambdaHandler, payload bool) {

	decodePayload = payload
	regHandler = lh

	C.gg_global_init(0)

	/* start the runtime in blocking mode. This blocks forever. */
	C.gg_runtime_start((C.Thandler)(unsafe.Pointer(C.handler)), C.uint(option))
}

func lambdaWriteResponse(payload []byte) error {

	errorCode := GreenGrassCode(
		C.gg_lambda_handler_write_response(
			unsafe.Pointer(&payload[0]),
			C.size_t(len(payload)),
		),
	)

	if errorCode == GreenGrassCodeSuccess {
		return nil
	}

	return fmt.Errorf(
		"Got error %d when writing response payload", errorCode,
	)

}

// lambdaWriteError will return an error for a lambda invocation.
func lambdaWriteError(errorMessage string) error {

	em := C.CString(errorMessage)

	defer C.free(unsafe.Pointer(em))

	errorCode := GreenGrassCode(C.gg_lambda_handler_write_error(em))

	if errorCode == GreenGrassCodeSuccess {
		return nil
	}

	return fmt.Errorf(
		"Got error %d when writing error response: '%s'", errorCode, errorMessage,
	)

}

// NewRequestReader creates a new request reader.
//
// .Example Usage
// [source,go]
// ....
// r := NewRequestReader()
// b := make([]byte, 256)
// for {
//
//	 n, err := 	r.Read(b)
//   if n > 0 {
//     // process the b[:n]
//   }
//
//   if err == io.EOF {
//	   break
//   }
//
//   if err != nil {
//	   panic(err)
//   }
// }
// ....
//
// Or use the `ReadAll` from _ioutil_ package.
//
// .Example ReadAll
// [source,go]
// ....
// if buf, err := ioutil.ReadAll(NewRequestReader()); err == nil {
//   // the complete request payload is in buf
// }
// ....
func NewRequestReader() *RequestReader {
	return &RequestReader{}
}

// RequestReader reads the lambda request data either
// iteratively (low memory footprint) or in full. It
// implements the `io.Reader` interface
type RequestReader struct{}

// Read implements the `io.Reader` interface method to read from
// the lambda request buffer.
//
// When the read from buffer is done a `io.EOF` is returned.
//
// NOTE: Even if `io.EOF` is returned, some data may exist in buffer to be handeled.
func (r *RequestReader) Read(b []byte) (n int, err error) {

	var amountRead C.size_t

	errCode := C.gg_lambda_handler_read(
		unsafe.Pointer(&b[0]),
		C.size_t(len(b)),
		&amountRead)

	n = int(amountRead)

	if n == 0 {
		err = io.EOF
		return
	}

	ggError := GreenGrassCode(errCode)

	if ggError != GreenGrassCodeSuccess {
		err = fmt.Errorf("failure to read from handler, code: %d", ggError)
	}

	return
}

var regHandler LambdaHandler
var decodePayload = false

//export handler
//
// This handler is the callback from the C environment that will call the registered
// `LambdaHandler` (if any).
func handler(ctx *C.gg_lambda_context) {

	if nil == regHandler {
		return
	}

	lc := LambdaContextSlim{
		FunctionARN:   C.GoString(ctx.function_arn),
		ClientContext: C.GoString(ctx.client_context),
	}

	if decodePayload {

		data, err := ioutil.ReadAll(NewRequestReader())

		if err != nil {

			Log(LogLevelError, "Failed to read data from request, error: %s", err.Error())
			return

		}

		lc.Payload = data
	}

	regHandler(&lc)
}
