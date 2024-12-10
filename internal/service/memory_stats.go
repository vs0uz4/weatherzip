package service

import (
	"math"

	"github.com/shirou/gopsutil/mem"
)

type MemoryService interface {
	GetMemoryStats() (uint64, uint64, uint64, uint64, float64, error)
}

type memoryService struct{}

func NewMemoryService() MemoryService {
	return &memoryService{}
}

func (s *memoryService) roundToOneDecimal(value float64) float64 {
	return math.Round(value*10) / 10
}

func (s *memoryService) GetMemoryStats() (uint64, uint64, uint64, uint64, float64, error) {
	memStats, err := mem.VirtualMemory()
	if err != nil {
		return 0, 0, 0, 0, 0, err
	}
	percentUsed := s.roundToOneDecimal(memStats.UsedPercent)
	return memStats.Total, memStats.Used, memStats.Free, memStats.Available, percentUsed, nil
}
