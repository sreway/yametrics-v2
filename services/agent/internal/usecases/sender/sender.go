package sender

import (
	"context"
	"crypto/x509"
	"errors"

	"github.com/sreway/yametrics-v2/pkg/tools/pem"
	"github.com/sreway/yametrics-v2/services/agent/config"

	"github.com/sreway/yametrics-v2/pkg/httpclient"
	"github.com/sreway/yametrics-v2/pkg/metric"
	log "github.com/sreway/yametrics-v2/pkg/tools/logger"
)

var ErrEmptyMetrics = errors.New("empty metrics data")

type UseCase struct {
	http *httpclient.Client
}

func (uc *UseCase) Send(ctx context.Context, endpoint string, m []metric.Metric) error {
	if len(m) == 0 {
		log.Warn(ErrEmptyMetrics.Error())
		return ErrEmptyMetrics
	}

	r, err := uc.http.R().SetContext(ctx).SetBody(&m).Post(endpoint)
	if err != nil {
		return httpclient.NewErrHTTPClient(r.StatusCode(), err.Error())
	}

	if r.StatusCode() != 200 {
		log.Error("Sender_Send: status code is not 200")
		return httpclient.NewErrHTTPClient(r.StatusCode(), "status code is not 200")
	}

	return nil
}

func New(cfg *config.Config) (*UseCase, error) {
	url := cfg.ServerHTTPScheme + "://" + cfg.ServerAddress
	if cfg.ServerPublicKey != "" {
		var certs []*x509.Certificate
		pemData, err := pem.ParsePEM(cfg.ServerPublicKey)
		if err != nil {
			return nil, err
		}

		for _, cert := range pemData.Certificate {
			x509Cert, err := x509.ParseCertificate(cert)
			if err != nil {
				return nil, err
			}
			certs = append(certs, x509Cert)
		}

		return &UseCase{
			http: httpclient.New(httpclient.WithBaseURL(url), httpclient.WithCerts(certs...)),
		}, nil
	}

	return &UseCase{
		http: httpclient.New(httpclient.WithBaseURL(url)),
	}, nil
}
