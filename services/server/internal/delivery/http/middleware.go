package http

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	log "github.com/sreway/yametrics-v2/pkg/tools/logger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/sreway/yametrics-v2/services/server/config"
)

func (d *Delivery) useMiddleware(cfg *config.DeliveryConfig, r chi.Router) {
	r.Use(middleware.Compress(cfg.CompressLevel, cfg.CompressTypes...))
	if cfg.TrustedSubnet != nil {
		r.Use(TrustedSubnet(cfg.TrustedSubnet))
	}
}

func TrustedSubnet(subnet *net.IPNet) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var errMsg string
			if rip := r.Header.Get("X-Real-IP"); rip != "" {
				ip := net.ParseIP(rip)
				if ip == nil {
					errMsg = fmt.Sprintf("can't parse ip from X-Real-IP header: %s", ip)
					trustedSubnetFailed(w, errMsg)
					return
				}
				if !subnet.Contains(ip) {
					errMsg = fmt.Sprintf("ip address %s not allowed", ip.String())
					trustedSubnetFailed(w, errMsg)
					return
				}
				next.ServeHTTP(w, r)
				return
			}
			errMsg = "missing X-Real-IP header"
			trustedSubnetFailed(w, errMsg)
		})
	}
}

func trustedSubnetFailed(w http.ResponseWriter, msg string) {
	var resp struct {
		Error string `json:"error"`
	}
	resp.Error = msg
	w.WriteHeader(http.StatusForbidden)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Error(err.Error())
	}
}
