package main

import (
	"context"
	"net/http"

	httpcontroller "github.com/juanpabloavilan/meli-interview-exercise/product-history-service/internal/app/adapters/primary/http"
	httpclient "github.com/juanpabloavilan/meli-interview-exercise/product-history-service/internal/app/adapters/secondary/http"
	mongorepo "github.com/juanpabloavilan/meli-interview-exercise/product-history-service/internal/app/adapters/secondary/mongodb"
	"github.com/juanpabloavilan/meli-interview-exercise/product-history-service/internal/app/config"
	"github.com/juanpabloavilan/meli-interview-exercise/product-history-service/internal/app/core/services"
	"github.com/juanpabloavilan/meli-interview-exercise/product-history-service/internal/app/infrastructure"
)

func main() {
	ctx := context.Background()
	config := config.NewAppConfig()
	config.GetSecrets()

	db, err := infrastructure.NewMongoDB(ctx, config.MongoConnectionString, config.MongoDatabase)
	if err != nil {
		panic(err.Error())
	}

	router := infrastructure.NewRouter()

	priceHistoryRepo := mongorepo.NewProductPriceHistoryRepo(db)
	priceStatsHTTPClient := httpclient.NewProductPriceStatsHTTPClient(config.PriceStatsBaseURL, http.DefaultClient)

	service := services.NewProductPriceHistoryService(priceHistoryRepo, priceStatsHTTPClient)
	controller := httpcontroller.NewProductPriceController(router, service)

	controller.SetRoutes(router)
	infrastructure.Run(router, config.GinPort)

}
