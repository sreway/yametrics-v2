package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/sreway/yametrics-v2/services/server/config"
)

func (d *Delivery) useMiddleware(cfg *config.HTTPConfig, r chi.Router) {
	r.Use(middleware.Compress(cfg.CompressLevel, cfg.CompressTypes...))
}
