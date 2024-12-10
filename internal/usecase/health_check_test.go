package usecase

import (
	"testing"
	"weatherzip/internal/service"

	"github.com/stretchr/testify/assert"
)

func TestGetHealth(t *testing.T) {
	cpuService := service.NewCPUService()
	memoryService := service.NewMemoryService()
	uptimeService := service.NewUptimeService()

	healthCheck := NewHealthCheckUseCase(cpuService, memoryService, uptimeService)

	healthStats, err := healthCheck.GetHealth()

	assert.NoError(t, err, "HealthCheck should not produce an error")
	assert.Equal(t, "pass", healthStats.Status, "Default status should be 'pass'")
	assert.NotEmpty(t, healthStats.CPU, "CPU stats should not be empty")
	assert.NotEmpty(t, healthStats.Memory, "Memory stats should not be empty")
	assert.NotEmpty(t, healthStats.Uptime, "Uptime should not be empty")
	assert.NotEmpty(t, healthStats.Message, "Message should not be empty")
}
