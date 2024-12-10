package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMemoryStats(t *testing.T) {
	memoryService := NewMemoryService()

	total, used, free, available, percentUsed, err := memoryService.GetMemoryStats()

	assert.NoError(t, err, "Getting memory stats should not produce an error")
	assert.Greater(t, total, uint64(0), "Total memory should be greater than 0")
	assert.GreaterOrEqual(t, used, uint64(0), "Used memory should be >= 0")
	assert.GreaterOrEqual(t, free, uint64(0), "Free memory should be >= 0")
	assert.GreaterOrEqual(t, available, uint64(0), "Available memory should be >= 0")

	assert.GreaterOrEqual(t, percentUsed, 0.0, "Percent used should be >= 0")
	assert.LessOrEqual(t, percentUsed, 100.0, "Percent used should be <= 100")
}
