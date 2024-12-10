package mock

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMockCPUService(t *testing.T) {
	mock := &MockCPUService{
		GetCPUStatsFunc: func() (int, []float64, error) {
			return 4, []float64{10.5, 20.3}, nil
		},
	}

	cores, percentUsed, err := mock.GetCPUStats()

	assert.NoError(t, err, "Expected no error from mock")
	assert.Equal(t, 4, cores, "Expected cores to match mock value")
	assert.Equal(t, []float64{10.5, 20.3}, percentUsed, "Expected percentUsed to match mock value")

	mock.GetCPUStatsFunc = func() (int, []float64, error) {
		return 0, nil, errors.New("mock error")
	}

	_, _, err = mock.GetCPUStats()
	assert.Error(t, err, "Expected error from mock")
}
