package web

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"weatherzip/internal/service"
	"weatherzip/internal/usecase"

	"github.com/stretchr/testify/assert"
)

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
