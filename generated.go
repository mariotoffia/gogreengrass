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

