package server

import (
	"context"
	"time"

	"github.com/sreway/yametrics-v2/services/server/config"

	"github.com/sreway/yametrics-v2/pkg/postgres"
	repoMemory "github.com/sreway/yametrics-v2/services/server/internal/repository/storage/memory/metric"
	repoPostgres "github.com/sreway/yametrics-v2/services/server/internal/repository/storage/postgres/metric"
	"github.com/sreway/yametrics-v2/services/server/internal/usecases/adapters/storage"

	log "github.com/sreway/yametrics-v2/pkg/tools/logger"
)

func InitPostgres(ctx context.Context, cfg config.PostgresConfig) (storage.Storage, error) {
	if cfg.DSN == "" {
		return nil, ErrEmptyDSN
	}
	pg, err := postgres.New(ctx, cfg.DSN)
	if err != nil {
		return nil, err
	}
	s := repoPostgres.New(pg)
	return s, nil
}

func InitMemory(cfg config.MemoryStorageConfig) (storage.Storage, error) {
	s, err := repoMemory.New(cfg.StoreFile)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Server) InitStorage(ctx context.Context) error {
	if s.config.Postgres.DSN != "" {
		pg, err := InitPostgres(ctx, s.config.Postgres)
		if err == nil {
			log.Info("Server: use postgres storage")
			s.storage = pg
			err = s.Migrate()
			log.Warn(err.Error())
			return nil
		}

		log.Warn(err.Error())
	}

	st, err := InitMemory(s.config.MemoryStorage)
	if err != nil {
		return err
	}

	s.storage = st
	mst := s.storage.(storage.MemoryStorage)
	if s.config.MemoryStorage.Restore {
		err = mst.Load()
		if err != nil {
			log.Info(err.Error())
		}
	}

	if s.config.MemoryStorage.StoreInterval != 0 {
		go func() {
			tick := time.NewTicker(s.config.MemoryStorage.StoreInterval)
			defer tick.Stop()
			for {
				select {
				case <-tick.C:
					err = mst.Store()
					if err != nil {
						log.Info(err.Error())
					}
					log.Info("Server: success store metrics")
				case <-ctx.Done():
					return
				}
			}
		}()
	}

	log.Info("Server: use memory storage")
	return nil
}
