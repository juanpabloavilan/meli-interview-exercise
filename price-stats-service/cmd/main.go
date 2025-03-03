package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/juanpabloavilan/meli-interview-exercise/price-stats-service/internal/app/adapters/controller"
	"github.com/juanpabloavilan/meli-interview-exercise/price-stats-service/internal/app/adapters/db"
	"github.com/juanpabloavilan/meli-interview-exercise/price-stats-service/internal/app/adapters/stream"
	"github.com/juanpabloavilan/meli-interview-exercise/price-stats-service/internal/app/config"
	"github.com/juanpabloavilan/meli-interview-exercise/price-stats-service/internal/app/core/models"
	"github.com/juanpabloavilan/meli-interview-exercise/price-stats-service/internal/app/core/services"
	"github.com/juanpabloavilan/meli-interview-exercise/price-stats-service/internal/app/infrastructure"
	"github.com/juanpabloavilan/meli-interview-exercise/price-stats-service/pkg/logger"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config := config.NewAppConfig()
	config.GetSecrets()

	customHandler := &logger.Handler{
		Handler: slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
		}),
	}

	logger := slog.New(customHandler)
	slog.SetDefault(logger)

	router := infrastructure.NewRouter()
	redisClient, err := infrastructure.NewRedisClient(ctx, config.RedisConnectionString)
	if err != nil {
		panic(err)
	}

	mongoDB, err := infrastructure.NewMongoDB(ctx, config.MongoConnectionString, config.MongoDBName)
	if err != nil {
		panic(err)
	}

	zindexRepo := db.NewZIndexPriceAnomalyRedisRepo(redisClient)

	anomalyDetectionAlgos := map[models.AnomalyDetectionAlgo]services.AnomalyDetector{
		models.ZINDEX: services.NewZIndexService(zindexRepo, config.ZScoreTreshold),
	}

	anomalyDetectionSvc := services.NewPriceAnomalyDetectionService(anomalyDetectionAlgos)

	productPriceStream := stream.NewProductPricesStream(
		mongoDB.Collection("ProductPriceHistory"),
		config.ProductPriceStreamBatchSize,
		config.ProductPriceStreamBatchWindow,
		anomalyDetectionSvc,
	)

	go productPriceStream.WatchAndProcess(ctx)

	controller := controller.NewPriceStatsController(router, anomalyDetectionSvc)

	controller.SetRoutes(router)
	infrastructure.RunHTTPServer(router, config.GinPort)
}
