package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Postgres struct {
	Pool *pgxpool.Pool
}

func New(ctx context.Context, dsn string) (*Postgres, error) {
	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("Postgres_New: %w", err)
	}

	pool, err := pgxpool.ConnectConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("Postgres_New: %w", err)
	}

	return &Postgres{
		pool,
	}, nil
}

func (p *Postgres) Close() {
	p.Pool.Close()
}
