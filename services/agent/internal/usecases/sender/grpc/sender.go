package grpc

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"

	"google.golang.org/grpc/credentials"

	"github.com/sreway/yametrics-v2/pkg/metric"
	log "github.com/sreway/yametrics-v2/pkg/tools/logger"
	"github.com/sreway/yametrics-v2/pkg/tools/pem"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/sreway/yametrics-v2/proto/metric/v1"
	"github.com/sreway/yametrics-v2/services/agent/config"
)

var ErrEmptyMetrics = errors.New("empty metrics data")

type UseCase struct {
	grpc pb.MetricServiceClient
	conn *grpc.ClientConn
}

func New(cfg *config.Config) (*UseCase, error) {
	var opts []grpc.DialOption
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

		tlsConfig := tls.Config{}
		tlsConfig.RootCAs = x509.NewCertPool()

		for _, cert := range certs {
			tlsConfig.RootCAs.AddCert(cert)
		}

		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(&tlsConfig)))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	conn, err := grpc.Dial(cfg.ServerAddress, opts...)
	if err != nil {
		return nil, err
	}
	client := pb.NewMetricServiceClient(conn)

	return &UseCase{
		grpc: client,
		conn: conn,
	}, nil
}

func (uc *UseCase) Send(ctx context.Context, m []metric.Metric) error {
	if len(m) == 0 {
		log.Warn(ErrEmptyMetrics.Error())
		return ErrEmptyMetrics
	}
	metrics := make([]*pb.Metric, 0, len(m))

	for _, i := range m {
		rm := new(pb.Metric)
		rm.Id = i.ID
		switch i.MType {
		case metric.CounterType:
			rm.Type = pb.Type_COUNTER
			rm.Delta = i.Delta.Value()
		case metric.GaugeType:
			rm.Type = pb.Type_GAUGE
			rm.Value = i.Value.Value()
		}
		rm.Hash = i.Hash
		metrics = append(metrics, rm)
	}

	resp, err := uc.grpc.BatchAdd(ctx, &pb.BatchAddMetricRequest{
		Metrics: metrics,
	})
	if err != nil {
		return err
	}

	if resp.Error != "" {
		return fmt.Errorf("sender.Send:%s", resp.Error)
	}

	return nil
}

func (uc *UseCase) Close() error {
	return uc.conn.Close()
}
