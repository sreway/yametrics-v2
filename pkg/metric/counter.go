package metric

import (
	"encoding/json"
	"fmt"
)

type Counter struct {
	value int64
}

func (c Counter) String() string {
	return fmt.Sprintf("%v", c.value)
}

func (c Counter) Value() int64 {
	return c.value
}

func (c *Counter) SetValue(v int64) {
	c.value = v
}

func (c *Counter) Inc(v int64) {
	c.value += v
}

func (c Counter) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.value)
}

func (c *Counter) Scan(value interface{}) error {
	v, ok := value.(int64)
	if !ok {
		return fmt.Errorf("Counter_Scan: incorrect type %v", value)
	}
	c.value = v

	return nil
}

func NewCounter(v int64) *Counter {
	return &Counter{
		value: v,
	}
}
