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

