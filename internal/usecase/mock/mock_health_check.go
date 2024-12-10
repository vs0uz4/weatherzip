package mock

import "weatherzip/internal/infra/web/health"

type MockHealthCheckUseCase struct {
	GetHealthFunc func() (health.HealthStats, error)
}

func (m *MockHealthCheckUseCase) GetHealth() (health.HealthStats, error) {
	return m.GetHealthFunc()
}
