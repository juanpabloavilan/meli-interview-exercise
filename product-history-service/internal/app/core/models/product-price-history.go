package models

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrInvalidProductPriceHistory = errors.New("invalid product price record")
)

type ProducPriceHistory struct {
	ItemID         string  `csv:"ITEM_ID" json:"itemID"`
	OrderCloseDate string  `csv:"ORD_CLOSED_DT" json:"orderCloseDate"`
	Price          float64 `csv:"PRICE" json:"price"`
}

func (p ProducPriceHistory) Validate() error {
	if strings.TrimSpace(p.ItemID) == "" {
		return fmt.Errorf("%w, itemID is required", ErrInvalidProductPriceHistory)
	}
	if strings.TrimSpace(p.OrderCloseDate) == "" {
		return fmt.Errorf("%w, orderCloseDate is required", ErrInvalidProductPriceHistory)
	}
	if _, err := time.Parse(time.DateOnly, p.OrderCloseDate); err != nil {
		return fmt.Errorf("%w, orderDate is not a valid date", ErrInvalidProductPriceHistory)
	}

	return nil
}
