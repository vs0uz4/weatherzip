package mock

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMockMemoryService(t *testing.T) {
	mock := &MockMemoryService{
		GetMemoryStatsFunc: func() (uint64, uint64, uint64, uint64, float64, error) {
			return 1024, 512, 256, 128, 50.5, nil
		},
	}

	total, used, free, available, percentUsed, err := mock.GetMemoryStats()

	assert.NoError(t, err, "Expected no error from mock")
	assert.Equal(t, uint64(1024), total, "Expected total to match mock value")
	assert.Equal(t, uint64(512), used, "Expected used to match mock value")
	assert.Equal(t, uint64(256), free, "Expected free to match mock value")
	assert.Equal(t, uint64(128), available, "Expected available to match mock value")
	assert.Equal(t, 50.5, percentUsed, "Expected percentUsed to match mock value")

	mock.GetMemoryStatsFunc = func() (uint64, uint64, uint64, uint64, float64, error) {
		return 0, 0, 0, 0, 0, errors.New("mock error")
	}

	_, _, _, _, _, err = mock.GetMemoryStats()
	assert.Error(t, err, "Expected error from mock")
}
