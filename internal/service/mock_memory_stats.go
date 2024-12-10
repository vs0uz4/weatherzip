package service

type MockMemoryService struct {
	GetMemoryStatsFunc func() (uint64, uint64, uint64, uint64, float64, error)
}

func (m *MockMemoryService) GetMemoryStats() (uint64, uint64, uint64, uint64, float64, error) {
	return m.GetMemoryStatsFunc()
}
