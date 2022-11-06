package http

import (
	"github.com/go-chi/chi/v5"

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

func (d *Delivery) Run(cfg *config.HTTPConfig) error {
	httpServer := httpserver.New(d.router, httpserver.Addr(cfg.Address))
	err := httpServer.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}
