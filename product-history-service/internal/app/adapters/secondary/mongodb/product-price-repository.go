package mongorepo

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/juanpabloavilan/meli-interview-exercise/product-history-service/internal/app/core/models"
	"github.com/juanpabloavilan/meli-interview-exercise/product-history-service/internal/app/ports"
)

const (
	ProducPriceHistoryCollection = "ProductPriceHistory"
)

type productPriceHistoryRepo struct {
	db *mongo.Database
}

func NewProductPriceHistoryRepo(db *mongo.Database) ports.ProductHistoryRepo {
	return &productPriceHistoryRepo{
		db: db,
	}
}

func (r *productPriceHistoryRepo) AddMany(ctx context.Context, items []models.ProducPriceHistory) error {
	itemsAny := make([]any, 0, len(items))
	for _, i := range items {
		itemsAny = append(itemsAny, i)
	}
	res, err := r.db.Collection(ProducPriceHistoryCollection).InsertMany(ctx, itemsAny, &options.InsertManyOptions{})
	if err != nil {
		slog.ErrorContext(ctx, "[productPriceHistoryRepo.AddMany] error", slog.Any("details", err))
		return ports.ErrAddMany{
			FailedRows: len(items),
		}
	}

	slog.InfoContext(ctx, "[productPriceHistoryRepo.AddMany] success", slog.Any("details", len(res.InsertedIDs)))

	return nil

}
