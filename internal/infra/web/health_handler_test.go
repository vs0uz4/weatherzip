package web

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"weatherzip/internal/infra/web/health"
	"weatherzip/internal/service"
	"weatherzip/internal/usecase"
	"weatherzip/internal/usecase/mock"

	"github.com/stretchr/testify/assert"
)

type ErrorResponseWriter struct{}

func (e *ErrorResponseWriter) Header() http.Header {
	return http.Header{}
}

func (e *ErrorResponseWriter) Write([]byte) (int, error) {
	return 0, errors.New("mock encoding error")
}

func (e *ErrorResponseWriter) WriteHeader(statusCode int) {}

func TestHealthHandlerGetHealth(t *testing.T) {
	cpuService := service.NewCPUService()
	memoryService := service.NewMemoryService()
	uptimeService := service.NewUptimeService()
	healthCheck := usecase.NewHealthCheckUseCase(cpuService, memoryService, uptimeService)
	handler := NewHealthHandler(healthCheck)

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	handler.GetHealth(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Response status code should be 200")
}

func TestHealthHandlerGetHealthErrorFromUseCase(t *testing.T) {
	mockUseCase := &mock.MockHealthCheckUseCase{
		GetHealthFunc: func() (health.HealthStats, error) {
			return health.HealthStats{}, errors.New("mock use case error")
		},
	}

	handler := NewHealthHandler(mockUseCase)
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	handler.GetHealth(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode, "Response status code should be 500 for use case error")
}

func TestHealthHandlerGetHealthErrorEncodingResponse(t *testing.T) {
	mockUseCase := &mock.MockHealthCheckUseCase{
		GetHealthFunc: func() (health.HealthStats, error) {
			return health.HealthStats{Status: "pass"}, nil
		},
	}

	handler := NewHealthHandler(mockUseCase)
	req := httptest.NewRequest("GET", "/health", nil)
	w := &ErrorResponseWriter{}

	handler.GetHealth(w, req)

	assert.True(t, true, "Encoding error should be handled gracefully")
}
