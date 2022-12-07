package delivery

import (
	"context"

	"github.com/sreway/yametrics-v2/services/server/config"
)

type Delivery interface {
	Run(ctx context.Context, cfg *config.DeliveryConfig) error
}
