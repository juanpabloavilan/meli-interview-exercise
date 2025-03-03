package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"runtime"

	"github.com/juanpabloavilan/meli-interview-exercise/product-history-service/internal/app/core/models"
	"github.com/juanpabloavilan/meli-interview-exercise/product-history-service/internal/app/ports"
	"golang.org/x/sync/errgroup"
)

type UpdateStatsBody struct {
	ItemID  string                      `json:"itemID"`
	History []models.ProducPriceHistory `json:"history"`
}

type productPriceStatsHTTPClient struct {
	client  *http.Client
	baseURL string
}

func NewProductPriceStatsHTTPClient(priceStatsURL string, client *http.Client) ports.PriceStatsClient {
	return &productPriceStatsHTTPClient{
		baseURL: priceStatsURL,
		client:  client,
	}
}

func (s productPriceStatsHTTPClient) UpdateStats(ctx context.Context, opts ports.UpdateStatsOpts) error {
	url := fmt.Sprintf("%s/product/:id", s.baseURL)

	eg := new(errgroup.Group)
	eg.SetLimit(runtime.NumCPU() * 2)

	for itemID, history := range opts.HistoryPerItem {
		eg.Go(func() error {
			body := UpdateStatsBody{
				ItemID:  itemID,
				History: history,
			}

			jsonBytes, err := json.Marshal(body)
			if err != nil {
				slog.ErrorContext(ctx, err.Error())
				return err
			}

			req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewReader(jsonBytes))
			if err != nil {
				slog.ErrorContext(ctx, err.Error())
				return err
			}
			req.Header.Set("Content-Type", "application/json")

			resp, err := s.client.Do(req)
			if err != nil {
				slog.ErrorContext(ctx, err.Error())
				return err
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				err = fmt.Errorf("failed to call price stats microservice")
				slog.ErrorContext(ctx, err.Error(), slog.Any("details", body))

				return err
			}

			slog.InfoContext(ctx, "success updating price stats", slog.Any("body", body))

			return nil
		})
	}

	return nil
}
