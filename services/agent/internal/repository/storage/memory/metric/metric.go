package metric

import (
	"sync"

	"github.com/sreway/yametrics-v2/pkg/metric"
)

type RepoMetric struct {
	metrics map[string]*metric.Metric
	mu      sync.RWMutex
}

func (r *RepoMetric) AddCounterWithoutLock(id string, v int64) {
	if val, ok := r.metrics[id]; !ok {
		r.metrics[id] = metric.New(id, metric.CounterType, v)
	} else {
		val.Delta.Inc(v)
	}
}

func (r *RepoMetric) AddCounter(id string, v int64) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if val, ok := r.metrics[id]; !ok {
		r.metrics[id] = metric.New(id, metric.CounterType, v)
	} else {
		val.Delta.Inc(v)
	}
}

func (r *RepoMetric) UpdateCounter(id string, v int64) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.metrics[id]; ok {
		r.metrics[id].Delta.SetValue(v)
	}
}

func (r *RepoMetric) AddGaugeWithoutLock(id string, v float64) {
	if val, ok := r.metrics[id]; !ok {
		r.metrics[id] = metric.New(id, metric.GaugeType, v)
	} else {
		val.Value.SetValue(v)
	}
}

func (r *RepoMetric) AddGauge(id string, v float64) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if val, ok := r.metrics[id]; !ok {
		r.metrics[id] = metric.New(id, metric.GaugeType, v)
	} else {
		val.Value.SetValue(v)
	}
}

func (r *RepoMetric) UpdateGauge(id string, v float64) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.metrics[id]; ok {
		r.metrics[id].Value.SetValue(v)
	}
}

func (r *RepoMetric) Get() []metric.Metric {
	r.mu.RLock()
	defer r.mu.RUnlock()

	m := []metric.Metric{}

	for _, v := range r.metrics {
		var cm *metric.Metric
		if v.MType == metric.CounterType {
			cm = metric.New(v.ID, metric.CounterType, v.Delta.Value())
		} else {
			cm = metric.New(v.ID, metric.GaugeType, v.Value.Value())
		}

		m = append(m, *cm)
	}

	return m
}

func (r *RepoMetric) Lock() {
	r.mu.Lock()
}

func (r *RepoMetric) Unlock() {
	r.mu.Unlock()
}

func (r *RepoMetric) RLock() {
	r.mu.RLock()
}

func (r *RepoMetric) RUnlock() {
	r.mu.RUnlock()
}

func New() *RepoMetric {
	return &RepoMetric{
		metrics: map[string]*metric.Metric{},
		mu:      sync.RWMutex{},
	}
}
