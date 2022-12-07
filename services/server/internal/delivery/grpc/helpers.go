package grpc

import (
	"fmt"

	"github.com/sreway/yametrics-v2/pkg/metric"
	pb "github.com/sreway/yametrics-v2/proto/metric/v1"
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
		return nil, fmt.Errorf("unknown metric type")
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
		return nil, fmt.Errorf("unknown metric type")
	}
	pbm.Hash = m.Hash

	return &pbm, nil
}
