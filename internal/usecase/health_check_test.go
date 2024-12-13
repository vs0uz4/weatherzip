package usecase

import (
	"errors"
	"testing"

	"github.com/vs0uz4/weatherzip/internal/service"
	"github.com/vs0uz4/weatherzip/internal/service/mock"

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

func TestGetHealthCpuServiceError(t *testing.T) {
	mockCPUService := &mock.MockCPUService{
		GetCPUStatsFunc: func() (int, []float64, error) {
			return 0, nil, errors.New("mock CPU error")
		},
	}
	mockMemoryService := service.NewMemoryService()
	mockUptimeService := service.NewUptimeService()

	healthCheck := NewHealthCheckUseCase(mockCPUService, mockMemoryService, mockUptimeService)

	healthStats, err := healthCheck.GetHealth()

	assert.NoError(t, err, "HealthCheck should not return an error")
	assert.Equal(t, "fail", healthStats.Status, "Status should be 'fail' when CPU service fails")
	assert.Equal(t, "Still alive, but not kicking!", healthStats.Message, "Message should reflect the 'fail' status")
}

func TestGetHealthMemoryServiceError(t *testing.T) {
	mockCPUService := service.NewCPUService()
	mockMemoryService := &mock.MockMemoryService{
		GetMemoryStatsFunc: func() (uint64, uint64, uint64, uint64, float64, error) {
			return 0, 0, 0, 0, 0, errors.New("mock memory error")
		},
	}
	mockUptimeService := service.NewUptimeService()

	healthCheck := NewHealthCheckUseCase(mockCPUService, mockMemoryService, mockUptimeService)

	healthStats, err := healthCheck.GetHealth()

	assert.NoError(t, err, "HealthCheck should not return an error")
	assert.Equal(t, "fail", healthStats.Status, "Status should be 'fail' when memory service fails")
	assert.Equal(t, "Still alive, but not kicking!", healthStats.Message, "Message should reflect the 'fail' status")
}

func TestGetHealthMessageWhenNotPass(t *testing.T) {
	mockCPUService := &mock.MockCPUService{
		GetCPUStatsFunc: func() (int, []float64, error) {
			return 0, nil, errors.New("mock CPU error")
		},
	}
	mockMemoryService := &mock.MockMemoryService{
		GetMemoryStatsFunc: func() (uint64, uint64, uint64, uint64, float64, error) {
			return 0, 0, 0, 0, 0, errors.New("mock memory error")
		},
	}
	mockUptimeService := service.NewUptimeService()

	healthCheck := NewHealthCheckUseCase(mockCPUService, mockMemoryService, mockUptimeService)

	healthStats, err := healthCheck.GetHealth()

	assert.NoError(t, err, "HealthCheck should not return an error")
	assert.Equal(t, "fail", healthStats.Status, "Status should be 'fail' when both services fail")
	assert.Equal(t, "Still alive, but not kicking!", healthStats.Message, "Message should be 'Still alive, but not kicking!' when status is 'fail'")
}
