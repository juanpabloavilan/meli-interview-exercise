package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/juanpabloavilan/meli-interview-exercise/product-history-service/internal/app/adapters/controllers"
	"github.com/juanpabloavilan/meli-interview-exercise/product-history-service/internal/app/adapters/repositories"
	"github.com/juanpabloavilan/meli-interview-exercise/product-history-service/internal/app/config"
	"github.com/juanpabloavilan/meli-interview-exercise/product-history-service/internal/app/core/services"
	"github.com/juanpabloavilan/meli-interview-exercise/product-history-service/internal/app/infrastructure"
	"github.com/juanpabloavilan/meli-interview-exercise/product-history-service/internal/app/pkg/logger"
)

func main() {
	ctx := context.Background()
	config := config.NewAppConfig()
	config.GetSecrets()

	db, err := infrastructure.NewMongoDB(ctx, config.MongoConnectionString, config.MongoDatabase)
	if err != nil {
		panic(err.Error())
	}

	customHandler := &logger.Handler{
		Handler: slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
		}),
	}

	logger := slog.New(customHandler)
	slog.SetDefault(logger)

	router := infrastructure.NewRouter()

	priceHistoryRepo := repositories.NewProductPriceHistoryRepo(db)
	service := services.NewProductPriceHistoryService(priceHistoryRepo)
	controller := controllers.NewProductPriceController(router, service)

	controller.SetRoutes(router)
	infrastructure.Run(router, config.GinPort)

}
