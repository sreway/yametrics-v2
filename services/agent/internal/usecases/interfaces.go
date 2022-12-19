package usecases

import (
	"context"
	"sync"
	"time"

	"github.com/sreway/yametrics-v2/pkg/metric"
)

type Collector interface {
	Collect(ctx context.Context, wg *sync.WaitGroup, polling time.Duration)
	Expose() []metric.Metric
	ResetCounter()
}

type Sender interface {
	Send(ctx context.Context, m []metric.Metric) error
	Close() error
}
