from ctypes import *
import json

# you need to copy sdk to folder greengrasssdk
import greengrasssdk

class GoString(Structure):
    _fields_ = [("p", c_char_p), ("n", c_longlong)]

# Load your library that you've built with 
# go build -o main.so -buildmode=c-shared main.go init.go
lib = cdll.LoadLibrary("./main.so")
lib.invokeJSON.restype = c_char_p

# Initialize the lambda for invocation (one-time only)
lib.setup()

# Creating a greengrass core sdk client
client = greengrasssdk.client("iot-data")

# The actual function that you need to bind the lambda entry point to
def function_handler(event, context):
    result = invokeJSON(context, event)
    return result


# publish may be called from go lambda to publish data
def publish(topic: str, queueFullPolicy: str, payload: str):
    client.publish(topic=topic,
                   queueFullPolicy=queueFullPolicy,
                   payload=payload)

# This is invoked by the python lambda function_handler
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
