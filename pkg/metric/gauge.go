package metric

import (
	"encoding/json"
	"fmt"
)

type Gauge struct {
	value float64
}

func (g Gauge) String() string {
	return fmt.Sprintf("%v", g.value)
}

func (g *Gauge) SetValue(v float64) {
	g.value = v
}

func (g Gauge) Value() float64 {
	return g.value
}

func (g Gauge) MarshalJSON() ([]byte, error) {
	return json.Marshal(g.value)
}

func (g *Gauge) Scan(value interface{}) error {
	v, ok := value.(float64)
	if !ok {
		return fmt.Errorf("Gauge_Scan: incorrect type %v", value)
	}
	g.value = v

	return nil
}

func NewGauge(v float64) *Gauge {
	return &Gauge{
		value: v,
	}
}
