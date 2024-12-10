package mock

import (
	"errors"
	"testing"
	"weatherzip/internal/infra/web/health"

	"github.com/stretchr/testify/assert"
)

func TestMockHealthCheckUseCase(t *testing.T) {
	mock := &MockHealthCheckUseCase{
		GetHealthFunc: func() (health.HealthStats, error) {
			return health.HealthStats{Status: "pass"}, nil
		},
	}

	healthStats, err := mock.GetHealth()

	assert.NoError(t, err, "Expected no error from mock")
	assert.Equal(t, "pass", healthStats.Status, "Expected status to match mock value")

	mock.GetHealthFunc = func() (health.HealthStats, error) {
		return health.HealthStats{}, errors.New("mock error")
	}

	_, err = mock.GetHealth()
	assert.Error(t, err, "Expected error from mock")
}
