package db

import (
	"context"
	"fmt"
	"log/slog"

	redis "github.com/redis/go-redis/v9"

	"github.com/juanpabloavilan/meli-interview-exercise/price-stats-service/internal/app/core/models"
	"github.com/juanpabloavilan/meli-interview-exercise/price-stats-service/internal/app/ports"
)

type zIndexPriceAnomalyRedisRepo struct {
	client *redis.Client
}

func NewZIndexPriceAnomalyRedisRepo(client *redis.Client) ports.PriceAnomalyStatsRepo[models.ZindexPriceAnomalyStats] {
	return &zIndexPriceAnomalyRedisRepo{
		client: client,
	}
}

func (r *zIndexPriceAnomalyRedisRepo) UpdateStats(ctx context.Context, stats models.ZindexPriceAnomalyStats) error {
	key := fmt.Sprintf("price-stats:%s", stats.ItemID)

	res := r.client.HSet(ctx, key, stats)
	if err := res.Err(); err != nil {
		slog.ErrorContext(ctx, err.Error())
		return err
	}

	return nil
}
func (r *zIndexPriceAnomalyRedisRepo) GetStats(ctx context.Context, itemID string) (*models.ZindexPriceAnomalyStats, error) {
	var stats models.ZindexPriceAnomalyStats

	key := fmt.Sprintf("price-stats:%s", itemID)
	err := r.client.HGetAll(ctx, key).Scan(&stats)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, err
	}

	return &stats, nil
}
