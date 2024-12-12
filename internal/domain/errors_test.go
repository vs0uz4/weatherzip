package domain

import (
	"errors"
	"fmt"
	"testing"
)

func TestNewUnexpectedStatusCodeError(t *testing.T) {
	statusCode := 500
	expectedMessage := fmt.Sprintf("unexpected status code: %d", statusCode)
	err := NewUnexpectedStatusCodeError(statusCode)

	if err.Error() != expectedMessage {
		t.Errorf("Expected error message %q, got %q", expectedMessage, err.Error())
	}
}

func TestNewFailedToCreateRequestError(t *testing.T) {
	originalErr := errors.New("invalid URL")
	expectedMessage := fmt.Sprintf("failed to create request: %v", originalErr)
	err := NewFailedToCreateRequestError(originalErr)

	if err.Error() != expectedMessage {
		t.Errorf("Expected error message %q, got %q", expectedMessage, err.Error())
	}
}

func TestNewFailedToMakeRequestError(t *testing.T) {
	originalErr := errors.New("connection refused")
	expectedMessage := fmt.Sprintf("failed to make request: %v", originalErr)
	err := NewFailedToMakeRequestError(originalErr)

	if err.Error() != expectedMessage {
		t.Errorf("Expected error message %q, got %q", expectedMessage, err.Error())
	}
}

func TestNewFailedToDecodeResponseError(t *testing.T) {
	originalErr := errors.New("invalid JSON format")
	expectedMessage := fmt.Sprintf("failed to decode response: %v", originalErr)
	err := NewFailedToDecodeResponseError(originalErr)

	if err.Error() != expectedMessage {
		t.Errorf("Expected error message %q, got %q", expectedMessage, err.Error())
	}
}
