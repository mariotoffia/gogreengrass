from ctypes import *
from glue import lib

# sudo apt-get install python3-dev
# https://www.datadoghq.com/blog/engineering/cgo-and-python/
# https://stackoverflow.com/questions/52967133/python3-ctypes-callback-causes-memory-leak-in-simple-example
# https://stackoverflow.com/questions/63442920/how-to-avoid-runtimewarning-memory-leak-in-callback-function-with-python-callba


@CFUNCTYPE(None, c_char_p, c_char_p, c_char_p)
def publishcb(topic: str, queueFullPolicy: str, payload: str):
    print("publishing")
    print("----------")
    print(topic.decode("utf-8"))
    print(queueFullPolicy.decode("utf-8"))
    print(payload.decode("utf-8"))

@CFUNCTYPE(None, c_char_p, c_char_p)
def getThingShadow(ctx: str, thingName: str):
    print("thing name: " + thingName.decode("utf-8"))    
    lib.set_process_buffer(ctx, "device shadow".encode("utf-8"))

@CFUNCTYPE(None, c_char_p, c_char_p, c_char_p)
def updateThingShadow(ctx: str, thingName: str, payload: str):
    print("thing name: " + thingName.decode("utf-8"))    
    print("shadow: " + payload.decode("utf-8"))
    lib.set_process_buffer(ctx, "updated device shadow".encode("utf-8"))

@CFUNCTYPE(None, c_char_p, c_char_p)
def deleteThingShadow(ctx: str, thingName: str):
    print("thing name: " + thingName.decode("utf-8"))    
    lib.set_process_buffer(ctx, "device shadow deleted".encode("utf-8"))

lib.initcb(
    publishcb,
    getThingShadow,
    updateThingShadow,
    deleteThingShadow,
    None
)

lib.setup()
