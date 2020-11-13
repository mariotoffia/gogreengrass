from ctypes import *
import json

from contextmock import _Context

def function_handler(event, context):
    result = invokeJSON(context, event)
    print("result: " + result)

class GoString(Structure):
    _fields_ = [("p", c_char_p), ("n", c_longlong)]

lib = cdll.LoadLibrary("./main.so")
lib.invokeJSON.restype = c_char_p


# Main code for lambda
lib.setup()

def publishcb(topic: str, queueFullPolicy: str, payload: str):
    print("publishing")
    print("----------")
    print(topic.decode("utf-8"))
    print(queueFullPolicy.decode("utf-8"))
    print(payload.decode("utf-8"))

CMPFUNC = CFUNCTYPE(None, c_char_p, c_char_p, c_char_p)
callback_publish = CMPFUNC(publishcb)
lib.initcb(callback_publish)

def invokeJSON(context: _Context,
               event: any,
               deadlineMS: str = '300000') -> str:

    cDump = json.dumps({
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

    eDump = json.dumps(event).encode("utf-8")
    cg = GoString(cDump, len(cDump))
    pg = GoString(eDump, len(eDump))

    r = lib.invokeJSON(cg, pg)
    return r.decode("utf-8")


c = _Context(
             'arn:aws:lambda:eu-central-1:033549287452:function:ggtest:7',
             '99994f79-2e36-4eb2-6584-b533cb3bc491'
             )
             
e = {"data": 44, "hello": "world"}

function_handler(e, c)
