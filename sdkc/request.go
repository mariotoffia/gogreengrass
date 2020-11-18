package sdkc

/*
#include "lib/greengrasssdk.h"
*/
import "C"
import (
	"fmt"
	"io/ioutil"
)

// RequestStatus is the return code from an API invocation
type RequestStatus int

const (
	// RequestStatusSuccess is when function call returns expected payload type
	RequestStatusSuccess = 0
	// RequestStatusHandled is when function call is successfull, however
	// lambda response with an error
	RequestStatusHandled = 1
	// RequestStatusUnhandled is when function call is unsuccessfull,
	// lambda exits abnormally
	RequestStatusUnhandled = 2
	// RequestStatusUnknown is when system encounters unknown error.
	// Check logs for more details
	RequestStatusUnknown = 3
	// RequestStatusRequestAgain is when function call is throttled, try again
	RequestStatusRequestAgain = 4
	// RequestStatusReservedMax is last valid const for this type
	RequestStatusReservedMax = 5
	// RequestStatusReservedPad is padding
	RequestStatusReservedPad = 0x7FFFFFFF
)

func createErrorFromRequestStatus(errorMessage string, status RequestStatus) error {
	if status == RequestStatusSuccess {
		return nil
	}

	err := fmt.Errorf("%s, status: %d", errorMessage, status)
	Log(LogLevelInfo, err.Error()+"\n")

	return err
}

// APIRequest is a wrapper to hold a pointer
// for the request duration. This is used
// for all GreenGrass API invocations such
// as device shadow, secrets manager, publish.
type APIRequest struct {
	request C.gg_request
	// Last error occurred.
	err error
}

// NewAPIRequest creates a new API request.
func (ar *APIRequest) initialize() {

	e := GreenGrassCode(C.gg_request_init(&ar.request))
	ar.err = createError("Failed to do gg_request_init()", e)

}

// Error returns the last error occurred.
func (ar *APIRequest) Error() error {
	return ar.err
}

// ClearError clears any error state
func (ar *APIRequest) ClearError() *APIRequest {
	ar.err = nil
	return ar
}

// Close finalizes the API call and free up resources.
func (ar *APIRequest) close() *APIRequest {

	if ar.err != nil {
		return ar
	}

	e := GreenGrassCode(C.gg_request_close(ar.request))
	ar.err = createError("Failed to do gg_request_close()", e)

	return ar

}

func (ar *APIRequest) handleRequestResponse(msg string, e GreenGrassCode, re RequestStatus) string {

	ar.err = createError("Failed to update shadow", e)
	if ar.err != nil {
		return ""
	}

	// read success msg or error msg
	s := ""
	if data, err := ioutil.ReadAll(NewRequestReader()); err == nil {
		s = string(data)
	}

	if e != RequestStatusSuccess {

		ar.err = createErrorFromRequestStatus(
			fmt.Sprintf("%s, msg: '%s'", msg, s), re,
		)

		return ""
	}

	// success: return the response
	return s
}
