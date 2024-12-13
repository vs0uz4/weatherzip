package middleware

import (
	"log"
	"net/http"
	"time"
)

type ResponseRecorder struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int
	errorMessage string
}

func (rr *ResponseRecorder) WriteHeader(code int) {
	rr.statusCode = code

	if code != http.StatusOK && rr.bytesWritten == 0 {
		rr.ResponseWriter.WriteHeader(code)
	}
}

func (rr *ResponseRecorder) Write(data []byte) (int, error) {
	size, err := rr.ResponseWriter.Write(data)
	rr.bytesWritten += size
	return size, err
}

func (rr *ResponseRecorder) ReadError() string {
	return rr.errorMessage
}

func (rr *ResponseRecorder) WriteError(message string) {
	rr.errorMessage = message
}

func ErrorLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rr := &ResponseRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(rr, r)

		if rr.statusCode >= 400 {
			duration := time.Since(start)
			log.Printf(`"%s %s %s" from %s - %d %dB in %v - Error: %s`,
				r.Method,
				r.URL.String(),
				r.Proto,
				r.RemoteAddr,
				rr.statusCode,
				rr.bytesWritten,
				duration,
				rr.errorMessage,
			)
		}
	})
}
