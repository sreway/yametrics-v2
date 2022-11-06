package server

import (
	"context"
	"errors"
	"time"

	"github.com/sreway/yametrics-v2/pkg/postgres"
	repoMemory "github.com/sreway/yametrics-v2/services/server/internal/repository/storage/memory/metric"
	repoPostgres "github.com/sreway/yametrics-v2/services/server/internal/repository/storage/postgres/metric"
	"github.com/sreway/yametrics-v2/services/server/internal/usecases/adapters/storage"

	log "github.com/sreway/yametrics-v2/pkg/tools/logger"
)

func (s *Server) InitStorage(ctx context.Context) error {
	if s.config.Postgres.DSN != "" {
		pg, err := postgres.New(ctx, s.config.Postgres.DSN)
		if err == nil {
			s.storage = repoPostgres.New(pg)
			err = s.Migrate()
			if err != nil {
				log.Warn(err.Error())
			}
			log.Info("Server: use postgres storage")
			return nil
		}

		log.Warn(err.Error())
	}

	st, err := repoMemory.New(s.config.MemoryStorage.StoreFile)
	if err != nil {
		return err
	}

	s.storage = st

	switch t := s.storage.(type) {
	case storage.MemoryStorage:
		if s.config.MemoryStorage.Restore {
			err = t.Load()
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
						err = t.Store()
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
	default:
		return errors.New("invalid storage")
	}
	log.Info("Server: use memory storage")
	return nil
}
