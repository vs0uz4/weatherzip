package service

import (
	"testing"

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
