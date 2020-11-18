package main

var glueGo = []byte(`package main

/*
typedef void (*publish) (char* topic, char *policy, char *payload);
typedef void (*getThingShadow) (char* thingName);
typedef void (*updateThingShadow) (char *ctx, char* thingName, char *payload);
typedef void (*deleteThingShadow) (char *ctx, char* thingName);
typedef void (*get_secret) (char *ctx, char* id, char *version, char *stage);
*/
import "C"

import (
	"github.com/mariotoffia/gogreengrass/sdk"
)

// GGDataplane is the functions to for dataplane against GGC API
var GGDataplane sdk.DataplaneClient

// GGSecretsManager is the functions for retrieving secrets via GGC API
var GGSecretsManager sdk.SecretsManager

//export initcb
func initcb(
	fnPublish C.publish,
	fnGetThingShadow C.getThingShadow,
	fnUpdateThingShadow C.updateThingShadow,
	fnDeleteThingShadow C.deleteThingShadow,
	fnGetSecret C.get_secret) {

	GGDataplane = sdk.NewDataplaneClient(
		fnPublish,
		fnGetThingShadow,
		fnUpdateThingShadow,
		fnDeleteThingShadow)

	GGSecretsManager = sdk.NewSecretsManagerClient(fnGetSecret)

}

//export setup
func setup() {
	once()
}

//export set_process_buffer
func set_process_buffer(ctx *C.char, buffer *C.char) {
	sdk.SetProcessBuffer(C.GoString(ctx), C.GoString(buffer))
}

//export invokeJSON
func invokeJSON(context string, payload string) *C.char {
	return C.CString(sdk.InvokeJSON(context, payload))
}
`)

var gluePy = []byte(`from ctypes import *
import json

# you need to copy sdk to folder greengrasssdk
import greengrasssdk

class GoString(Structure):
    _fields_ = [("p", c_char_p), ("n", c_longlong)]

# Creating a greengrass core sdk client
client = greengrasssdk.client("iot-data")
sm = greengrasssdk.client("secretsmanager")

# Load your library that you've built with 
# go build -o main.so -buildmode=c-shared main.go init.go
lib = cdll.LoadLibrary("./main.so")
lib.invokeJSON.restype = c_char_p

# The actual function that you need to bind the lambda entry point to
def function_handler(event, context):
    result = invokeJSON(context, event)
    return result


@CFUNCTYPE(None, c_char_p, c_char_p, c_char_p)
def publishcb(topic: str, queueFullPolicy: str, payload: str):
     client.publish(topic=topic.decode("utf-8"),
                   queueFullPolicy=queueFullPolicy.decode("utf-8"),
                   payload=payload.decode("utf-8"))

@CFUNCTYPE(None, c_char_p, c_char_p)
def getThingShadow(ctx: str, thingName: str):
    result = client.get_thing_shadow(thingName=thingName).payload
    lib.set_process_buffer(ctx, result)

@CFUNCTYPE(None, c_char_p, c_char_p, c_char_p)
def updateThingShadow(ctx: str, thingName: str, payload: str):
    result = client.update_thing_shadow(
        thingName=thingName,
        payload=payload).payload

    lib.set_process_buffer(ctx, result)

@CFUNCTYPE(None, c_char_p, c_char_p)
def deleteThingShadow(ctx: str, thingName: str):
    result = client.delete_thing_shadow(thingName=thingName).payload
    lib.set_process_buffer(ctx, result)

@CFUNCTYPE(None, c_char_p, c_char_p, c_char_p, c_char_p)
def getSecret(ctx: str, secretId: str, versionId: str, versionStage: str):
    ret  = sm.get_secret_value(SecretId=secretId, VersionId=versionId, VersionStage=versionStage)

    o = json.dumps({
        'arn': ret.ARN,
        'name':ret.Name,
        'version': ret.VersionId,
        'bin': ret.SecretBinary.encode("utf-8"),
        'secret': ret.SecretString,
        'stages': ret.VersionStages,
        'created':'2020-18-02T18:34:00'
    }) 
    ## todo: - created: ret.CreatedDate (datetime)

    lib.set_process_buffer(ctx, o.encode("utf-8"))

lib.initcb(
    publishcb,
    getThingShadow,
    updateThingShadow,
    deleteThingShadow,
    getSecret
)

# Initialize the lambda for invocation (one-time only)
lib.setup()

# Invokes the golang lambda handler (non binary version).
def invokeJSON(context: any,
               event: any,
               deadlineMS: str = '300000') -> str:

    c = json.dumps({
        'aws_request_id': context.aws_request_id,
        'client_context': {
            'client': context.client_context.client,
            'custom': context.client_context.custom,
            'env': context.client_context.env
        },
        'function_name': context.function_name,
        'function_version': context.function_version,
        'identity': context.identity,
        'invoked_function_arn': context.invoked_function_arn,
        'headers': {
            'Lambda-Runtime-Deadline-Ms': deadlineMS
        }
    }).encode("utf-8")

    e = json.dumps(event).encode("utf-8")

    goContext = GoString(c, len(c))
    goEvent = GoString(e, len(e))

    r = lib.invokeJSON(goContext, goEvent)
    return r.decode("utf-8")
`)

var soFile = []byte(`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`)

