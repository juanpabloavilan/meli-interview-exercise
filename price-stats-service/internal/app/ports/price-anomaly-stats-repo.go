package ports

import (
	"context"
)

type PriceAnomalyStatsRepo[T any] interface {
	UpdateStats(ctx context.Context, stats T) error
	GetStats(ctx context.Context, itemID string) (*T, error)
}
