package sdkc

/*
#include "lib/greengrasssdk.h"
*/
import "C"
import (
	"encoding/json"
	"fmt"
	"unsafe"
)

// QueueFullPolicyOption specifies what to do when queue is full.
type QueueFullPolicyOption int

const (
	// QueueFullPolicyOptionBestEffort sets publishing at best effort.
	QueueFullPolicyOptionBestEffort = 0
	// QueueFullPolicyOptionAllOrError specifies that GGC will either deliver
	// messages to all targets and return request status GG_REQUEST_SUCCESS or
	// deliver to no targets and return a request status GG_REQUEST_AGAIN
	QueueFullPolicyOptionAllOrError = 1
	// QueueFullPolicyOptionReservedMax max in enum
	QueueFullPolicyOptionReservedMax = 2
	// QueueFullPolicyOptionReservedPad is reserved.
	QueueFullPolicyOptionReservedPad = 0x7FFFFFFF
)

// QueueAPI handles publishing to MQTT through the local API
type QueueAPI struct {
	APIRequest
}

// NewQueueAPI creates a new MQTT client
func NewQueueAPI() *QueueAPI {
	return &QueueAPI{}
}

// Publish will publish a payload on provided topic.
func (q *QueueAPI) Publish(topic string, option QueueFullPolicyOption, payload []byte) {

	if q.err != nil {
		return
	}

	tn := C.CString(topic)

	q.initialize()

	defer func() {
		q.close()
		C.free(unsafe.Pointer(tn))
	}()

	var opts C.gg_publish_options
	e := GreenGrassCode(C.gg_publish_options_init(&opts))

	q.err = createError("Failed to init publish options", e)

	if q.err != nil {
		return
	}

	e = GreenGrassCode(
		C.gg_publish_options_set_queue_full_policy(
			opts,
			C.gg_queue_full_policy_options(option),
		))

	q.err = createError("Failed to set publish options", e)

	if q.err != nil {
		return
	}

	res := C.gg_request_result{}
	e = GreenGrassCode(
		C.gg_publish_with_options(
			q.request, tn,
			unsafe.Pointer(&payload[0]), C.size_t(len(payload)),
			opts, &res,
		),
	)

	q.handleRequestResponse(
		fmt.Sprintf("Failed to post to topic '%s'", topic),
		e, RequestStatus(res.request_status),
	)

}

// PublishObject will publish the object, marshalled as
// _JSON_, onto specified topic.
func (q *QueueAPI) PublishObject(
	topic string,
	policy QueueFullPolicyOption,
	object interface{}) *QueueAPI {

	if nil == object {
		return q
	}

	data, err := json.Marshal(object)
	if nil == err {
		q.Publish(topic, policy, data)
	} else {
		q.err = fmt.Errorf("Failed to marshal object, error: %s", err.Error())
	}

	return q
}
