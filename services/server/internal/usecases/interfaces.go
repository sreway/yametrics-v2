package usecases

import (
	"context"

	"github.com/sreway/yametrics-v2/pkg/metric"
)

type (
	Metric interface {
		Add(ctx context.Context, m *metric.Metric) error
		BatchAdd(ctx context.Context, m []*metric.Metric) error
		Get(ctx context.Context, id string, t metric.Type) (*metric.Metric, error)
		GetMany(ctx context.Context) ([]metric.Metric, error)
		StorageCheck(ctx context.Context) error
	}
)
