package grpc

import (
	"context"
	"fmt"

	"github.com/sreway/yametrics-v2/pkg/metric"
	log "github.com/sreway/yametrics-v2/pkg/tools/logger"
	pb "github.com/sreway/yametrics-v2/proto/metric/v1"
	"github.com/sreway/yametrics-v2/services/server/internal/usecases"
)

type MetricServer struct {
	pb.UnimplementedMetricServiceServer
	metrics usecases.Metric
}

func (s *MetricServer) Add(ctx context.Context, in *pb.AddMetricRequest) (*pb.AddMetricResponse, error) {
	var m *metric.Metric
	response := new(pb.AddMetricResponse)

	switch in.Metric.Type {
	case pb.Type_GAUGE:
		m = metric.New(in.Metric.Id, metric.GaugeType, in.Metric.Value)
	case pb.Type_COUNTER:
		m = metric.New(in.Metric.Id, metric.CounterType, in.Metric.Delta)
	default:
		msg := "unknown metric type"
		response.Error = msg
		log.Error(msg)
		return response, nil
	}

	err := s.metrics.Add(ctx, m)
	if err != nil {
		log.Error(err.Error())
		response.Error = err.Error()
	}

	return response, err
}

func (s *MetricServer) BatchAdd(ctx context.Context, in *pb.BatchAddMetricRequest) (*pb.BatchAddMetricResponse, error) {
	response := new(pb.BatchAddMetricResponse)
	metrics := make([]*metric.Metric, 0, len(in.Metrics))
	for _, i := range in.Metrics {
		m, err := NewMetric(i)
		if err != nil {
			response.Error = err.Error()
			log.Error(err.Error())
			return response, err
		}
		metrics = append(metrics, m)
	}

	err := s.metrics.BatchAdd(ctx, metrics)
	if err != nil {
		response.Error = err.Error()
		log.Error(err.Error())
		return response, err
	}

	return response, nil
}

func (s *MetricServer) Get(ctx context.Context, in *pb.GetMetricRequest) (*pb.GetMetricResponse, error) {
	var mtype metric.Type
	response := new(pb.GetMetricResponse)
	switch in.Type {
	case pb.Type_GAUGE:
		mtype = metric.GaugeType
	case pb.Type_COUNTER:
		mtype = metric.CounterType
	default:
		response.Error = "unknown metric type"
		log.Error(response.Error)
		return response, fmt.Errorf(response.Error)
	}

	m, err := s.metrics.Get(ctx, in.Id, mtype)
	if err != nil {
		response.Error = err.Error()
		log.Error(err.Error())
		return response, err
	}

	pbm, err := NewProtobufMetric(m)
	if err != nil {
		response.Error = err.Error()
		log.Error(err.Error())
		return response, err
	}
	response.Metric = pbm
	return response, nil
}

func (s *MetricServer) GetMany(ctx context.Context, _ *pb.GetManyMetricRequest) (*pb.GetManyMetricResponse, error) {
	response := new(pb.GetManyMetricResponse)
	metrics, err := s.metrics.GetMany(ctx)
	if err != nil {
		response.Error = err.Error()
		log.Error(err.Error())
		return response, err
	}

	pbmetrics := make([]*pb.Metric, 0, len(metrics))
	for _, i := range metrics {
		var pbm *pb.Metric
		pbm, err = NewProtobufMetric(&i)
		if err != nil {
			response.Error = err.Error()
			log.Error(err.Error())
			return response, err
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
		response.Error = err.Error()
		log.Error(err.Error())
		return response, err
	}
	return response, nil
}
