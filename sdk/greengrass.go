package sdk

/*
#include <stdlib.h>

typedef void (*publish) (char* topic, char *policy, char *payload);

static inline void call_c_func(publish ptr,char* topic, char *policy, char *payload) {
	(ptr)(topic, policy, payload);
}
*/
import "C"

import (
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

// GreenGrassFunctions is a interface that implements functions to
// reach the green grass environment.
type GreenGrassFunctions interface {
	// Publish will publish the data onto specified topic.
	Publish(topic string, policy QueueFullPolicy, data string)
}

type ggFunctions struct {
	fnPublish *[0]byte
}

// NewGreenGrassInterface creates a new interface bases on callback functions.
func NewGreenGrassInterface(fnPublish *[0]byte) GreenGrassFunctions {
	return &ggFunctions{
		fnPublish: fnPublish,
	}
}

// Publish will publish the data onto specified topic.
func (ggf *ggFunctions) Publish(topic string, policy QueueFullPolicy, data string) {

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

	C.call_c_func(ggf.fnPublish, ct, cp, cd)

}
