package config

import (
	"os"
)

var (
	MONGO_CONN_STRING        = "MONGO_CONN_STRING"
	MONGO_DATABASE           = "MONGO_DATABASE"
	GIN_PORT                 = "GIN_PORT"
	PRICE_STATS_SVC_BASE_URL = "PRICE_STATS_SVC_BASE_URL"
)

type AppConfig struct {
	MongoConnectionString string
	MongoDatabase         string
	GinPort               string
	PriceStatsBaseURL     string
}

func NewAppConfig() *AppConfig {
	return &AppConfig{}
}

func (config *AppConfig) GetSecrets() {
	connString, ok := os.LookupEnv(MONGO_CONN_STRING)
	if !ok {
		panic(MONGO_CONN_STRING + " env variable not found ")
	}

	database, ok := os.LookupEnv(MONGO_DATABASE)
	if !ok {
		panic(MONGO_DATABASE + " env variable not found ")
	}

	ginPort, ok := os.LookupEnv(GIN_PORT)
	if !ok {
		panic(GIN_PORT + " env variable not found ")
	}

	priceStatsBaseURL, ok := os.LookupEnv(PRICE_STATS_SVC_BASE_URL)
	if !ok {
		panic(PRICE_STATS_SVC_BASE_URL + " env variable not found ")
	}

	config.MongoConnectionString = connString
	config.GinPort = ginPort
	config.MongoDatabase = database
	config.PriceStatsBaseURL = priceStatsBaseURL

}
