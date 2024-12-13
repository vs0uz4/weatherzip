package middleware

import (
	"testing"
)

func TestResponseRecorderWriteError(t *testing.T) {
	rr := &ResponseRecorder{}
	errorMessage := "test error message"

	rr.WriteError(errorMessage)

	if rr.errorMessage != errorMessage {
		t.Errorf("Expected errorMessage to be %q, got %q", errorMessage, rr.errorMessage)
	}
}
