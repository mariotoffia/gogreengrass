package sdkc

/*
#include "lib/greengrasssdk.h"
*/
import "C"

import "unsafe"

// LambdaInvokeType specifies if the invoke is request / response or event based
type LambdaInvokeType int

const (
	// InvokeEvent specified that the invoke is asynchronously
	InvokeEvent = 0
	// InvokeRequestResponse specifies that the invoke is synchronously (default)
	InvokeRequestResponse = 1
	// InvokeReservedMax is reserved
	InvokeReservedMax = 2
	// InvokeReservedPad is reserved
	InvokeReservedPad = 0x7FFFFFFF
)

// LambdaInvokerAPI is to invoke a lambda and possibly handle
// the result.
type LambdaInvokerAPI struct {
	APIRequest
	response string
}

// NewLambdaInvokerAPI creates a new instance of the API to invoke
// GGC lambdas.
func NewLambdaInvokerAPI() *LambdaInvokerAPI {
	return &LambdaInvokerAPI{}
}

// GetLastResponse returns the last success response from an API call.
func (li *LambdaInvokerAPI) GetLastResponse() string {
	return li.response
}

// Invoke invokes the specified lambda either request/response or event style.
//
// The _functionARN_ is the full lambda ARN to be invoked and the _clientContext_
// is a base64-encoded null-terminated json string.
//
// .Example clientContext
// [source,json]
// ....
//  {"custom": {"value": "key"} } <1>
// ....
// <1> This translates to _clientContext_ of "eyAiY3VzdG9tIjp7ICJ2YWx1ZSI6ICJrZXkiIH19"
//
// The _invokeType_ specifies if this is a fire and forget (event) or if it is a request /
// response. The _functionVersion_ is a string representing which version e.g. "2"
func (li *LambdaInvokerAPI) Invoke(
	functionARN, clientContext, functionVersion string,
	invokeType LambdaInvokeType,
	payload []byte) *LambdaInvokerAPI {

	if li.err != nil {
		return li
	}

	fa := C.CString(functionARN)
	cc := C.CString(clientContext)
	fv := C.CString(functionVersion)

	defer func() {
		li.close()
		C.free(unsafe.Pointer(fa))
		C.free(unsafe.Pointer(cc))
		C.free(unsafe.Pointer(fv))
	}()

	li.initialize()

	opts := C.gg_invoke_options{
		function_arn:     fa,
		customer_context: cc,
		qualifier:        fv,
		_type:            C.gg_invoke_type(invokeType),
		payload:          unsafe.Pointer(&payload[0]),
		payload_size:     C.size_t(len(payload)),
	}

	res := C.gg_request_result{}
	e := GreenGrassCode(C.gg_invoke(li.request, &opts, &res))

	li.response = li.handleRequestResponse(
		"Failed to invoke lambda",
		e, RequestStatus(res.request_status),
	)

	return li
}
