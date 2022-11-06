package collector

import (
	"context"
	"sync"
	"time"

	repo "github.com/sreway/yametrics-v2/services/agent/internal/repository/storage/memory/metric"
	"github.com/sreway/yametrics-v2/services/agent/internal/usecases/adapters/storage"

	"github.com/sreway/yametrics-v2/pkg/metric"
	log "github.com/sreway/yametrics-v2/pkg/tools/logger"
)

type UseCase struct {
	secretKey string
	repo      storage.MemoryStorage
}

func (uc *UseCase) Collect(ctx context.Context, wg *sync.WaitGroup, polling time.Duration) {
	go uc.CollectRuntime(ctx, wg, polling)
	go uc.CollectUtil(ctx, wg)
}

func (uc *UseCase) Expose() []metric.Metric {
	metrics := uc.repo.Get()

	if uc.secretKey != "" {
		for idx, item := range metrics {
			sign := item.CalcHash(uc.secretKey)
			metrics[idx].Hash = sign
		}
	}

	return metrics
}

func (uc *UseCase) ResetCounter() {
	uc.repo.UpdateCounter("PollCount", 0)
}

func (uc *UseCase) CollectRuntime(ctx context.Context, wg *sync.WaitGroup, i time.Duration) {
	go func() {
		defer wg.Done()

		tick := time.NewTicker(i)
		defer tick.Stop()

		for {
			select {
			case <-tick.C:
				uc.SetRuntimeStats()
				log.Info("Agent: success collect runtime metrics")
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (uc *UseCase) CollectUtil(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	data := make(chan float64)
	stop := make(chan struct{})

	go func() {
		defer wg.Done()
		for {
			select {
			case <-stop:
				return
			default:
			}

			select {
			case <-ctx.Done():
				return
			case <-stop:
				return
			case data <- GetCPUPercent(10*time.Second, false):
			}
		}
	}()

	for {
		select {
		case cu := <-data:
			uc.SetMemmoryStats()
			uc.SetCPUStats(cu)
			log.Info("Agent: success collect utilization metrics")
		case <-ctx.Done():
			close(stop)
			return
		}
	}
}

func New(secretKey string) *UseCase {
	return &UseCase{
		secretKey: secretKey,
		repo:      repo.New(),
	}
}
