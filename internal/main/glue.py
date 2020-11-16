from ctypes import c_char_p,c_longlong, cdll, Structure

class GoString(Structure):
    _fields_ = [("p", c_char_p), ("n", c_longlong)]

lib = cdll.LoadLibrary("./main.so")

lib.set_process_buffer.argtypes = [c_char_p, c_char_p]