package grpc

import (
	"fmt"

	"github.com/sreway/yametrics-v2/pkg/metric"
	pb "github.com/sreway/yametrics-v2/proto/metric/v1"
	domain "github.com/sreway/yametrics-v2/services/server/internal/domain/metric"
)

func NewMetric(in *pb.Metric) (*metric.Metric, error) {
	m := metric.Metric{
		ID: in.Id,
	}
	switch in.Type {
	case pb.Type_GAUGE:
		m.MType = metric.GaugeType
		m.Value = metric.NewGauge(in.Value)
	case pb.Type_COUNTER:
		m.MType = metric.CounterType
		m.Delta = metric.NewCounter(in.Delta)
	default:
		err := fmt.Errorf("%w: %d", domain.ErrInvalidMetricType, in.Type)
		return nil, domain.NewMetricErr(in.Id, err)
	}

	m.Hash = in.Hash

	return &m, nil
}

func NewProtobufMetric(m *metric.Metric) (*pb.Metric, error) {
	pbm := pb.Metric{
		Id: m.ID,
	}

	switch m.MType {
	case metric.GaugeType:
		pbm.Type = pb.Type_GAUGE
		pbm.Value = m.Value.Value()
	case metric.CounterType:
		pbm.Type = pb.Type_COUNTER
		pbm.Delta = m.Delta.Value()
	default:
		err := fmt.Errorf("%w: %s", domain.ErrInvalidMetricType, m.MType)
		return nil, domain.NewMetricErr(m.ID, err)
	}
	pbm.Hash = m.Hash

	return &pbm, nil
}
