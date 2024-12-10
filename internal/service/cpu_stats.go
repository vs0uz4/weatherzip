package service

import (
	"math"
	"runtime"

	"github.com/shirou/gopsutil/cpu"
)

type CPUService interface {
	GetCPUStats() (int, []float64, error)
}

type cpuService struct{}

func NewCPUService() CPUService {
	return &cpuService{}
}

func (s *cpuService) roundToOneDecimal(value float64) float64 {
	return math.Round(value*10) / 10
}

func (s *cpuService) GetCPUStats() (int, []float64, error) {
	cores := runtime.NumCPU()
	percentUsed, err := cpu.Percent(0, true)
	if err != nil {
		return 0, nil, err
	}

	for i, val := range percentUsed {
		percentUsed[i] = s.roundToOneDecimal(val)
	}

	return cores, percentUsed, nil
}
