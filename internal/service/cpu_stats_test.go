package service

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetCPUStats(t *testing.T) {
	cpuService := NewCPUService()

	cores, percentUsed, err := cpuService.GetCPUStats()

	assert.NoError(t, err, "Getting CPU stats should not produce an error")
	assert.Greater(t, cores, 0, "Cores should be greater than 0")
	assert.NotEmpty(t, percentUsed, "Percent used should not be empty")

	for _, usage := range percentUsed {
		assert.GreaterOrEqual(t, usage, 0.0, "Usage percentage should be >= 0")
		assert.LessOrEqual(t, usage, 100.0, "Usage percentage should be <= 100")
	}
}

func TestGetCPUStatsError(t *testing.T) {
	mockCPUService := &cpuService{
		cpuPercentFunc: func(interval time.Duration, percpu bool) ([]float64, error) {
			return nil, errors.New("mock error")
		},
	}

	cores, percentUsed, err := mockCPUService.GetCPUStats()

	assert.Error(t, err, "An error should occur")
	assert.Equal(t, 0, cores, "Cores should be 0 on error")
	assert.Nil(t, percentUsed, "Percent used should be nil on error")
}
