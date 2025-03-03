package services

import (
	"context"
	"fmt"
	"log/slog"
	"math"

	"github.com/aws/smithy-go/ptr"
	"github.com/juanpabloavilan/meli-interview-exercise/price-stats-service/internal/app/core/models"
	"github.com/juanpabloavilan/meli-interview-exercise/price-stats-service/internal/app/ports"
	"github.com/juanpabloavilan/meli-interview-exercise/price-stats-service/pkg/logger"
)

type ZIndex struct {
	repo      ports.PriceAnomalyStatsRepo[models.ZindexPriceAnomalyStats]
	threshold int
}

func NewZIndexService(r ports.PriceAnomalyStatsRepo[models.ZindexPriceAnomalyStats], zScoreThreshold int) AnomalyDetector {
	return &ZIndex{
		repo:      r,
		threshold: zScoreThreshold,
	}
}

func (s *ZIndex) UpdateStats(ctx context.Context, item models.ItemPriceHistory) error {
	ctx = logger.AppendCtx(ctx, slog.String("algorithm", string(models.ZINDEX)))

	// prevStats refers to previous stats
	prevStats, err := s.repo.GetStats(ctx, item.ItemID)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
	}

	prev := *prevStats

	slog.InfoContext(ctx, "PREV stats", slog.Any("stats", prev))
	// current refers to current stats
	current := models.ZindexPriceAnomalyStats{
		ItemID: item.ItemID,
	}

	var t1, t2, t3, t4 float64
	for _, x := range item.PriceHistory {
		current.N++
		current.Mean = prev.Mean + (x-prev.Mean)/current.N
		if current.N > 1 {
			t1 = (current.N - 2) * prev.Variance
			t2 = (current.N - 1) * math.Pow(prev.Mean-current.Mean, 2)
			t3 = math.Pow((x - current.Mean), 2)
			t4 = current.N - 1
			current.Variance = (t1 + t2 + t3) / t4
			current.StDev = math.Sqrt(current.Variance)
		}
		prev = current
	}

	if err = s.repo.UpdateStats(ctx, current); err != nil {
		return err
	}

	slog.InfoContext(ctx, "update stats correctly")

	return nil
}

func (s *ZIndex) DetectAnomaly(ctx context.Context, itemID string, price float64) (*bool, error) {
	stats, err := s.repo.GetStats(ctx, itemID)
	if err != nil {
		return nil, fmt.Errorf("%w failed to get stats", ErrDetectAnomaly)
	}

	zIndex := (price - stats.Mean) / stats.StDev
	if math.Abs(zIndex) > float64(s.threshold) {
		return ptr.Bool(true), nil
	}

	return ptr.Bool(false), nil

}
