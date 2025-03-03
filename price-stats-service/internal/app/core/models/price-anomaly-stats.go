package models

import "errors"

type AnomalyDetectionAlgo string

const (
	ZINDEX AnomalyDetectionAlgo = "ZINDEX"
)

var (
	ErrInvalidDetectionAlgo = errors.New("invalid detection algorithm")
)

func AnomalyDetectionAlgoFromString(v string) (AnomalyDetectionAlgo, error) {
	if v == "" || v == string(ZINDEX) {
		return ZINDEX, nil
	}

	return "", ErrInvalidDetectionAlgo
}

type ItemPriceHistory struct {
	ItemID       string
	PriceHistory []float64
}

type ItemPriceMongoEvent struct {
	ID struct {
		data string `bson:"_data"`
	} `bson:"_id"`
	FullDocument struct {
		ItemID string  `bson:"itemid"`
		Price  float64 `bson:"price"`
	} `bson:"fullDocument"`
}

type ItemPrice struct {
	ItemID string  `bson:"itemid"`
	Price  float64 `bson:"price"`
}

type ZindexPriceAnomalyStats struct {
	ItemID   string  `redis:"itemID"`
	Mean     float64 `redis:"mean"`
	StDev    float64 `redis:"stdev"`
	Variance float64 `redis:"variance"`
	N        float64 `redis:"n"`
}
