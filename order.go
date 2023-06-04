package robinhood

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// OrderSide is which side of the trade an order is on
type OrderSide int

// MarshalJSON implements json.Marshaler
func (o OrderSide) MarshalJSON() ([]byte, error) {
	return []byte("\"" + strings.ToLower(o.String()) + "\""), nil
}

// Buy/Sell
//
//go:generate stringer -type OrderSide
const (
	Sell OrderSide = iota + 1
	Buy
)

// OrderType represents a Limit or Market order
type OrderType int

// MarshalJSON implements json.Marshaler
func (o OrderType) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", strings.ToLower(o.String()))), nil
}

// Well-known order types. Default is Market.
//
//go:generate stringer -type OrderType
const (
	Market OrderType = iota
	Limit
)

const (
	StopTrigger = "stop"
	ImmTrigger  = "immediate"

	TrailTypePrice = "price"
)

type TrailPeg struct {
	Price struct {
		Amount       float64 `json:"amount,string,omitempty"`
		CurrencyCode string  `json:"currency_code,omitempty"`
	} `json:"price,omitempty"`
	Type string `json:"type,omitempty"`
}

type RhOrder struct {
	Account       string    `json:"account,omitempty"`
	ExtendedHours bool      `json:"extended_hours"`
	Instrument    string    `json:"instrument,omitempty"`
	Price         float64   `json:"price,string,omitempty"`
	Quantity      float64   `json:"quantity,string,omitempty"`
	RefID         string    `json:"ref_id,omitempty"`
	Side          string    `json:"side,omitempty"`
	Symbol        string    `json:"symbol,omitempty"`
	TimeInForce   string    `json:"time_in_force,omitempty"`
	Trigger       string    `json:"trigger,omitempty"`
	Type          string    `json:"type,omitempty"`
	StopPrice     float64   `json:"stop_price,string,omitempty"`
	TrailingPeg   *TrailPeg `json:"trailing_peg,omitempty"`

	OverrideDayTradeChecks bool `json:"override_day_trade_checks,omitempty"`
	OverrideDtbpChecks     bool `json:"override_dtbp_checks,omitempty"`

	OrderFormVersion int `json:"order_form_version"`
}

func (c *Client) CreateOrder(i *Instrument) *RhOrder {

	newOrd := RhOrder{
		Account:          c.Account.URL,
		Symbol:           i.Symbol,
		Instrument:       i.URL,
		TimeInForce:      strings.ToLower(GTC.String()),
		Type:             strings.ToLower(Market.String()),
		Trigger:          ImmTrigger,
		OrderFormVersion: 2,
	}
	return &newOrd
}

// Order places an order for a given instrument. Cancellation of the given
// context cancels only the _http request_ and not any orders that may have
// been created regardless of the cancellation.
func (c *Client) SubmitOrder(ctx context.Context, rhOrd *RhOrder) (*OrderOutput, error) {

	// no fractional pennies.
	usdPrice := int(rhOrd.Price * 100)
	rhOrd.Price = float64(usdPrice) / 100.0
	rhOrd.Side = strings.ToLower(rhOrd.Side)
	rhOrd.Type = strings.ToLower(rhOrd.Type)
	rhOrd.TimeInForce = strings.ToLower(rhOrd.TimeInForce)

	bs, err := json.Marshal(rhOrd)
	if err != nil {
		return nil, err
	}

	post, err := http.NewRequest("POST", EPOrders, bytes.NewReader(bs))
	if err != nil {
		return nil, fmt.Errorf("error creating POST http.Request: %w", err)
	}

	post.Header.Add("Content-Type", "application/json")

	out := OrderOutput{}
	err = c.DoAndDecode(ctx, post, &out)
	if err != nil {
		return &out, err
	}

	return &out, nil
}

