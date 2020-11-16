package sdk

/*
#include <stdlib.h>

typedef void (*publish) (char* topic, char *policy, char *payload);
typedef void (*getThingShadow) (char *ctx, char* thingName);
typedef void (*deleteThingShadow) (char *ctx, char* thingName);
typedef void (*updateThingShadow) (char *ctx, char* thingName, char *payload);

static inline void call_c_publish(publish ptr,char* topic, char *policy, char *payload) {
	(ptr)(topic, policy, payload);
}

static inline void call_c_getThingShadow(getThingShadow ptr,char *ctx,char* thingName) {
	(ptr)(ctx, thingName);
}

static inline void call_c_updateThingShadow(updateThingShadow ptr,char *ctx,char* thingName, char *payload) {
	(ptr)(ctx, thingName, payload);
}

static inline void call_c_deleteThingShadow(deleteThingShadow ptr,char *ctx,char* thingName) {
	(ptr)(ctx, thingName);
}
*/
import "C"

import (
	"encoding/json"
	"fmt"
	"unsafe"
)

// QueueFullPolicy specifies what to do when queue is full while
// publishing data to queue
type QueueFullPolicy string

const (
	// QueueFullPolicyAllOrException - TODO:
	QueueFullPolicyAllOrException QueueFullPolicy = "AllOrException"
	// QueueFullPolicyBestEffort - TODO:
	QueueFullPolicyBestEffort QueueFullPolicy = "BestEffort"
)

// DataplaneClient is a interface that implements functions to
// reach the green grass environment.
type DataplaneClient interface {
	// Publish will publish the data onto specified topic.
	Publish(topic string, policy QueueFullPolicy, data string)
	// Publish will publish the object, marshalled as _JSON_, onto specified topic.
	PublishObject(topic string, policy QueueFullPolicy, object interface{})
	// GetThingShadow retrieves the thingshadow state for a specific
	// _thingName_.
	//
	// This shadow is synchronized, when connected, with the IoT Core
	// device shadow (default, not named).
	GetThingShadow(thingName string) (string, error)
	// UpdateThingShadow updates the shadow with state specified in the
	// _payload_ data. If it succeeds it will return the result of the
	// operation, otherwise an error is returned.
	//
	// This data is then synchronized with cloud IoT Core when connected.
	UpdateThingShadow(thingName string, payload string) (string, error)
	// DeleteThingShadow deletes the thing shadow and returns the output
	// from the DeleteThingShadow operation if successful. Otherwise
	// and error is returned.
	//
	// This operation is replicated to cloud IoT Core when connected.
	DeleteThingShadow(thingName string) (string, error)
}

type ggDataplane struct {
	fnPublish           *[0]byte
	fnGetThingShadow    *[0]byte
	fnUpdateThingShadow *[0]byte
	fnDeleteThingShadow *[0]byte
}

// NewDataplaneClient creates a new dataplane client to communicate with
// the greengrass dataplane API.
func NewDataplaneClient(
	fnPublish *[0]byte,
	fnGetThingShadow *[0]byte,
	fnUpdateThingShadow *[0]byte,
	fnDeleteThingShadow *[0]byte) DataplaneClient {

	return &ggDataplane{
		fnPublish:           fnPublish,
		fnGetThingShadow:    fnGetThingShadow,
		fnUpdateThingShadow: fnUpdateThingShadow,
		fnDeleteThingShadow: fnDeleteThingShadow,
	}

}

// Publish will publish the data onto specified topic.
func (ggf *ggDataplane) Publish(topic string, policy QueueFullPolicy, data string) {

	if nil == ggf.fnPublish {
		return
	}

	ct := C.CString(topic)
	cp := C.CString(string(policy))
	cd := C.CString(data)

	defer func() {
		C.free(unsafe.Pointer(ct))
		C.free(unsafe.Pointer(cp))
		C.free(unsafe.Pointer(cd))
	}()

	C.call_c_publish(ggf.fnPublish, ct, cp, cd)

}

// Publish will publish the object, marshalled as _JSON_, onto specified topic.
func (ggf *ggDataplane) PublishObject(topic string, policy QueueFullPolicy, object interface{}) {
	if nil == object {
		return
	}

	if data, err := json.Marshal(object); nil == err {
		ggf.Publish(topic, policy, string(data))
	}
}

// GetThingShadow retrieves the thingshadow state for a specific
// _thingName_.
//
// This shadow is synchronized, when connected, with the IoT Core
// device shadow (default, not named).
func (ggf *ggDataplane) GetThingShadow(thingName string) (string, error) {

	if nil == ggf.fnGetThingShadow {

		return "",
			fmt.Errorf(
				"Could not fetch thing shadow for %s since no func defined", thingName,
			)

	}

	ctx := getContext()

	ct := C.CString(thingName)
	cc := C.CString(ctx)

	defer func() {
		C.free(unsafe.Pointer(ct))
		C.free(unsafe.Pointer(cc))
	}()

	C.call_c_getThingShadow(ggf.fnGetThingShadow, cc, ct)

	buf := getProcessBuffer(ctx)
	if "" == buf {
		return "", nil /*TODO: error*/
	}

	return buf, nil
}

// UpdateThingShadow updates the shadow with state specified in the
// _payload_ data. If it succeeds it will return the result of the
// operation, otherwise an error is returned.
//
// This data is then synchronized with cloud IoT Core when connected.
func (ggf *ggDataplane) UpdateThingShadow(thingName string, payload string) (string, error) {

	if nil == ggf.fnUpdateThingShadow {

		return "",
			fmt.Errorf(
				"Could not fetch thing shadow for %s since no func defined", thingName,
			)

	}

	ctx := getContext()

	ct := C.CString(thingName)
	cc := C.CString(ctx)
	cp := C.CString(string(payload))

	defer func() {
		C.free(unsafe.Pointer(ct))
		C.free(unsafe.Pointer(cc))
		C.free(unsafe.Pointer(cp))
	}()

	C.call_c_updateThingShadow(ggf.fnUpdateThingShadow, cc, ct, cp)

	buf := getProcessBuffer(ctx)
	if "" == buf {
		return "", nil /*TODO: error*/
	}

	return buf, nil
}

// DeleteThingShadow deletes the thing shadow and returns the output
// from the DeleteThingShadow operation if successful. Otherwise
// and error is returned.
//
// This operation is replicated to cloud IoT Core when connected.
func (ggf *ggDataplane) DeleteThingShadow(thingName string) (string, error) {

	if nil == ggf.fnDeleteThingShadow {

		return "",
			fmt.Errorf(
				"Could not delete thing shadow for %s since no func defined", thingName,
			)

	}

	ctx := getContext()

	ct := C.CString(thingName)
	cc := C.CString(ctx)

	defer func() {
		C.free(unsafe.Pointer(ct))
		C.free(unsafe.Pointer(cc))
	}()

	C.call_c_deleteThingShadow(ggf.fnDeleteThingShadow, cc, ct)

	buf := getProcessBuffer(ctx)
	if "" == buf {
		return "", nil /*TODO: error*/
	}

	return buf, nil

}
