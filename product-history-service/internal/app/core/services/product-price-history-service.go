package services

import (
	"context"
	"encoding/csv"
	"io"
	"log/slog"

	"golang.org/x/sync/errgroup"

	"github.com/juanpabloavilan/meli-interview-exercise/product-history-service/internal/app/core/models"
	filereader "github.com/juanpabloavilan/meli-interview-exercise/product-history-service/internal/app/pkg/file-reader"
	"github.com/juanpabloavilan/meli-interview-exercise/product-history-service/internal/app/ports"
)

const (
	maxBatchSize = 10000
)

type productPriceHistoryService struct {
	repo ports.ProductHistoryRepo
}

func NewProductPriceHistoryService(r ports.ProductHistoryRepo) ProductPriceHistoryService {
	return &productPriceHistoryService{
		repo: r,
	}
}

func (s *productPriceHistoryService) ImportFromCSVFile(ctx context.Context, reader io.Reader) error {
	var (
		rows         []models.ProducPriceHistory
		row          *models.ProducPriceHistory
		item         any
		csvReader    *csv.Reader
		unmarshaller filereader.CSVRowUnmarshaller[models.ProducPriceHistory]
		rowNumber    int64
		err          error
	)
	csvReader = csv.NewReader(reader)
	csvReader.ReuseRecord = true

	unmarshaller, err = filereader.NewCSVRowUnmarshaller[models.ProducPriceHistory](csvReader)
	if err != nil {
		return err
	}

	rowNumber = 0

	for {
		item, err = unmarshaller.ReadUnmarshalCSVRow()
		if err == io.EOF {
			slog.InfoContext(ctx, "END reading file", slog.Int64("numberOfRows", rowNumber))
			break
		}
		if err != nil {
			slog.ErrorContext(ctx, err.Error())
			return err
		}

		rowNumber++

		row = item.(*models.ProducPriceHistory)

		err = row.Validate()
		if err != nil {
			slog.ErrorContext(ctx, err.Error(), slog.Any("row", row))
			return err
		}

		rows = append(rows, *row)

	}

	if err := s.loadRowsInBatches(ctx, rows); err != nil {
		slog.ErrorContext(ctx, err.Error())
		return err
	}

	return nil

}

func (s *productPriceHistoryService) loadRowsInBatches(ctx context.Context, rows []models.ProducPriceHistory) error {
	if len(rows) < maxBatchSize {
		err := s.repo.AddMany(ctx, rows)
		if err != nil {
			return err
		}
		return nil
	}

	eg := new(errgroup.Group)
	intervals := s.getIntervals(len(rows))

	for i, interval := range intervals {
		eg.Go(func() error {
			start := interval[0]
			end := interval[1]
			slog.InfoContext(ctx, "BATCH LOAD", slog.Any("START", start), slog.Any("END", end), slog.Any("BATCH_ID", i))
			err := s.repo.AddMany(ctx, rows[start:end])
			if err != nil {
				return err
			}

			return nil

		})
	}

	if err := eg.Wait(); err != nil {
		slog.ErrorContext(ctx, err.Error())
	}

	return nil
}

func (s *productPriceHistoryService) getIntervals(length int) [][2]int {
	start := 0
	end := maxBatchSize
	intervals := make([][2]int, 0)

	numOfIntervals := length/maxBatchSize + 1

	for i := 0; i < numOfIntervals; i++ {
		intervals = append(intervals, [2]int{start, end})

		start = end
		end += maxBatchSize
		if end > length {
			end = length
		}
	}

	return intervals
}
