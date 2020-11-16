from glue import GoString, lib
from ctypes import c_char_p
import json
from contextmock import _Context

def function_handler(event, context):
    result = invokeJSON(context, event)
    print("result: " + result)

# Main code for lambda
lib.invokeJSON.restype = c_char_p
lib.setup()

def invokeJSON(context,event, deadlineMS: str = '300000') -> str:

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
