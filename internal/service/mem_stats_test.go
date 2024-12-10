package service

import (
	"errors"
	"testing"

	"github.com/shirou/gopsutil/mem"
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

func TestGetMemoryStatsError(t *testing.T) {
	mockMemoryService := &memoryService{
		virtualMemoryFunc: func() (*mem.VirtualMemoryStat, error) {
			return nil, errors.New("mock error")
		},
	}

	total, used, free, available, percentUsed, err := mockMemoryService.GetMemoryStats()

	assert.Error(t, err, "An error should occur")
	assert.Equal(t, uint64(0), total, "Total memory should be 0 on error")
	assert.Equal(t, uint64(0), used, "Used memory should be 0 on error")
	assert.Equal(t, uint64(0), free, "Free memory should be 0 on error")
	assert.Equal(t, uint64(0), available, "Available memory should be 0 on error")
	assert.Equal(t, 0.0, percentUsed, "Percent used should be 0.0 on error")
}
