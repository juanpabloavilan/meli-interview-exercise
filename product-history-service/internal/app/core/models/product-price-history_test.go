package models

import "testing"

func TestProducPriceHistory_Validate(t *testing.T) {
	type fields struct {
		ItemID         string
		OrderCloseDate string
		Price          float64
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				ItemID:         "MLB3836655204",
				OrderCloseDate: "2024-02-18",
				Price:          114.1751,
			},
		},
		{
			name: "invalid item id",
			fields: fields{
				ItemID:         " ",
				OrderCloseDate: "2024-02-18",
				Price:          114.1751,
			},
			wantErr: true,
		},
		{
			name: "invalid close date",
			fields: fields{
				ItemID:         "MLB3836655204",
				OrderCloseDate: "invalid",
				Price:          114.1751,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := ProducPriceHistory{
				ItemID:         tt.fields.ItemID,
				OrderCloseDate: tt.fields.OrderCloseDate,
				Price:          tt.fields.Price,
			}
			if err := p.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("ProducPriceHistory.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
