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

// GreenGrassLogLevel specifies the verbosity of the log output.
type GreenGrassLogLevel int

const (
	// GreenGrassLogLevelReservedNotSet is not used.
	GreenGrassLogLevelReservedNotSet = 0
	//GreenGrassLogLevelDebug specifies debug output
	GreenGrassLogLevelDebug = 1
	//GreenGrassLogLevelInfo specifies info output
	GreenGrassLogLevelInfo = 2
	// GreenGrassLogLevelWarn specifies warn output
	GreenGrassLogLevelWarn = 3
	//GreenGrassLogLevelError will error output
	GreenGrassLogLevelError = 4
	// GreenGrassLogLevelFatal is fatal. System will exist
	GreenGrassLogLevelFatal = 5
	//GreenGrassLogLevelReservedMax is last enum
	GreenGrassLogLevelReservedMax = 6
	//GreenGrassLogLevelReservedPad for padding
	GreenGrassLogLevelReservedPad = 0x7FFFFFFF
)

// Log logs out using a printf format.
func Log(level GreenGrassLogLevel, format string, args ...interface{}) {

	f := C.CString(fmt.Sprintf(format, args...))
	C.logwrapper(C.int(level), f)

	C.free(unsafe.Pointer(f))
}
