package service

type MockCPUService struct {
	GetCPUStatsFunc func() (int, []float64, error)
}

func (m *MockCPUService) GetCPUStats() (int, []float64, error) {
	return m.GetCPUStatsFunc()
}
