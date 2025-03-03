package ports

import (
	"context"
	"fmt"

	"github.com/juanpabloavilan/meli-interview-exercise/product-history-service/internal/app/core/models"
)

type ErrAddMany struct {
	FailedRows int
}

func (e ErrAddMany) Error() string {
	return fmt.Sprintf("failed to add many : failed rows %d", e.FailedRows)
}

type ProductHistoryRepo interface {
	AddMany(ctx context.Context, items []models.ProducPriceHistory) error
}
