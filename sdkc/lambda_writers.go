package sdkc

// ErrorResponseWriter writes a error response failure
type ErrorResponseWriter struct {
}

// NewErrorResponseWriter creates a new error response writer.
func NewErrorResponseWriter() *ErrorResponseWriter {
	return &ErrorResponseWriter{}
}

func (erw *ErrorResponseWriter) Write(p []byte) (n int, err error) {

	if err = lambdaWriteError(string(p)); err == nil {
		n = len(p)
	}

	return
}

// ResponseWriter is writes a success response payload
type ResponseWriter struct {
}

// NewResponseWriter creates a new response writer.
func NewResponseWriter() *ResponseWriter {
	return &ResponseWriter{}
}

func (rw *ResponseWriter) Write(p []byte) (n int, err error) {

	if err = lambdaWriteResponse(p); err == nil {
		n = len(p)
	}

	return
}
