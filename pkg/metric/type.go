package metric

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
)

const (
	CounterType Type = "counter"
	GaugeType   Type = "gauge"
)

type (
	Type string
)

type (
	Metric struct {
		ID    string   `json:"id"`
		MType Type     `json:"type"`
		Delta *Counter `json:"delta,omitempty"`
		Value *Gauge   `json:"value,omitempty"`
		Hash  string   `json:"hash,omitempty"`
	}
)

func (m *Metric) CalcHash(key string) string {
	if m.MType == CounterType {
		return calcHash[int64](key, m.ID, m.MType, m.Delta.Value())
	}

	return calcHash[float64](key, m.ID, m.MType, m.Value.Value())
}

func New[T int64 | float64](id string, t Type, v T) *Metric {
	m := new(Metric)
	m.ID = id

	if t == CounterType {
		m.MType = CounterType
		cv, ok := any(v).(int64)
		if ok {
			m.Delta = NewCounter(cv)
		}
		return m
	}

	m.MType = GaugeType
	gv, ok := any(v).(float64)
	if ok {
		m.Value = NewGauge(gv)
	}

	return m
}

func calcHash[T int64 | float64](key, id string, t Type, v T) string {
	var msg string

	cv, ok := any(v).(int64)
	if ok {
		msg = fmt.Sprintf("%s:%v:%d", id, t, cv)
	}

	gv, ok := any(v).(float64)
	if ok {
		msg = fmt.Sprintf("%s:%v:%f", id, t, gv)
	}

	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(msg))
	hash := h.Sum(nil)
	return fmt.Sprintf("%x", hash)
}

func (m *Metric) UnmarshalJSON(data []byte) error {
	var err error

	type AliasType Metric
	aliasValue := &struct {
		*AliasType
		Delta int64   `json:"delta,omitempty"`
		Value float64 `json:"value,omitempty"`
	}{
		AliasType: (*AliasType)(m),
	}

	if err = json.Unmarshal(data, aliasValue); err != nil {
		return err
	}

	if m.MType == CounterType {
		m.Delta = &Counter{value: aliasValue.Delta}
	}

	if m.MType == GaugeType {
		m.Value = &Gauge{value: aliasValue.Value}
	}

	return nil
}

func (m *Metric) GetStrValue() string {
	if m.MType == CounterType {
		return m.Delta.String()
	}
	return m.Value.String()
}

func (t *Type) Valid() bool {
	allowed := []Type{CounterType, GaugeType}
	for _, v := range allowed {
		if v == *t {
			return true
		}
	}
	return false
}
