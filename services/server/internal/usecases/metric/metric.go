package metric

import (
	"context"
	"sort"

	domain "github.com/sreway/yametrics-v2/services/server/internal/domain/metric"
	"github.com/sreway/yametrics-v2/services/server/internal/usecases/adapters/storage"

	"github.com/sreway/yametrics-v2/pkg/metric"
)

type UseCase struct {
	secretKey string
	storage   storage.Storage
}

func (uc *UseCase) Add(ctx context.Context, m *metric.Metric) error {
	if uc.secretKey != "" && m.Hash != m.CalcHash(uc.secretKey) {
		return domain.NewMetricErr(m.ID, domain.ErrInvalidMetricHash)
	}
	return uc.storage.Add(ctx, m)
}

func (uc *UseCase) BatchAdd(ctx context.Context, m []*metric.Metric) error {
	if uc.secretKey != "" {
		for _, item := range m {
			if item.Hash != item.CalcHash(uc.secretKey) {
				return domain.NewMetricErr(item.ID, domain.ErrInvalidMetricHash)
			}
		}
	}

	return uc.storage.BatchAdd(ctx, m)
}

func (uc *UseCase) Get(ctx context.Context, id string, t metric.Type) (*metric.Metric, error) {
	m, err := uc.storage.Get(ctx, id, t)
	if err != nil {
		return nil, err
	}

	if uc.secretKey != "" {
		m.Hash = m.CalcHash(uc.secretKey)
	}

	return m, nil
}

func (uc *UseCase) GetMany(ctx context.Context) ([]metric.Metric, error) {
	metrics, err := uc.storage.GetMany(ctx)
	if err != nil {
		return nil, err
	}

	sort.Slice(metrics, func(i, j int) bool {
		if metrics[i].ID != metrics[j].ID {
			return metrics[i].ID < metrics[j].ID
		}
		return metrics[i].ID < metrics[j].ID
	})

	if uc.secretKey != "" {
		for _, item := range metrics {
			item.Hash = item.CalcHash(uc.secretKey)
		}
	}

	return metrics, nil
}

func (uc *UseCase) StorageCheck(ctx context.Context) error {
	return uc.storage.StorageCheck(ctx)
}

func New(s storage.Storage, secretKey string) *UseCase {
	uc := &UseCase{
		secretKey: secretKey,
		storage:   s,
	}
	return uc
}
