package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sreway/yametrics-v2/pkg/metric"
	log "github.com/sreway/yametrics-v2/pkg/tools/logger"
	pb "github.com/sreway/yametrics-v2/proto/metric/v1"
	domain "github.com/sreway/yametrics-v2/services/server/internal/domain/metric"
	"github.com/sreway/yametrics-v2/services/server/internal/usecases"
)

type MetricServer struct {
	pb.UnimplementedMetricServiceServer
	metrics usecases.Metric
}

func (s *MetricServer) Add(ctx context.Context, in *pb.AddMetricRequest) (*pb.AddMetricResponse, error) {
	var m *metric.Metric
	response := new(pb.AddMetricResponse)
	m, err := NewMetric(in.Metric)
	if err != nil {
		return response, HandelErrMetric(err)
	}

	err = s.metrics.Add(ctx, m)
	if err != nil {
		return response, HandelErrMetric(err)
	}

	return response, nil
}

func (s *MetricServer) BatchAdd(ctx context.Context, in *pb.BatchAddMetricRequest) (*pb.BatchAddMetricResponse, error) {
	response := new(pb.BatchAddMetricResponse)
	metrics := make([]*metric.Metric, 0, len(in.Metrics))
	for _, i := range in.Metrics {
		m, err := NewMetric(i)
		if err != nil {
			return response, HandelErrMetric(err)
		}
		metrics = append(metrics, m)
	}

	err := s.metrics.BatchAdd(ctx, metrics)
	if err != nil {
		return response, HandelErrMetric(err)
	}

	return response, nil
}

func (s *MetricServer) Get(ctx context.Context, in *pb.GetMetricRequest) (*pb.GetMetricResponse, error) {
	var mtype metric.Type
	response := new(pb.GetMetricResponse)

	m, err := s.metrics.Get(ctx, in.Id, mtype)
	if err != nil {
		return response, HandelErrMetric(err)
	}
	pbm, err := NewProtobufMetric(m)
	if err != nil {
		return response, HandelErrMetric(err)
	}
	response.Metric = pbm
	return response, nil
}

func (s *MetricServer) GetMany(ctx context.Context, _ *pb.GetManyMetricRequest) (*pb.GetManyMetricResponse, error) {
	response := new(pb.GetManyMetricResponse)
	metrics, err := s.metrics.GetMany(ctx)
	if err != nil {
		return response, HandelErrMetric(err)
	}

	pbmetrics := make([]*pb.Metric, 0, len(metrics))
	for _, i := range metrics {
		var pbm *pb.Metric
		pbm, err = NewProtobufMetric(&i)
		if err != nil {
			return response, HandelErrMetric(err)
		}
		pbmetrics = append(pbmetrics, pbm)
	}

	response.Metrics = pbmetrics
	return response, nil
}

func (s *MetricServer) StorageCheck(ctx context.Context, _ *pb.StorageCheckMetricRequest) (
	*pb.StorageCheckMetricResponse, error,
) {
	response := new(pb.StorageCheckMetricResponse)
	err := s.metrics.StorageCheck(ctx)
	if err != nil {
		log.Error(err.Error())
		return response, status.Error(codes.Unavailable, err.Error())
	}
	return response, nil
}

func HandelErrMetric(err error) error {
	log.Error(err.Error())

	switch {
	case errors.Is(err, domain.ErrInvalidMetricHash):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, domain.ErrInvalidMetricType):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, domain.ErrInvalidMetricValue):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, domain.ErrMetricNotFound):
		return status.Error(codes.NotFound, err.Error())
	default:
		return status.Error(codes.Unknown, err.Error())
	}
}
