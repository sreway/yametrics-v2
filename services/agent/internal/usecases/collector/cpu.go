package collector

import (
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
)

func (uc *UseCase) SetCPUStats(cu float64) {
	uc.repo.AddGauge("CPUutilization1", cu)
}

func GetCPUPercent(i time.Duration, percpu bool) float64 {
	percent, _ := cpu.Percent(i, percpu)
	return percent[0]
}
