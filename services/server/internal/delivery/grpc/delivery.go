package grpc

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	log "github.com/sreway/yametrics-v2/pkg/tools/logger"
	pb "github.com/sreway/yametrics-v2/proto/metric/v1"
	"github.com/sreway/yametrics-v2/services/server/config"
	"github.com/sreway/yametrics-v2/services/server/internal/usecases"
)

type Delivery struct {
	s *grpc.Server
}

func New(uc usecases.Metric, cfg *config.DeliveryConfig) (*Delivery, error) {
	var (
		s   *grpc.Server
		opt []grpc.ServerOption
	)
	d := new(Delivery)

	if cfg.CryptoCrt != "" && cfg.CryptoKey != "" {
		tlsCreds, err := credentials.NewServerTLSFromFile(cfg.CryptoCrt, cfg.CryptoKey)
		if err != nil {
			return nil, err
		}
		opt = append(opt, grpc.Creds(tlsCreds))
	}

	s = grpc.NewServer(opt...)
	pb.RegisterMetricServiceServer(s, &MetricServer{metrics: uc})
	d.s = s
	return d, nil
}

func (d *Delivery) Run(ctx context.Context, cfg *config.DeliveryConfig) error {
	serverCtx, stopServer := context.WithCancel(context.Background())
	defer stopServer()

	listen, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return err
	}

	go func() {
		<-ctx.Done()
		d.s.GracefulStop()
		log.Info("Delivery: graceful shutdown grpc server")
		stopServer()
	}()

	err = d.s.Serve(listen)
	if err != nil {
		return err
	}
	<-serverCtx.Done()
	return nil
}
