package server

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"

	log "github.com/sreway/yametrics-v2/pkg/tools/logger"

	// migrate tools
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func (s *Server) Migrate() error {
	if s.config.Postgres.MigrateURL != "" {
		m, err := migrate.New(s.config.Postgres.MigrateURL, s.config.Postgres.DSN)
		defer func() {
			_, _ = m.Close()
		}()

		if err != nil {
			return err
		}

		err = m.Up()
		if err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("Server_Postgres_Migrate: %w", err)
		}

		if errors.Is(err, migrate.ErrNoChange) {
			return errors.New("Server_Postgres_Migrate: no change")
		}
		log.Info("Server_Postgres_Migrate: success apply migrations")
		return nil
	}

	return errors.New("Server_Postgres_Migrate: missing migrate url")
}
