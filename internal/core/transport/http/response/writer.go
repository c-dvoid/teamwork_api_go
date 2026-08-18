package core_http_response

import "net/http"

var (
	StatusCodeUninitialized = -1
)

type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
		statusCode:     StatusCodeUninitialized,
	}
}

func (rw *ResponseWriter) WriteHeader(statusCode int) {
	rw.ResponseWriter.WriteHeader(statusCode)
	rw.statusCode = statusCode
}

// т.к Write перехватывает случаи, когда WriteHeader не был вызван явно
// (например, сторонними хендлерами вроде httpSwagger.WrapHandler),
// повторяется стандартное поведение net/http, тем самым первым вызывая Write() без
// предварительного WriteHeader, что приводит к неявному статусу 200 OK.
func (rw *ResponseWriter) Write(b []byte) (int, error) {
	if rw.statusCode == StatusCodeUninitialized {
		rw.WriteHeader(http.StatusOK)
	}
	return rw.ResponseWriter.Write(b)
}

func (rw *ResponseWriter) GetStatusCode() int {
	if rw.statusCode == StatusCodeUninitialized {
		panic("no status code set")
	}
	return rw.statusCode
}
