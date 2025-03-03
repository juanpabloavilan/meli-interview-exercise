package controller

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/juanpabloavilan/meli-interview-exercise/price-stats-service/internal/app/core/models"
	"github.com/juanpabloavilan/meli-interview-exercise/price-stats-service/internal/app/core/services"
)

type PriceStatsController struct {
	router                  *gin.Engine
	anomalyDetectionService services.AnomalyDetectionService
}

type UpdateStatsRequest struct {
	ItemID  string `json:"itemID"`
	History []struct {
		ItemID         string  `json:"itemID" binding:"required"`
		OrderCloseDate string  `json:"ordeCloseDate"`
		Price          float64 `json:"price" binding:"required"`
	} `json:"history"`
}

type DetectAnomalyRequest struct {
	Query struct {
		ItemID string  `json:"itemID" binding:"required"`
		Price  float64 `json:"price" binding:"required"`
	} `json:"query" binding:"required"`
	Algorithm string `json:"algorithm"`
}

func NewPriceStatsController(r *gin.Engine, svc services.AnomalyDetectionService) *PriceStatsController {
	return &PriceStatsController{
		router:                  r,
		anomalyDetectionService: svc,
	}
}

func (ctrl PriceStatsController) SetRoutes(r *gin.Engine) {
	v1 := r.Group("api/v1")
	{
		v1.POST("product-stats/event/detect-price-anomaly", ctrl.DetectAnomaly)
		v1.PUT("product-stats/event/update-price-stats", ctrl.UpdateStats)
	}
}

func (ctrl PriceStatsController) UpdateStats(ctx *gin.Context) {
	var payload UpdateStatsRequest

	if err := ctx.BindJSON(payload); err != nil {
		slog.ErrorContext(ctx, err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"message":    err.Error(),
		})
	}

	// item := models.ItemPriceHistory{
	// 	ItemID:       payload.ItemID,
	// 	PriceHistory: nil,
	// }

	// ctrl.anomalyDetectionService.UpdateStats(ctx, payload)

}

func (ctrl PriceStatsController) DetectAnomaly(ctx *gin.Context) {
	var payload DetectAnomalyRequest

	if err := ctx.BindJSON(&payload); err != nil {
		slog.ErrorContext(ctx, err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"message":    "itemID is required",
		})
		return
	}

	detectionAlgo, err := models.AnomalyDetectionAlgoFromString(payload.Algorithm)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusNotFound,
			"message":    "failed to detect anomaly",
		})
		return
	}
	err = ctrl.anomalyDetectionService.SetStrategy(detectionAlgo)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusNotFound,
			"message":    "failed to detect anomaly",
		})
		return
	}

	isAnomaly, err := ctrl.anomalyDetectionService.DetectAnomaly(ctx, payload.Query.ItemID, payload.Query.Price)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusNotFound,
			"message":    "failed to detect anomaly",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"item_id":     payload.Query.ItemID,
		"price":       payload.Query.Price,
		"anomaly":     *isAnomaly,
		"metadata":    []string{},
		"status_code": http.StatusOK,
	})
}
