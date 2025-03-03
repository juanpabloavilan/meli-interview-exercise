package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/juanpabloavilan/meli-interview-exercise/price-stats-service/internal/app/core/models"
	"github.com/juanpabloavilan/meli-interview-exercise/price-stats-service/pkg/logger"
	"golang.org/x/sync/errgroup"
)

var (
	ErrUpdateStats          = errors.New("failed to update stats")
	ErrDetectAnomaly        = errors.New("failed to detect anomaly")
	ErrInvalidDetectionAlgo = errors.New("invalid anomaly detection algorithm")
)

type priceAnomalyDetectionService struct {
	strategy    AnomalyDetector
	strategyMap map[models.AnomalyDetectionAlgo]AnomalyDetector
}

func NewPriceAnomalyDetectionService(strategies map[models.AnomalyDetectionAlgo]AnomalyDetector) AnomalyDetectionService {
	return &priceAnomalyDetectionService{
		strategyMap: strategies,
	}
}

func (s *priceAnomalyDetectionService) SetStrategy(algorithm models.AnomalyDetectionAlgo) error {
	strategy, ok := s.strategyMap[algorithm]
	if !ok {
		return fmt.Errorf("%w strategy does not exists", ErrInvalidDetectionAlgo)
	}

	s.strategy = strategy
	return nil
}

func (s *priceAnomalyDetectionService) UpdateStats(ctx context.Context, item models.ItemPriceHistory) error {
	eg, ctx := errgroup.WithContext(ctx)
	ctx = logger.AppendCtx(ctx, slog.String("ITEM_ID", item.ItemID))

	for algorithm, anomalyDetector := range s.strategyMap {
		eg.Go(func() error {
			slog.InfoContext(ctx, "updating stats for "+string(algorithm))
			return anomalyDetector.UpdateStats(ctx, item)

		})
	}
	if err := eg.Wait(); err != nil {
		slog.ErrorContext(ctx, err.Error())
	}

	slog.InfoContext(ctx, "success updating stats")

	return nil
}

func (s *priceAnomalyDetectionService) DetectAnomaly(ctx context.Context, itemID string, price float64) (*bool, error) {
	if s.strategy == nil {
		return nil, fmt.Errorf("%w not selected strategy", ErrInvalidDetectionAlgo)
	}

	isAnomaly, err := s.strategy.DetectAnomaly(ctx, itemID, price)
	if err != nil {
		slog.ErrorContext(ctx, "failed to detect anomaly", "details", err.Error())
		return nil, ErrDetectAnomaly
	}

	slog.InfoContext(ctx, "success detecting anomaly")

	return isAnomaly, nil

}
