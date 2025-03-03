package ports

import (
	"context"

	"github.com/juanpabloavilan/meli-interview-exercise/product-history-service/internal/app/core/models"
)

type UpdateStatsOpts struct {
	HistoryPerItem map[string][]models.ProducPriceHistory
}

type PriceStatsClient interface {
	UpdateStats(ctx context.Context, opts UpdateStatsOpts) error
}
