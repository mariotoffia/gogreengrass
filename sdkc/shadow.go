package sdkc

/*
#include "lib/greengrasssdk.h"
*/
import "C"

import (
	"unsafe"
)

// TODO: remove these links when done
// https://github.com/aws/aws-greengrass-core-sdk-c/tree/master/aws-greengrass-core-sdk-c-example
// https://golang.org/cmd/cgo/
// https://docs.aws.amazon.com/greengrass/latest/developerguide/lambda-functions.html#lambda-executables

// ShadowAPI encapsulates communication with the device shadow API.
type ShadowAPI struct {
	APIRequest
	response string
}

// GetLastResponse returns the last success response from an API call.
func (sa *ShadowAPI) GetLastResponse() string {
	return sa.response
}

// NewShadowAPI creates a new Shadow manager.
func NewShadowAPI() *ShadowAPI {
	return &ShadowAPI{}
}

// Update will update the device shadow locally on the greengrass
// system.
func (sa *ShadowAPI) Update(thingName, payload string) *ShadowAPI {

	if sa.err != nil {
		return sa
	}

	tn := C.CString(thingName)
	pl := C.CString(payload)

	defer func() {
		sa.close()
		C.free(unsafe.Pointer(tn))
		C.free(unsafe.Pointer(pl))
	}()

	res := C.gg_request_result{}

	sa.APIRequest.initialize()

	e := GreenGrassCode(C.gg_update_thing_shadow(sa.request, tn, pl, &res))

	sa.response = sa.handleRequestResponse(
		"Failed to update shadow",
		e, RequestStatus(res.request_status),
	)

	return sa
}

// Get will get the device shadow from local device.
func (sa *ShadowAPI) Get(thingName string) *ShadowAPI {

	if sa.err != nil {
		return sa
	}

	tn := C.CString(thingName)

	defer func() {
		C.free(unsafe.Pointer(tn))
		sa.APIRequest.close()
	}()

	res := C.gg_request_result{}

	sa.APIRequest.initialize()

	e := GreenGrassCode(C.gg_get_thing_shadow(sa.request, tn, &res))

	sa.response = sa.handleRequestResponse(
		"Failed to get shadow",
		e, RequestStatus(res.request_status),
	)

	return sa
}

// Delete will remove the device shadow from local device.
func (sa *ShadowAPI) Delete(thingName string) *ShadowAPI {

	if sa.err != nil {
		return sa
	}

	tn := C.CString(thingName)

	defer func() {
		C.free(unsafe.Pointer(tn))
		sa.APIRequest.close()
	}()

	res := C.gg_request_result{}

	sa.APIRequest.initialize()

	e := GreenGrassCode(C.gg_delete_thing_shadow(sa.request, tn, &res))

	sa.response = sa.handleRequestResponse(
		"Failed to delete shadow",
		e, RequestStatus(res.request_status),
	)

	return sa
}
