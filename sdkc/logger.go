package sdkc

/*
#include<stdlib.h>
#include "lib/greengrasssdk.h"

extern void logwrapper(int level, char *log) {
	gg_log(level, log);
}
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// LogLevel specifies the verbosity of the log output.
type LogLevel int

const (
	// LogLevelReservedNotSet is not used.
	LogLevelReservedNotSet = 0
	//LogLevelDebug specifies debug output
	LogLevelDebug = 1
	//LogLevelInfo specifies info output
	LogLevelInfo = 2
	// LogLevelWarn specifies warn output
	LogLevelWarn = 3
	//LogLevelError will error output
	LogLevelError = 4
	// LogLevelFatal is fatal. System will exist
	LogLevelFatal = 5
	//LogLevelReservedMax is last enum
	LogLevelReservedMax = 6
	//LogLevelReservedPad for padding
	LogLevelReservedPad = 0x7FFFFFFF
)

// Log logs out using a printf format.
func Log(level LogLevel, format string, args ...interface{}) {

	f := C.CString(fmt.Sprintf(format, args...))
	C.logwrapper(C.int(level), f)

	C.free(unsafe.Pointer(f))
}
