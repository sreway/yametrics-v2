package http

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	log "github.com/sreway/yametrics-v2/pkg/tools/logger"

	"github.com/sreway/yametrics-v2/services/server/config"
	"github.com/sreway/yametrics-v2/services/server/internal/usecases"

	"github.com/sreway/yametrics-v2/pkg/httpserver"
)

type Delivery struct {
	metrics usecases.Metric
	router  *chi.Mux
}

func New(uc usecases.Metric, cfg *config.HTTPConfig) *Delivery {
	d := &Delivery{
		metrics: uc,
	}
	d.router = d.initRouter(cfg)
	return d
}

func (d *Delivery) Run(ctx context.Context, cfg *config.HTTPConfig) error {
	httpServer := httpserver.New(d.router, httpserver.Addr(cfg.Address))
	serverCtx, stopServer := context.WithCancel(context.Background())

	go func() {
		<-ctx.Done()
		err := httpServer.Shutdown(serverCtx)
		if err != nil {
			log.Fatal(err.Error())
		}
		log.Info("Delivery: graceful shutdown http server")
		stopServer()
	}()

	err := httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	<-serverCtx.Done()
	return nil
}
