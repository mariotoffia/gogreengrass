from ctypes import *
import json
from contextmock import _Context

# define class GoString to map:
# C type struct { const char *p; GoInt n; }
class GoString(Structure):
    _fields_ = [("p", c_char_p), ("n", c_longlong)]

lib = cdll.LoadLibrary("./main.so")
lib.invokeJSON.restype = c_char_p

# Main code for lambda
lib.setup()

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




c = _Context('arn:aws:lambda:eu-central-1:033549287452:function:ggtest:7','99994f79-2e36-4eb2-6584-b533cb3bc491')
result = invokeJSON(c, { "data": 19, "hello": "world" })

print("result: " + result)