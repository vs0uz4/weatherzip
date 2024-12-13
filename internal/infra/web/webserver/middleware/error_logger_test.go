package middleware

import (
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestResponseRecorderWriteHeader(t *testing.T) {
	mockWriter := httptest.NewRecorder()
	rr := &ResponseRecorder{ResponseWriter: mockWriter}

	statusCode := http.StatusNotFound
	rr.WriteHeader(statusCode)

	if rr.statusCode != statusCode {
		t.Errorf("Expected statusCode %d, got %d", statusCode, rr.statusCode)
	}

	if mockWriter.Code != statusCode {
		t.Errorf("Expected underlying writer statusCode %d, got %d", statusCode, mockWriter.Code)
	}
}

func TestResponseRecorderWrite(t *testing.T) {
	mockWriter := httptest.NewRecorder()
	rr := &ResponseRecorder{ResponseWriter: mockWriter}

	data := []byte("test data")
	written, err := rr.Write(data)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if written != len(data) {
		t.Errorf("Expected %d bytes written, got %d", len(data), written)
	}

	if rr.bytesWritten != len(data) {
		t.Errorf("Expected bytesWritten %d, got %d", len(data), rr.bytesWritten)
	}

	if mockWriter.Body.String() != string(data) {
		t.Errorf("Expected body %q, got %q", string(data), mockWriter.Body.String())
	}
}

func TestResponseRecorderWriteError(t *testing.T) {
	rr := &ResponseRecorder{}
	errorMessage := "test error message"

	rr.WriteError(errorMessage)

	if rr.ReadError() != errorMessage {
		t.Errorf("Expected errorMessage to be %q, got %q", errorMessage, rr.ReadError())
	}
}

func TestErrorLogger(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		expectedLogged bool
	}{
		{"Logs Error for 404", http.StatusNotFound, true},
		{"Logs Error for 500", http.StatusInternalServerError, true},
		{"Does Not Log for 200", http.StatusOK, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockWriter := httptest.NewRecorder()
			mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				w.Write([]byte("test"))
			})

			logger := ErrorLogger(mockHandler)

			var logs strings.Builder
			log.SetOutput(&logs)
			defer log.SetOutput(nil)

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			logger.ServeHTTP(mockWriter, req)

			logged := strings.Contains(logs.String(), "Error:")
			if logged != tt.expectedLogged {
				t.Errorf("Expected log presence %v, got %v", tt.expectedLogged, logged)
			}
		})
	}
}