// OrderOutput is the response from the Order api
type OrderOutput struct {
	ID                 string    `json:"id"`
	RefID              string    `json:"ref_id"`
	URL                string    `json:"url"`
	Account            string    `json:"account"`
	Position           string    `json:"position"`
	CancelURL          string    `json:"cancel"`
	Instrument         string    `json:"instrument"`
	CumulativeQuantity float64   `json:"cumulative_quantity,string"`
	AveragePrice       float64   `json:"average_price,string"`
	Fees               float64   `json:"fees,string"`
	State              string    `json:"state"`
	Type               string    `json:"type"`
	Side               string    `json:"side"`
	TimeInForce        string    `json:"time_in_force"`
	Trigger            string    `json:"trigger"`
	Price              float64   `json:"price,string"`
	StopPrice          float64   `json:"stop_price,string"`
	Quantity           float64   `json:"quantity,string"`
	RejectReason       string    `json:"reject_reason"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	LastTransactionAt  time.Time `json:"last_transaction_at"`
	Executions         []struct {
		Price                  float64   `json:"price,string"`
		Quantity               float64   `json:"quantity,string"`
		SettlementDate         string    `json:"settlement_date"`
		Timestamp              time.Time `json:"timestamp"`
		ID                     string    `json:"id"`
		IpoAccessExecutionRank float64   `json:"ipo_access_execution_rank"`
	} `json:"executions"`
	ExtendedHours          bool      `json:"extended_hours"`
	OverrideDtbpChecks     bool      `json:"override_dtbp_checks"`
	OverrideDayTradeChecks bool      `json:"override_day_trade_checks"`
	ResponseCategory       string    `json:"response_category"`
	StopTriggeredAt        time.Time `json:"stop_triggered_at"`
	TrailingPeg            struct {
		Type       string `json:"type"`
		Percentage int    `json:"percentage"`
		Price      struct {
			Amount       float64 `json:"amount,string,omitempty"`
			CurrencyCode string  `json:"currency_code,omitempty"`
		} `json:"price,omitempty"`
	} `json:"trailing_peg"`
	LastTrailPrice struct {
		Amount       float64 `json:"amount,string"`
		CurrencyCode string  `json:"currency_code"`
		CurrencyID   string  `json:"currency_id"`
	} `json:"last_trail_price"`
	LastTrailPriceUpdatedAt time.Time `json:"last_trail_price_updated_at"`
	DollarBasedAmount       struct {
		Amount       float64 `json:"amount,string"`
		CurrencyCode string  `json:"currency_code"`
		CurrencyID   string  `json:"currency_id"`
	} `json:"dollar_based_amount"`
	TotalNotional struct {
		Amount       float64 `json:"amount,string"`
		CurrencyCode string  `json:"currency_code"`
		CurrencyID   string  `json:"currency_id"`
	} `json:"total_notional"`
	ExecutedNotional struct {
		Amount       float64 `json:"amount,string"`
		CurrencyCode string  `json:"currency_code"`
		CurrencyID   string  `json:"currency_id"`
	} `json:"executed_notional"`
	InvestmentScheduleID        string  `json:"investment_schedule_id"`
	IsIpoAccessOrder            bool    `json:"is_ipo_access_order"`
	IpoAccessCancellationReason string  `json:"ipo_access_cancellation_reason"`
	IpoAccessLowerCollaredPrice float64 `json:"ipo_access_lower_collared_price,string"`
	IpoAccessUpperCollaredPrice float64 `json:"ipo_access_upper_collared_price,string"`
	IpoAccessUpperPrice         float64 `json:"ipo_access_upper_price,string"`
	IpoAccessLowerPrice         float64 `json:"ipo_access_lower_price,string"`
	IsIpoAccessPriceFinalized   bool    `json:"is_ipo_access_price_finalized"`
}

// Update returns any errors and updates the item with any recent changes.
func (o *OrderOutput) Update(ctx context.Context, client *Client) error {
	return client.GetAndDecode(ctx, o.URL, o)
}

// Cancel attempts to cancel an odrer
func (o OrderOutput) Cancel(ctx context.Context, client *Client) error {
	post, err := http.NewRequest("POST", o.CancelURL, nil)
	if err != nil {
		return err
	}

	var o2 OrderOutput
	err = client.DoAndDecode(ctx, post, &o2)
	if err != nil {
		return errors.Wrap(err, "could not decode response")
	}

	if o2.RejectReason != "" {
		return errors.New(o2.RejectReason)
	}
	return nil
}

func (c *Client) CancelOrderById(ctx context.Context, id string) error {

	var ordUrl = EPOrders + id + "/cancel/"

	post, err := http.NewRequest("POST", ordUrl, nil)
	if err != nil {
		return err
	}

	var o2 OrderOutput
	err = c.DoAndDecode(ctx, post, &o2)
	if err != nil {
		return errors.Wrap(err, "could not decode response")
	}

	if o2.RejectReason != "" {
		return errors.New(o2.RejectReason)
	}
	return nil
}

// RecentOrders returns any recent orders made by this client.
func (c *Client) RecentOrders(ctx context.Context) ([]OrderOutput, error) {
	var o struct {
		Results []OrderOutput
	}
	err := c.GetAndDecode(ctx, EPOrders, &o)
	if err != nil {
		return o.Results, err
	}

	return o.Results, nil
}

func (c *Client) GetOrders(ctx context.Context, nextUrl *string, pgSize int64, stateFilter string) ([]OrderOutput, string, error) {
	var o struct {
		Results []OrderOutput
		Next    string
	}

	url := EPOrders
	if pgSize != 0 {
		url = url + fmt.Sprintf("?page_size=%d", pgSize)
		if len(stateFilter) > 0 {
			url = url + fmt.Sprintf("&state=%s", stateFilter)
		}
	} else {
		if len(stateFilter) > 0 {
			url = url + fmt.Sprintf("?state=%s", stateFilter)
		}
	}

	if nextUrl != nil {
		url = *nextUrl
	}

	err := c.GetAndDecode(ctx, url, &o)
	if err != nil {
		return o.Results, o.Next, err
	}

	return o.Results, o.Next, nil
}

// AllOrders returns all orders made by this client.
func (c *Client) AllOrders(ctx context.Context) ([]OrderOutput, error) {
	var o struct {
		Results []OrderOutput
	}

	url := EPOrders
	for {
		select {
		case <-ctx.Done():
			return o.Results, ctx.Err()
		default:
		}

		var tmp struct {
			Results []OrderOutput
			Next    string
		}
		err := c.GetAndDecode(ctx, url, &tmp)

		if err != nil {
			return o.Results, err
		}

		url = tmp.Next
		o.Results = append(o.Results, tmp.Results...)

		if url == "" {
			break
		}
	}

	return o.Results, nil
}
