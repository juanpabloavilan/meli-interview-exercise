package stream

import (
	"context"
	"log/slog"
	"runtime"
	"time"

	"github.com/juanpabloavilan/meli-interview-exercise/price-stats-service/internal/app/core/models"
	"github.com/juanpabloavilan/meli-interview-exercise/price-stats-service/internal/app/core/services"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/sync/errgroup"
)

type ProductPricesStreamProcessor struct {
	collection *mongo.Collection
	// batchSize specifies a maximum amount of records to collect before triggering the flush
	batchSize int
	// batchWindow specifies a maximum amount of time to wait before triggering the flush
	batchWindow    time.Duration
	anomalyService services.AnomalyDetectionService
}

func NewProductPricesStream(collection *mongo.Collection, batchSize int, batchWindow time.Duration, anomalySvc services.AnomalyDetectionService) *ProductPricesStreamProcessor {
	return &ProductPricesStreamProcessor{
		collection:     collection,
		batchSize:      batchSize,
		batchWindow:    batchWindow,
		anomalyService: anomalySvc,
	}
}

func (s *ProductPricesStreamProcessor) WatchAndProcess(ctx context.Context) {
	stream, err := s.collection.Watch(ctx, mongo.Pipeline{
		{{
			Key: "$match", Value: bson.M{
				"operationType": "insert",
			},
		}},
		{{
			Key: "$project", Value: bson.M{
				"fullDocument.itemid": 1,
				"fullDocument.price":  1,
			},
		}},
	})
	if err != nil {
		panic(err)
	}

	go s.iterateChangeStream(ctx, stream)

}

func (s *ProductPricesStreamProcessor) iterateChangeStream(ctx context.Context, stream *mongo.ChangeStream) {
	var (
		events     = make(chan models.ItemPrice)
		itemPrice  models.ItemPrice
		mongoEvent models.ItemPriceMongoEvent
	)

	defer func() {
		close(events)
		if err := stream.Close(ctx); err != nil {
			slog.ErrorContext(ctx, "failed to close stream", slog.String("details", err.Error()))
		}
	}()

	go s.splitIntoBatches(ctx, events)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			for stream.Next(ctx) {
				if err := stream.Decode(&mongoEvent); err != nil {
					slog.ErrorContext(ctx, err.Error())
				}
				itemPrice.ItemID = mongoEvent.FullDocument.ItemID
				itemPrice.Price = mongoEvent.FullDocument.Price
				events <- itemPrice
			}
		}
	}
}

func (s *ProductPricesStreamProcessor) splitIntoBatches(ctx context.Context, eventsCh <-chan models.ItemPrice) {
	var (
		total               = 0
		processedBatchItems = 0
		rightPtr            = 0
		batchBuffer         = make([]models.ItemPrice, s.batchSize)
	)

	var event models.ItemPrice

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.Tick(s.batchWindow):
			// batch window time span is reached
			slog.DebugContext(ctx, "flushing window TICKER 1 SECOND", "batch", s.batchSize, "processed", processedBatchItems, "total", total)
			s.processBatch(ctx, batchBuffer[0:rightPtr])
			rightPtr, processedBatchItems = 0, 0

		case event = <-eventsCh:
			batchBuffer[rightPtr] = event

			if processedBatchItems == s.batchSize-1 {
				// batch window size is reached
				slog.DebugContext(ctx, "flushing window BATCH SIZE REACHED", "batch", s.batchSize, "processed", processedBatchItems, "total", total)
				s.processBatch(ctx, batchBuffer)
				rightPtr, processedBatchItems = 0, 0
			} else {
				rightPtr++
				processedBatchItems++
			}
			total++
		}
	}
}

func (s *ProductPricesStreamProcessor) processBatch(ctx context.Context, eventBatch []models.ItemPrice) {
	pricesGroupedByItemID := make(map[string][]float64)
	for _, e := range eventBatch {
		pricesGroupedByItemID[e.ItemID] = append(pricesGroupedByItemID[e.ItemID], e.Price)
	}

	eg, ctx := errgroup.WithContext(ctx)
	eg.SetLimit(runtime.NumCPU())

	for itemID, prices := range pricesGroupedByItemID {
		eg.Go(func() error {
			err := s.anomalyService.UpdateStats(ctx, models.ItemPriceHistory{
				ItemID:       itemID,
				PriceHistory: prices,
			})

			if err != nil {
				slog.ErrorContext(ctx, err.Error())
			}

			return nil
		})
	}

}
