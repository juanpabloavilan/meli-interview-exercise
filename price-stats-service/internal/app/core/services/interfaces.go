package services

import (
	"context"

	"github.com/juanpabloavilan/meli-interview-exercise/price-stats-service/internal/app/core/models"
)

type AnomalyDetector interface {
	UpdateStats(ctx context.Context, item models.ItemPriceHistory) error
	DetectAnomaly(ctx context.Context, itemID string, price float64) (*bool, error)
}
