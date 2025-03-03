package services

import (
	"context"
	"io"
)

type ProductPriceHistoryService interface {
	ImportFromCSVFile(ctx context.Context, reader io.Reader) error
}
