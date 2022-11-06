package metric

import (
	"context"
	"errors"
	"fmt"

	domain "github.com/sreway/yametrics-v2/services/server/internal/domain/metric"

	"github.com/jackc/pgx/v4"

	"github.com/sreway/yametrics-v2/pkg/postgres"

	"github.com/sreway/yametrics-v2/pkg/metric"
)

type RepoMetric struct {
	db *postgres.Postgres
}

func (r *RepoMetric) Add(ctx context.Context, m *metric.Metric) error {
	var query string

	tx, err := r.db.Pool.BeginTx(ctx, pgx.TxOptions{})
	defer func() {
		_ = tx.Rollback(ctx)
	}()
	if err != nil {
		return err
	}

	switch m.MType {
	case metric.CounterType:
		query = "INSERT INTO metrics (name, type, delta) VALUES ($1, $2, $3)" +
			" ON CONFLICT ON CONSTRAINT uniq_name_type DO UPDATE set delta = $3 + metrics.delta"
		_, err = tx.Exec(ctx, query, m.ID, m.MType, m.Delta.Value())
		if err != nil {
			return fmt.Errorf("PostgresStorage_Add: %w", err)
		}
	case metric.GaugeType:
		query = "INSERT INTO metrics (name, type, value) VALUES ($1, $2, $3) " +
			"ON CONFLICT ON CONSTRAINT uniq_name_type DO UPDATE set value=$3"
		_, err = tx.Exec(ctx, query, m.ID, m.MType, m.Value.Value())
		if err != nil {
			return fmt.Errorf("PostgresStorage_Add: %w", err)
		}
	default:
		return fmt.Errorf("PostgresStorage_Add: %w", domain.ErrInvalidMetricType)
	}

	return tx.Commit(ctx)
}

func (r *RepoMetric) BatchAdd(ctx context.Context, m []*metric.Metric) error {
	var query string

	tx, err := r.db.Pool.BeginTx(ctx, pgx.TxOptions{})
	defer func() {
		_ = tx.Rollback(ctx)
	}()
	if err != nil {
		return err
	}

	for _, item := range m {
		switch item.MType {
		case metric.CounterType:
			query = "INSERT INTO metrics (name, type, delta) VALUES ($1, $2, $3)" +
				" ON CONFLICT ON CONSTRAINT uniq_name_type DO UPDATE set delta = $3 + metrics.delta"
			_, err = tx.Exec(ctx, query, item.ID, item.MType, item.Delta.Value())
			if err != nil {
				return fmt.Errorf("PostgresStorage_BatchAdd: %w", err)
			}
		case metric.GaugeType:
			query = "INSERT INTO metrics (name, type, value) VALUES ($1, $2, $3) " +
				"ON CONFLICT ON CONSTRAINT uniq_name_type DO UPDATE set value=$3"
			_, err = tx.Exec(ctx, query, item.ID, item.MType, item.Value.Value())
			if err != nil {
				return fmt.Errorf("PostgresStorage_BatchAdd: %w", err)
			}
		default:
			return fmt.Errorf("PostgresStorage_BatchAdd: %w", domain.ErrInvalidMetricType)
		}
	}

	return tx.Commit(ctx)
}

func (r *RepoMetric) Get(ctx context.Context, id string, t metric.Type) (*metric.Metric, error) {
	var query string

	m := metric.Metric{}

	query = "SELECT delta, value FROM metrics WHERE name = $1 and type = $2"
	err := r.db.Pool.QueryRow(ctx, query, id, t).Scan(&m.Delta, &m.Value)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("PostgresStorage_Get: %w", domain.ErrMetricNotFound)
		}
		return nil, fmt.Errorf("PostgresStorage_Get: %w", err)
	}

	m.ID = id
	m.MType = t

	return &m, nil
}

func (r *RepoMetric) GetMany(ctx context.Context) ([]metric.Metric, error) {
	metrics := []metric.Metric{}
	query := "SELECT name, type, delta, value FROM metrics"
	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("PostgresStorage_GetMany: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var m metric.Metric
		err = rows.Scan(&m.ID, &m.MType, &m.Delta, &m.Value)
		if err != nil {
			return nil, fmt.Errorf("PostgresStorage_GetMany: %w", err)
		}
		metrics = append(metrics, m)
	}

	return metrics, nil
}

func (r *RepoMetric) Close() error {
	r.db.Close()
	return nil
}

func (r *RepoMetric) StorageCheck(ctx context.Context) error {
	if err := r.db.Pool.Ping(ctx); err != nil {
		return fmt.Errorf("RepoMetric_: %w", err)
	}
	return nil
}

func New(db *postgres.Postgres) *RepoMetric {
	return &RepoMetric{
		db,
	}
}
