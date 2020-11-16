from ctypes import *
from glue import lib
import json

@CFUNCTYPE(None, c_char_p, c_char_p, c_char_p, c_char_p)
def getSecret(ctx: str, secretId: str, versionId: str, versionStage: str):
    print("secret id: " + secretId.decode("utf-8"))    
    print("version: " + versionId.decode("utf-8"))
    print("stage: " + versionStage.decode("utf-8"))


    sample = json.dumps({
        'arn': 'arn:zyx',
        'name':secretId.decode("utf-8"),
        'version': versionId.decode("utf-8"),
        'bin': 'eyJ1c2VyIjoibmlzc2UifQ==',
        'secret': 'my secret',
        'stages': ['1','2'],
        'created':'2020-18-02T18:34:00'
    })

    lib.set_process_buffer(ctx, sample.encode("utf-8"))

lib.initcb(
    None,
    None,
    None,
    None,
    getSecret
)

lib.setup()
