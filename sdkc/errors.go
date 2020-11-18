package sdkc

// GreenGrassCode is the error codes that the runtime returns
type GreenGrassCode int

/**
 * @brief Greengrass SDK error enum
 *
 * Enumeration of return values from the gg_* functions within the SDK.
 */
const (
	// GreenGrassCodeSuccess returned when success.
	GreenGrassCodeSuccess GreenGrassCode = 0
	// GreenGrassCodeOutOfMemory is returned when process is out of memory
	GreenGrassCodeOutOfMemory GreenGrassCode = 1
	// GreenGrassCodeInvalidParameter is returned when input parameter is invalid
	GreenGrassCodeInvalidParameter GreenGrassCode = 2
	// GreenGrassCodeInvalidState is returned when SDK is in an invalid state
	GreenGrassCodeInvalidState GreenGrassCode = 3
	// GreenGrassCodeInternalFailure is returned when SDK encounters internal failure
	GreenGrassCodeInternalFailure GreenGrassCode = 4
	// GreenGrassCode is returned when process gets signal to terminate
	GreenGrassCodeTerminate   GreenGrassCode = 5
	GreenGrassCodeReservedMax GreenGrassCode = 6
	GreenGrassCodeReservedPad GreenGrassCode = 0x7FFFFFFF
)
