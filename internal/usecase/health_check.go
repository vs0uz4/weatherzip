package usecase

import (
	"time"
	"weatherzip/internal/infra/web/health"
	"weatherzip/internal/service"
)

type HealthCheckUseCase interface {
	GetHealth() (health.HealthStats, error)
}

type healthCheckUseCase struct {
	cpuService    service.CPUService
	memoryService service.MemoryService
	uptimeService service.UptimeService
}

func NewHealthCheckUseCase(cpu service.CPUService, memory service.MemoryService, uptime service.UptimeService) HealthCheckUseCase {
	return &healthCheckUseCase{cpu, memory, uptime}
}

func (h *healthCheckUseCase) GetHealth() (health.HealthStats, error) {
	start := time.Now()

	healthStats := health.HealthStats{
		Status:   "pass",
		Duration: "",
	}

	cores, percentUsed, err := h.cpuService.GetCPUStats()
	if err != nil {
		healthStats.Status = "fail"
	} else {
		healthStats.CPU = health.CPUStats{
			Cores:       cores,
			PercentUsed: percentUsed,
		}
	}

	total, used, free, available, percentUsedMem, err := h.memoryService.GetMemoryStats()
	if err != nil {
		healthStats.Status = "fail"
	} else {
		healthStats.Memory = health.MemoryStats{
			Total:       total,
			Used:        used,
			Free:        free,
			Available:   available,
			PercentUsed: percentUsedMem,
		}
	}

	if healthStats.Status == "pass" {
		healthStats.Message = "Alive and kicking!"
	} else {
		healthStats.Message = "Still alive, but not kicking!"
	}

	healthStats.Uptime = h.uptimeService.GetUptime()
	healthStats.Duration = time.Since(start).String()
	healthStats.Time = time.Now().Format(time.RFC3339)

	return healthStats, nil
}
