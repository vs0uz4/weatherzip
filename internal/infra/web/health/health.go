package health

type CPUStats struct {
	Cores       int       `json:"cores"`
	PercentUsed []float64 `json:"percent_used"`
}

type MemoryStats struct {
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	Free        uint64  `json:"free"`
	Available   uint64  `json:"available"`
	PercentUsed float64 `json:"percent_used"`
}

type HealthStats struct {
	CPU      CPUStats    `json:"cpu"`
	Memory   MemoryStats `json:"memory"`
	Uptime   string      `json:"uptime"`
	Duration string      `json:"duration"`
	Status   string      `json:"status"`
	Message  string      `json:"message"`
	Time     string      `json:"time"`
}
