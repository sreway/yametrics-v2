package http

import (
	"github.com/go-chi/chi/v5"

	"github.com/sreway/yametrics-v2/services/server/config"
)

func (d *Delivery) initRouter(cfg *config.HTTPConfig) *chi.Mux {
	router := chi.NewRouter()
	d.useMiddleware(cfg, router)

	router.Route("/", func(r chi.Router) {
		r.Get("/", d.Index)
	})

	router.Route("/update", func(r chi.Router) {
		r.Post("/", d.UpdateMetricJSON)
		r.Post("/{type}/{id}/{value}", d.UpdateMetric)
	})

	router.Route("/updates", func(r chi.Router) {
		r.Post("/", d.BatchMetrics)
	})

	router.Route("/value", func(r chi.Router) {
		r.Post("/", d.GetMetricJSON)
		r.Get("/{type}/{id}", d.GetMetric)
	})

	router.Route("/ping", func(r chi.Router) {
		r.Get("/", d.HealthCheckStorage)
	})

	return router
}
