package metric

import (
	"encoding/json"

	"github.com/sreway/yametrics-v2/pkg/metric"
)

type Metrics map[string]*metric.Metric

func (m Metrics) MarshalJSON() ([]byte, error) {
	type aliasType = struct {
		Counter map[string]*metric.Metric `json:"counter"`
		Gauge   map[string]*metric.Metric `json:"gauge"`
	}

	aliasValue := aliasType{
		Counter: map[string]*metric.Metric{},
		Gauge:   map[string]*metric.Metric{},
	}

	for _, item := range m {
		if item.MType == metric.CounterType {
			aliasValue.Counter[item.ID] = item
		} else {
			aliasValue.Gauge[item.ID] = item
		}
	}

	return json.Marshal(aliasValue)
}

func (m Metrics) UnmarshalJSON(data []byte) error {
	type aliasType = struct {
		Counter map[string]*metric.Metric `json:"counter"`
		Gauge   map[string]*metric.Metric `json:"gauge"`
	}

	aliasValue := aliasType{
		Counter: map[string]*metric.Metric{},
		Gauge:   map[string]*metric.Metric{},
	}

	if err := json.Unmarshal(data, &aliasValue); err != nil {
		return err
	}

	for _, counter := range aliasValue.Counter {
		m[counter.ID] = counter
	}

	for _, gauge := range aliasValue.Gauge {
		m[gauge.ID] = gauge
	}

	return nil
}
