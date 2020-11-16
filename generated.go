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

