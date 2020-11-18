package sdkc

/*
#cgo LDFLAGS: -L./ -laws-greengrass-core-sdk-c -Wl,-rpath=./
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

//export handler
func handler(ctx *C.gg_lambda_context) {
	clientContext := C.GoString(ctx.client_context)
	functionARN := C.GoString(ctx.function_arn)

	fmt.Printf("%s, %s\n", clientContext, functionARN)

	buff, _ := ioutil.ReadAll(NewRequestReader())
	fmt.Println(string(buff))

	// TODO: Invoke real handler.
	// https://github.com/aws/aws-greengrass-core-sdk-c/tree/master/aws-greengrass-core-sdk-c-example
	// https://github.com/lxwagn/using-go-with-c-libraries/blob/master/cgo.go
}

// RuntimeOption specifies how the runtime shall behave or initialized
type RuntimeOption int

const (
	// RuntimeOptionSingleThread will start the runtime and register the lambda function and
	// process all request on the caller thread.
	RuntimeOptionSingleThread RuntimeOption = 0
	// RuntimeOptionSeparateThread will start the runtime and register the lambda function in
	// a new thread. When the caller thread / main thread exits this runtime thread also exits.
	RuntimeOptionSeparateThread RuntimeOption = 0
)

// Start will start the runtime and register the lambda function.
func Start(option RuntimeOption) {

	C.gg_global_init(0)

	/* start the runtime in blocking mode. This blocks forever. */
	C.gg_runtime_start((C.Thandler)(unsafe.Pointer(C.handler)), C.uint(option))
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
