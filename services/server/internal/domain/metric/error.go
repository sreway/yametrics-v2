package metric

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidMetricValue = errors.New("invalid metric value")
	ErrInvalidMetricType  = errors.New("invalid metric type")
	ErrInvalidMetricHash  = errors.New("invalid metric hash")
	ErrMetricNotFound     = errors.New("metric not found")
)

type (
	ErrMetric struct {
		ID    string
		error error
	}
)

func NewMetricErr(id string, err error) error {
	return &ErrMetric{
		ID:    id,
		error: err,
	}
}

func (e *ErrMetric) Error() string {
	return fmt.Sprintf("Metric_Error[%s]: %s", e.ID, e.error)
}

func (e ErrMetric) Is(err error) bool {
	return errors.Is(e.error, err)
}
