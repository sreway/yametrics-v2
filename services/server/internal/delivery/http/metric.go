package http

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"html/template"
	"net/http"
	"strconv"
	"time"

	domain "github.com/sreway/yametrics-v2/services/server/internal/domain/metric"

	"github.com/go-chi/chi/v5"

	"github.com/sreway/yametrics-v2/pkg/metric"
	log "github.com/sreway/yametrics-v2/pkg/tools/logger"
)

var (
	//go:embed templates/index.gohtml
	templatesFS   embed.FS
	templateFiles = map[string]string{
		"/": "templates/index.gohtml",
	}
)

func (d *Delivery) Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	p, ok := templateFiles[r.URL.Path]
	if !ok {
		log.Error(ErrTemplateNotFound.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	m, err := d.metrics.GetMany(r.Context())
	if err != nil {
		HandelErrMetric(w, err)
	}

	tmpl, err := template.ParseFS(templatesFS, p)
	if err != nil {
		log.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = tmpl.Execute(w, m)
	if err != nil {
		log.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (d *Delivery) GetMetric(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	t := chi.URLParam(r, "type")
	id := chi.URLParam(r, "id")

	m, err := d.metrics.Get(r.Context(), id, metric.Type(t))
	if err != nil {
		HandelErrMetric(w, err)
		return
	}

	_, err = w.Write([]byte(m.GetStrValue()))
	if err != nil {
		HandelErrMetric(w, err)
		return
	}
}

func (d *Delivery) GetMetricJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	reqm := metric.Metric{}

	if err := decoder.Decode(&reqm); err != nil {
		log.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
	}

	m, err := d.metrics.Get(r.Context(), reqm.ID, reqm.MType)
	if err != nil {
		HandelErrMetric(w, err)
		return
	}

	if err = json.NewEncoder(w).Encode(m); err != nil {
		log.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (d *Delivery) UpdateMetricJSON(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	reqm := metric.Metric{}

	if err := decoder.Decode(&reqm); err != nil {
		log.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
	}

	err := d.metrics.Add(r.Context(), &reqm)
	if err != nil {
		HandelErrMetric(w, err)
		return
	}
}

func (d *Delivery) UpdateMetric(w http.ResponseWriter, r *http.Request) {
	var m *metric.Metric

	t := chi.URLParam(r, "type")
	id := chi.URLParam(r, "id")
	v := chi.URLParam(r, "value")

	switch metric.Type(t) {
	case metric.CounterType:
		iv, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			HandelErrMetric(w, domain.ErrInvalidMetricValue)
			return
		}

		m = metric.New(id, metric.CounterType, iv)

	case metric.GaugeType:
		fv, err := strconv.ParseFloat(v, 64)
		if err != nil {
			HandelErrMetric(w, domain.ErrInvalidMetricValue)
			return
		}

		m = metric.New(id, metric.GaugeType, fv)

	default:
		HandelErrMetric(w, domain.ErrInvalidMetricType)
		return
	}

	err := d.metrics.Add(r.Context(), m)
	if err != nil {
		HandelErrMetric(w, err)
		return
	}
}

func (d *Delivery) BatchMetrics(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	reqm := []*metric.Metric{}

	if err := decoder.Decode(&reqm); err != nil {
		log.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
	}

	err := d.metrics.BatchAdd(r.Context(), reqm)
	if err != nil {
		HandelErrMetric(w, err)
		return
	}
}

func (d *Delivery) HealthCheckStorage(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()

	err := d.metrics.StorageCheck(ctx)
	if err != nil {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}
}

func HandelErrMetric(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrTemplateNotFound):
		w.WriteHeader(http.StatusNotFound)
	case errors.Is(err, domain.ErrInvalidMetricHash):
		w.WriteHeader(http.StatusBadRequest)
	case errors.Is(err, domain.ErrInvalidMetricType):
		w.WriteHeader(http.StatusNotImplemented)
	case errors.Is(err, domain.ErrInvalidMetricValue):
		w.WriteHeader(http.StatusBadRequest)
	case errors.Is(err, domain.ErrMetricNotFound):
		w.WriteHeader(http.StatusNotFound)
	default:
		w.WriteHeader(http.StatusNotImplemented)
	}
	log.Error(err.Error())
}
