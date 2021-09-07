package robinhood

import (
	"context"
	"time"
)

type CryptoHolding struct {
	AccountID string `json:"account_id"`
	CostBases []struct {
		CurrencyID        string  `json:"currency_id"`
		DirectCostBasis   float64 `json:"direct_cost_basis,string"`
		DirectQuantity    float64 `json:"direct_quantity,string"`
		ID                string  `json:"id"`
		IntradayCostBasis float64 `json:"intraday_cost_basis,string"`
		IntradayQuantity  float64 `json:"intraday_quantity,string"`
		MarkedCostBasis   float64 `json:"marked_cost_basis,string"`
		MarkedQuantity    float64 `json:"marked_quantity,string"`
	} `json:"cost_bases"`
	CreatedAt time.Time `json:"created_at"`
	Currency  struct {
		BrandColor string  `json:"brand_color"`
		Code       string  `json:"code"`
		ID         string  `json:"id"`
		Increment  float64 `json:"increment,string"`
		Name       string  `json:"name"`
		Type       string  `json:"type"`
	} `json:"currency"`
	ID                  string    `json:"id"`
	Quantity            float64   `json:"quantity,string"`
	QuantityAvailable   float64   `json:"quantity_available,string"`
	QuantityHeldForBuy  float64   `json:"quantity_held_for_buy,string"`
	QuantityHeldForSell float64   `json:"quantity_held_for_sell,string"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// GetCryptoHoldings returns crypto portfolio info
func (c *Client) GetCryptoHoldings(ctx context.Context) ([]CryptoHolding, error) {
	var p struct{ Results []CryptoHolding }
	var hdlUrl = EPCryptoHoldings
	err := c.GetAndDecode(ctx, hdlUrl, &p)
	return p.Results, err
}