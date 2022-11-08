package http

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sreway/yametrics-v2/services/server/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sreway/yametrics-v2/pkg/metric"

	repoMemory "github.com/sreway/yametrics-v2/services/server/internal/repository/storage/memory/metric"
	metricService "github.com/sreway/yametrics-v2/services/server/internal/usecases/metric"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path, body string) *http.Response {
	reader := strings.NewReader(body)
	url := fmt.Sprintf("%s%s", ts.URL, path)
	req := httptest.NewRequest(method, url, reader)
	req.RequestURI = ""
	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	err = resp.Body.Close()
	require.NoError(t, err)
	return resp
}

func newTestStorage(t *testing.T, path string) *repoMemory.RepoMetric {
	s, err := repoMemory.New(path)
	require.NoError(t, err)
	return s
}

func TestDelivery_UpdateMetric(t *testing.T) {
	type want struct {
		statusCode int
	}

	type args struct {
		uri    string
		method string
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "send counter",
			args: args{
				uri:    "/update/counter/PollCount/100",
				method: http.MethodPost,
			},
			want: want{
				statusCode: 200,
			},
		},

		{
			name: "send gauge",
			args: args{
				uri:    "/update/gauge/RandomValue/10.8",
				method: http.MethodPost,
			},
			want: want{
				statusCode: 200,
			},
		},

		{
			name: "invalid value",
			args: args{
				uri:    "/update/counter/PollCount/none",
				method: http.MethodPost,
			},
			want: want{
				statusCode: 400,
			},
		},

		{
			name: "invalid type",
			args: args{
				uri:    "/update/unknown/PollCount/100",
				method: http.MethodPost,
			},
			want: want{
				statusCode: 501,
			},
		},

		{
			name: "invalid uri",
			args: args{
				uri:    "/update/unknown",
				method: http.MethodPost,
			},
			want: want{
				statusCode: 404,
			},
		},
	}

	cfg, err := config.New()
	assert.NoError(t, err)

	store := newTestStorage(t, "")

	for _, tt := range tests {
		ms := metricService.New(store, cfg.SecretKey)
		d := New(ms, &cfg.HTTP)
		ts := httptest.NewServer(d.router)

		t.Run(tt.name, func(t *testing.T) {
			resp := testRequest(t, ts, tt.args.method, tt.args.uri, ``)
			defer resp.Body.Close()
			assert.Equal(t, tt.want.statusCode, resp.StatusCode)
			err = resp.Body.Close()
			require.NoError(t, err)
		})
	}
}

func TestDelivery_GetMetric(t *testing.T) {
	type want struct {
		statusCode int
	}

	type args struct {
		uri    string
		method string
	}

	type fields struct {
		storageData []*metric.Metric
	}

	tests := []struct {
		args   args
		name   string
		fields fields
		want   want
	}{
		{
			name: "get counter",
			args: args{
				uri:    "/value/counter/PollCount",
				method: http.MethodGet,
			},
			fields: fields{
				storageData: []*metric.Metric{
					metric.New("PollCount", metric.CounterType, int64(21)),
				},
			},

			want: want{
				statusCode: 200,
			},
		},

		{
			name: "get gauge",
			args: args{
				uri:    "/value/gauge/testGauge",
				method: http.MethodGet,
			},

			fields: fields{
				storageData: []*metric.Metric{
					metric.New("testGauge", metric.GaugeType, float64(21)),
				},
			},

			want: want{
				statusCode: 200,
			},
		},

		{
			name: "non existent counter",
			args: args{
				uri:    "/value/counter/testCounter",
				method: http.MethodGet,
			},

			want: want{
				statusCode: 404,
			},
		},

		{
			name: "non existent gauge",
			args: args{
				uri:    "/value/gauge/testCounter",
				method: http.MethodGet,
			},

			want: want{
				statusCode: 404,
			},
		},

		{
			name: "invalid type",
			args: args{
				uri:    "/value/unknown/testCounter",
				method: http.MethodGet,
			},

			want: want{
				statusCode: 501,
			},
		},

		{
			name: "invalid uri",
			args: args{
				uri:    "/value/unknown",
				method: http.MethodGet,
			},

			want: want{
				statusCode: 404,
			},
		},
	}

	ctx := context.Background()
	cfg, err := config.New()
	assert.NoError(t, err)
	store := newTestStorage(t, "")

	for _, tt := range tests {
		ms := metricService.New(store, cfg.SecretKey)
		d := New(ms, &cfg.HTTP)
		ts := httptest.NewServer(d.router)
		err = store.BatchAdd(ctx, tt.fields.storageData)
		assert.NoError(t, err)

		t.Run(tt.name, func(t *testing.T) {
			resp := testRequest(t, ts, tt.args.method, tt.args.uri, ``)
			defer resp.Body.Close()
			assert.Equal(t, tt.want.statusCode, resp.StatusCode)
		})
	}
}

