package robinhood

import (
	"bytes"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"math"
	"strings"
	"time"

	"encoding/json"

	"net/http"
)

const (
	Tradable    = "tradable"
	NonTradable = "untradable"
)

// CryptoOrder is the payload to create a crypto currency order
type CryptoOrder struct {
	AccountID      string  `json:"account_id"`
	CurrencyPairID string  `json:"currency_pair_id"`
	Price          float64 `json:"price,string"`
	Quantity       float64 `json:"quantity,string"`
	RefID          string  `json:"ref_id"`
	Side           string  `json:"side"`
	TimeInForce    string  `json:"time_in_force"`
	Type           string  `json:"type"`

	AmountInDollars float64 `json:"-"`
}

// CryptoOrderOutput holds the response from api
type CryptoOrderOutput struct {
	AccountID          string    `json:"account_id"`
	AveragePrice       float64   `json:"average_price,string"`
	CancelURL          string    `json:"cancel_url"`
	CreatedAt          time.Time `json:"created_at"`
	CumulativeQuantity float64   `json:"cumulative_quantity,string"`
	CurrencyPairID     string    `json:"currency_pair_id"`
	EnteredPrice       float64   `json:"entered_price,string"`
	Executions         []struct {
		EffectivePrice float64   `json:"effective_price,string"`
		ID             string    `json:"id"`
		Quantity       float64   `json:"quantity,string"`
		Timestamp      time.Time `json:"timestamp"`
	} `json:"executions"`
	ID                      string      `json:"id"`
	InitiatorID             interface{} `json:"initiator_id"`
	InitiatorType           interface{} `json:"initiator_type"`
	LastTransactionAt       time.Time   `json:"last_transaction_at"`
	Price                   float64     `json:"price,string"`
	Quantity                float64     `json:"quantity,string"`
	RejectReason            string      `json:"reject_reason"`
	RefID                   string      `json:"ref_id"`
	RoundedExecutedNotional float64     `json:"rounded_executed_notional,string"`
	Side                    string      `json:"side"`
	State                   string      `json:"state"`
	StopPrice               float64     `json:"stop_price,string"`
	TimeInForce             string      `json:"time_in_force"`
	Type                    string      `json:"type"`
	UpdatedAt               time.Time   `json:"updated_at"`

	client *Client
}

func (c *Client) CreateCryptoOrder(currId string) *CryptoOrder {

	newOrd := CryptoOrder{
		AccountID:       c.CryptoAccount.ID,
		CurrencyPairID:  currId,
		Price:           0,
		Quantity:        0,
		TimeInForce:     strings.ToLower(GTC.String()),
		Type:            strings.ToLower(Market.String()),
		AmountInDollars: 0,
	}
	return &newOrd
}

// CryptoOrder will actually place the order
func (c *Client) SubmitCryptoOrder(ctx context.Context, o *CryptoOrder) (*CryptoOrderOutput, error) {

	if o.Quantity == 0 {
		o.Quantity = math.Round(o.AmountInDollars / o.Price)
	}
	payload, err := json.Marshal(o)

	if err != nil {
		return nil, err
	}

	post, err := http.NewRequest("POST", EPCryptoOrders, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("could not create Crypto http.Request: %w", err)
	}

	post.Header.Add("Content-Type", "application/json")

	var out CryptoOrderOutput
	err = c.DoAndDecode(ctx, post, &out)
	out.client = c
	return &out, err
}

// Cancel will cancel the order.
func (o *CryptoOrderOutput) Cancel(ctx context.Context) error {
	post, err := http.NewRequest("POST", o.CancelURL, nil)
	if err != nil {
		return err
	}

	var output CryptoOrderOutput
	err = o.client.DoAndDecode(ctx, post, &output)

	if err != nil {
		return errors.Wrap(err, "could not decode response")
	}

	if output.RejectReason != "" {
		return errors.New(output.RejectReason)
	}

	return nil
}

func (c *Client) GetCryptoOrders(ctx context.Context, nextUrl *string, pgSize int64) ([]CryptoOrderOutput, string, error) {
	var o struct {
		Results []CryptoOrderOutput
		Next    string
	}

	url := EPCryptoOrders
	if pgSize != 0 {
		url = url + fmt.Sprintf("?page_size=%d", pgSize)
	}
	if nextUrl != nil {
		url = *nextUrl
	}

	err := c.GetAndDecode(ctx, url, &o)
	if err != nil {
		return o.Results, o.Next, err
	}

	for i := range o.Results {
		o.Results[i].client = c
	}

	return o.Results, o.Next, nil
}

// Update returns any errors and updates the item with any recent changes.
func (o *CryptoOrderOutput) Update(ctx context.Context) error {
	ordUrl := EPCryptoOrders + fmt.Sprintf("%s", o.ID)
	return o.client.GetAndDecode(ctx, ordUrl, o)
}
