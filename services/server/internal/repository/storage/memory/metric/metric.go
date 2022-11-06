package metric

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	domain "github.com/sreway/yametrics-v2/services/server/internal/domain/metric"

	"github.com/sreway/yametrics-v2/pkg/metric"
	log "github.com/sreway/yametrics-v2/pkg/tools/logger"
)

type RepoMetric struct {
	db      Metrics
	fileObj *os.File
	mu      sync.RWMutex
}

func (r *RepoMetric) Add(ctx context.Context, m *metric.Metric) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_ = ctx

	val, ok := r.db[m.ID]

	switch {
	case !m.MType.Valid():
		return domain.NewMetricErr(m.ID, domain.ErrInvalidMetricType)
	case ok && val.MType != m.MType:
		return domain.NewMetricErr(m.ID, domain.ErrInvalidMetricType)
	case ok && val.MType == metric.CounterType:
		val.Delta.Inc(m.Delta.Value())
		return nil
	default:
		r.db[m.ID] = m
	}

	return nil
}

func (r *RepoMetric) BatchAdd(ctx context.Context, m []*metric.Metric) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_ = ctx

	for _, item := range m {
		val, ok := r.db[item.ID]
		switch {
		case !item.MType.Valid(), ok && val.MType != item.MType:
			return domain.NewMetricErr(item.ID, domain.ErrInvalidMetricType)
		case ok && val.MType == metric.CounterType:
			val.Delta.Inc(item.Delta.Value())
		default:
			r.db[item.ID] = item
		}
	}

	return nil
}

func (r *RepoMetric) Get(ctx context.Context, id string, t metric.Type) (*metric.Metric, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_ = ctx

	val, ok := r.db[id]

	switch {
	case !t.Valid():
		return nil, domain.NewMetricErr(id, domain.ErrInvalidMetricType)
	case !ok, val.MType != t:
		return nil, domain.NewMetricErr(id, domain.ErrMetricNotFound)
	default:
		var cm *metric.Metric
		if val.MType == metric.CounterType {
			cm = metric.New(val.ID, metric.CounterType, val.Delta.Value())
		} else {
			cm = metric.New(val.ID, metric.GaugeType, val.Value.Value())
		}

		return cm, nil
	}
}

func (r *RepoMetric) GetMany(ctx context.Context) ([]metric.Metric, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_ = ctx

	m := []metric.Metric{}

	for _, v := range r.db {
		// copy metric with values
		var cm *metric.Metric
		if v.MType == metric.CounterType {
			cm = metric.New(v.ID, metric.CounterType, v.Delta.Value())
		} else {
			cm = metric.New(v.ID, metric.GaugeType, v.Value.Value())
		}

		m = append(m, *cm)
	}

	return m, nil
}

func (r *RepoMetric) Close() error {
	err := r.fileObj.Close()
	if err != nil {
		return fmt.Errorf("MemoryStorage_Close: %w", err)
	}

	return nil
}

func (r *RepoMetric) Load() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := json.NewDecoder(r.fileObj).Decode(&r.db); err != nil {
		return fmt.Errorf("MemoryStorage_Load: %w: cant't decode metrics", ErrLoadMetrics)
	}

	log.Info("MemoryStorage_Load: success load metrics")

	return nil
}

func (r *RepoMetric) Store() error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	err := r.fileObj.Truncate(0)
	if err != nil {
		return fmt.Errorf("MemoryStorage_Store: %w cat't truncate file", ErrStoreMetrics)
	}

	_, err = r.fileObj.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("MemoryStorage_Store: %w cat't seek file", ErrStoreMetrics)
	}

	if err = json.NewEncoder(r.fileObj).Encode(r.db); err != nil {
		return fmt.Errorf("MemoryStorage_Store: %w: cant't encode metrics", ErrStoreMetrics)
	}

	log.Info("MemoryStorage_Store: success save metrics to file")

	return nil
}

func (r *RepoMetric) StorageCheck(ctx context.Context) error {
	_ = ctx
	return fmt.Errorf("MemoryStorage_StorageCheck: %s", "PostgresStorage unavailable")
}

func OpenFile(path string) (*os.File, error) {
	flag := os.O_RDWR | os.O_CREATE
	fileObj, err := os.OpenFile(path, flag, 0o644)
	if err != nil {
		return nil, fmt.Errorf("MemoryStorage_File_New: %w %s", ErrOpenFile, path)
	}

	return fileObj, nil
}

func New(filePath string) (*RepoMetric, error) {
	s := RepoMetric{}
	s.db = map[string]*metric.Metric{}

	if filePath != "" {
		fileObj, err := OpenFile(filePath)
		if err != nil {
			return nil, err
		}

		s.fileObj = fileObj
	}

	return &s, nil
}