func TestDelivery_UpdateMetricJSON(t *testing.T) {
	type want struct {
		statusCode int
	}

	type args struct {
		uri    string
		method string
		body   string
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "update counter",
			args: args{
				uri:    "/update/",
				method: http.MethodPost,
				body:   `{"id":"testCounter","type":"counter","delta":1}`,
			},

			want: want{
				statusCode: 200,
			},
		},
		{
			name: "update gauge",
			args: args{
				uri:    "/update/",
				method: http.MethodPost,
				body:   `{"id":"testGauge","type":"gauge","value":20}`,
			},

			want: want{
				statusCode: 200,
			},
		},

		{
			name: "incorrect metric data",
			args: args{
				uri:    "/update/",
				method: http.MethodPost,
				body:   ``,
			},

			want: want{
				statusCode: 400,
			},
		},

		{
			name: "incorrect method",
			args: args{
				uri:    "/update/",
				method: http.MethodGet,
				body:   `{"id":"testGauge","type":"gauge","value":20}`,
			},

			want: want{
				statusCode: 405,
			},
		},
	}

	cfg, err := config.New()
	assert.NoError(t, err)

	store := newTestStorage(t, "")
	assert.NoError(t, err)

	for _, tt := range tests {
		ms := metricService.New(store, cfg.SecretKey)
		d := New(ms, &cfg.HTTP)
		ts := httptest.NewServer(d.router)

		t.Run(tt.name, func(t *testing.T) {
			resp := testRequest(t, ts, tt.args.method, tt.args.uri, tt.args.body)
			defer resp.Body.Close()
			assert.Equal(t, tt.want.statusCode, resp.StatusCode)
			err = resp.Body.Close()
			require.NoError(t, err)
		})
	}
}

func TestDelivery_GetMetricJSON(t *testing.T) {
	type want struct {
		statusCode int
	}

	type args struct {
		uri    string
		method string
		body   string
	}

	type fields struct {
		storageData []*metric.Metric
	}

	tests := []struct {
		name   string
		args   args
		fields fields
		want   want
	}{
		{
			name: "counter value",
			args: args{
				uri:    "/value/",
				method: http.MethodPost,
				body:   `{"id":"PollCounter","type":"counter"}`,
			},

			fields: fields{
				storageData: []*metric.Metric{
					metric.New("PollCounter", metric.CounterType, int64(1)),
				},
			},

			want: want{
				statusCode: 200,
			},
		},

		{
			name: "gauge value",
			args: args{
				uri:    "/value/",
				method: http.MethodPost,
				body:   `{"id":"testGauge","type":"gauge"}`,
			},

			fields: fields{
				storageData: []*metric.Metric{
					metric.New("testGauge", metric.GaugeType, float64(1)),
				},
			},

			want: want{
				statusCode: 200,
			},
		},

		{
			name: "incorrect type",
			args: args{
				uri:    "/value/",
				method: http.MethodPost,
				body:   `{"id":"testGauge","type":"incorrect"}`,
			},

			want: want{
				statusCode: 501,
			},
		},

		{
			name: "not exist value",
			args: args{
				uri:    "/value/",
				method: http.MethodPost,
				body:   `{"id":"not-exist-id","type":"gauge"}`,
			},
			want: want{
				statusCode: 404,
			},
		},
	}

	ctx := context.Background()
	cfg, err := config.New()
	assert.NoError(t, err)

	store := newTestStorage(t, "")
	assert.NoError(t, err)

	for _, tt := range tests {
		ms := metricService.New(store, cfg.SecretKey)
		d := New(ms, &cfg.HTTP)
		ts := httptest.NewServer(d.router)
		err = store.BatchAdd(ctx, tt.fields.storageData)
		assert.NoError(t, err)

		t.Run(tt.name, func(t *testing.T) {
			resp := testRequest(t, ts, tt.args.method, tt.args.uri, tt.args.body)
			defer resp.Body.Close()
			assert.Equal(t, tt.want.statusCode, resp.StatusCode)
		})
	}
}

func TestDelivery_BatchMetrics(t *testing.T) {
	type want struct {
		statusCode int
	}

	type args struct {
		uri    string
		method string
		body   string
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "send batch metrics",
			args: args{
				uri:    "/updates/",
				method: http.MethodPost,
				body:   `[{"id":"tGauge","type":"gauge","value":20},{"id":"tCounter","type":"counter","delta":1}]`,
			},

			want: want{
				statusCode: 200,
			},
		},

		{
			name: "send incorrect batch metrics",
			args: args{
				uri:    "/updates/",
				method: http.MethodPost,
				body:   ``,
			},

			want: want{
				statusCode: 400,
			},
		},
	}

	cfg, err := config.New()
	assert.NoError(t, err)

	store := newTestStorage(t, "")
	assert.NoError(t, err)

	for _, tt := range tests {
		ms := metricService.New(store, cfg.SecretKey)
		d := New(ms, &cfg.HTTP)
		ts := httptest.NewServer(d.router)

		t.Run(tt.name, func(t *testing.T) {
			resp := testRequest(t, ts, tt.args.method, tt.args.uri, tt.args.body)
			defer resp.Body.Close()
			assert.Equal(t, tt.want.statusCode, resp.StatusCode)
			err = resp.Body.Close()
			require.NoError(t, err)
		})
	}
}

func TestDelivery_HealthCheckStorage(t *testing.T) {
	type want struct {
		statusCode int
	}

	type args struct {
		uri    string
		method string
	}

	type fields struct {
		file string
	}

	tests := []struct {
		name   string
		args   args
		fields fields
		want   want
	}{
		{
			name: "get ping for memory storage",
			args: args{
				uri:    "/ping",
				method: http.MethodGet,
			},

			fields: fields{
				file: "",
			},

			want: want{
				statusCode: 501,
			},
		},
	}

	cfg, err := config.New()
	assert.NoError(t, err)

	for _, tt := range tests {
		store := newTestStorage(t, tt.fields.file)
		assert.NoError(t, err)

		ms := metricService.New(store, cfg.SecretKey)
		d := New(ms, &cfg.HTTP)
		ts := httptest.NewServer(d.router)
		assert.NoError(t, err)

		t.Run(tt.name, func(t *testing.T) {
			resp := testRequest(t, ts, tt.args.method, tt.args.uri, ``)
			defer resp.Body.Close()
			assert.Equal(t, tt.want.statusCode, resp.StatusCode)
		})
	}
}
