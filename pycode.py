from ctypes import *
import json

# context = _Context(self.function_arn, invocation_id, client_context)
# class _Context:
#     def __init__(self, function_arn, invocation_id, client_context_encoded=None):
#         arn_fields = FunctionArnFields(function_arn)

#         self.invoked_function_arn = function_arn
#         self.function_name = arn_fields.name
#         self.function_version = arn_fields.qualifier
#         self.aws_request_id = invocation_id

#         if client_context_encoded:
#             client_context_map = json.loads(base64.b64decode(client_context_encoded).decode('utf-8'))
#             client_context_client = None
#             if 'client' in client_context_map:
#                 client_context_client = ClientContextClient(**client_context_map['client'])
#             client_context_custom = None
#             if 'custom' in client_context_map:
#                 client_context_custom = client_context_map['custom']
#             client_context_env = None
#             if 'env' in client_context_map:
#                 client_context_env = client_context_map['env']

#             self.client_context = ClientContext(
#                 client_context_client,
#                 client_context_custom,
#                 client_context_env
#             )
#         else:
#             self.client_context = None

#         self.identity = None

        ## skip self.memory_limit_in_mb
        ## skip self.log_group_name
        ## skip self.log_stream_name

# define class GoString to map:
# C type struct { const char *p; GoInt n; }
class GoString(Structure):
    _fields_ = [("p", c_char_p), ("n", c_longlong)]

lib = cdll.LoadLibrary("./main.so")

# Main code for lambda
lib.setup()

# Per invocation
context = json.dumps({
	'aws_request_id': '99994f79-2e36-4eb2-6584-b533cb3bc491',
	'client_context': {
        'client': None, 
        'custom': {'subject': 'invoke/ggtest'}, 
        'env': None,
    },
	'function_name': 'ggtest',
	'function_version': '7',
	'identity': None,
	'invoked_function_arn': 'arn:aws:lambda:eu-central-1:033549287452:function:ggtest:7',
    'headers': {
        'Lambda-Runtime-Deadline-Ms': '300000'
    }
}).encode("utf-8")

event = json.dumps({
    "data": 19,
    "hello": "world"
}).encode("utf-8")


cg = GoString(context, len(context))
pg = GoString(event, len(event))
lib.invokeJSON(cg, pg)