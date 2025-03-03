package config

import (
	"os"
	"strconv"
	"time"
)

var (
	REDIS_CONN_STRING                 = "REDIS_CONN_STRING"
	MONGO_CONN_STRING                 = "MONGO_CONN_STRING"
	GIN_PORT                          = "GIN_PORT"
	PRODUCT_PRICE_STREAM_BATCH_SIZE   = "PRODUCT_PRICE_STREAM_BATCH_SIZE"
	PRODUCT_PRICE_STREAM_BATCH_WINDOW = "PRODUCT_PRICE_STREAM_BATCH_WINDOW"
	Z_SCORE_THRESHOLD                 = "Z_SCORE_THRESHOLD"
	MONGO_DB_NAME                     = "MONGO_DB_NAME"
)

type AppConfig struct {
	MongoConnectionString         string
	MongoDBName                   string
	RedisConnectionString         string
	GinPort                       string
	ProductPriceStreamBatchSize   int
	ProductPriceStreamBatchWindow time.Duration
	ZScoreTreshold                int
}

func NewAppConfig() *AppConfig {
	return &AppConfig{}
}

func (config *AppConfig) GetSecrets() {
	redisConnString, ok := os.LookupEnv(REDIS_CONN_STRING)
	if !ok {
		panic(REDIS_CONN_STRING + " env variable not found ")
	}

	ginPort, ok := os.LookupEnv(GIN_PORT)
	if !ok {
		panic(GIN_PORT + " env variable not found ")
	}

	mongoConnString, ok := os.LookupEnv(MONGO_CONN_STRING)
	if !ok {
		panic(MONGO_CONN_STRING + " env variable not found ")
	}

	batchSizeStr, ok := os.LookupEnv(PRODUCT_PRICE_STREAM_BATCH_SIZE)
	if !ok {
		panic(PRODUCT_PRICE_STREAM_BATCH_SIZE + " env variable not found ")
	}
	batchSize, err := strconv.Atoi(batchSizeStr)
	if err != nil {
		panic(PRODUCT_PRICE_STREAM_BATCH_SIZE + err.Error())
	}

	batchWindowStr, ok := os.LookupEnv(PRODUCT_PRICE_STREAM_BATCH_WINDOW)
	if !ok {
		panic(PRODUCT_PRICE_STREAM_BATCH_WINDOW + " env variable not found ")
	}
	batchWindow, err := time.ParseDuration(batchWindowStr)
	if err != nil {
		panic(PRODUCT_PRICE_STREAM_BATCH_WINDOW + err.Error())
	}

	zScoreTresholdStr, ok := os.LookupEnv(Z_SCORE_THRESHOLD)
	if !ok {
		panic(Z_SCORE_THRESHOLD + " env variable not found ")
	}
	zScoreTreshold, err := strconv.Atoi(zScoreTresholdStr)
	if err != nil {
		panic(Z_SCORE_THRESHOLD + err.Error())
	}
	dbName, ok := os.LookupEnv(MONGO_DB_NAME)
	if !ok {
		panic(MONGO_DB_NAME + " env variable not found ")
	}

	config.RedisConnectionString = redisConnString
	config.GinPort = ginPort
	config.MongoConnectionString = mongoConnString
	config.ProductPriceStreamBatchSize = batchSize
	config.ProductPriceStreamBatchWindow = batchWindow
	config.ZScoreTreshold = zScoreTreshold
	config.MongoDBName = dbName

}
