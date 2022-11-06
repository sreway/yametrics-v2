package storage

import "github.com/sreway/yametrics-v2/pkg/metric"

type Storage interface {
	AddCounter(id string, v int64)
	UpdateCounter(id string, v int64)
	AddGauge(id string, v float64)
	UpdateGauge(id string, v float64)
	Get() []metric.Metric
}

type MemoryStorage interface {
	Storage
	AddCounterWithoutLock(id string, v int64)
	AddGaugeWithoutLock(id string, v float64)
	RLock()
	RUnlock()
	Lock()
	Unlock()
}
