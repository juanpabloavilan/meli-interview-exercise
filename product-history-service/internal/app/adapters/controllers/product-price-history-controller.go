package controllers

import (
	"log/slog"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/juanpabloavilan/meli-interview-exercise/product-history-service/internal/app/core/services"
)

type ProductPriceController struct {
	router                     *gin.Engine
	productPriceHistoryService services.ProductPriceHistoryService
}

func NewProductPriceController(r *gin.Engine, svc services.ProductPriceHistoryService) *ProductPriceController {
	return &ProductPriceController{
		router:                     r,
		productPriceHistoryService: svc,
	}
}

func (ctrl ProductPriceController) SetRoutes(r *gin.Engine) {
	v1 := r.Group("api/v1")
	{
		v1.POST("history/product/file-import", ctrl.ImportFromCSVFile)
	}
}

func (ctrl ProductPriceController) ImportFromCSVFile(ctx *gin.Context) {
	var payload struct {
		FileData *multipart.FileHeader `form:"file" binding:"required"`
	}

	if err := ctx.Bind(&payload); err != nil {
		slog.ErrorContext(ctx, "[ProductPriceController.ImportFromCSVFile] invalid payload", slog.Any("details", err))
		return
	}

	file, err := payload.FileData.Open()
	if err != nil {
		slog.ErrorContext(ctx, "failed to open file", slog.Any("details", err.Error()))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"message":    "invalid file",
		})
	}
	defer file.Close()

	if err := ctrl.productPriceHistoryService.ImportFromCSVFile(ctx, file); err != nil {
		slog.ErrorContext(ctx, "[ProductPriceController.ImportFromCSVFile] failed to import from csv file", slog.Any("details", err))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusInternalServerError,
			"message":    "Internal Server Error",
		})
	}

	ctx.JSON(http.StatusNoContent, nil)

}
