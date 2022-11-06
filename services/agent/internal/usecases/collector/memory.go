package collector

import (
	"github.com/shirou/gopsutil/v3/mem"
)

func (uc *UseCase) SetMemmoryStats() {
	stats, _ := mem.VirtualMemory()
	uc.repo.RLock()
	defer uc.repo.RUnlock()

	uc.repo.AddGaugeWithoutLock("TotalMemory", float64(stats.Total))
	uc.repo.AddGaugeWithoutLock("FreeMemory", float64(stats.Free))
}
