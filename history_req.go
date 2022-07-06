package robinhood

import (
	"context"
	"fmt"
)

const (
	TimeFrame3m = "3month"
	TimeFrame1m = "month"
	TimeFrame1y = "year"
)

type HistoryResponse struct {
	Title  interface{} `json:"title"`
	Weight interface{} `json:"weight"`
	Lines  []struct {
		Segments []struct {
			Points []struct {
				X          float64 `json:"x"`
				Y          float64 `json:"y"`
				CursorData struct {
					Label struct {
						Value string `json:"value"`
						Color struct {
							Light string `json:"light"`
							Dark  string `json:"dark"`
						} `json:"color"`
					} `json:"label"`
					PrimaryValue struct {
						Value string `json:"value"`
						Color struct {
							Light string `json:"light"`
							Dark  string `json:"dark"`
						} `json:"color"`
					} `json:"primary_value"`
					SecondaryValue struct {
						Main struct {
							Value string `json:"value"`
							Color struct {
								Light string `json:"light"`
								Dark  string `json:"dark"`
							} `json:"color"`
							Icon string `json:"icon"`
						} `json:"main"`
						StringFormat string      `json:"string_format"`
						Description  interface{} `json:"description"`
					} `json:"secondary_value"`
					TertiaryValue  interface{} `json:"tertiary_value"`
					PriceChartData struct {
						DollarValue struct {
							CurrencyCode string `json:"currency_code"`
							CurrencyID   string `json:"currency_id"`
							Amount       string `json:"amount"`
						} `json:"dollar_value"`
						DollarValueForReturn struct {
							CurrencyCode string `json:"currency_code"`
							CurrencyID   string `json:"currency_id"`
							Amount       string `json:"amount"`
						} `json:"dollar_value_for_return"`
						DollarValueForRateOfReturn struct {
							CurrencyCode string `json:"currency_code"`
							CurrencyID   string `json:"currency_id"`
							Amount       string `json:"amount"`
						} `json:"dollar_value_for_rate_of_return"`
					} `json:"price_chart_data"`
				} `json:"cursor_data"`
			} `json:"points"`
			Styles struct {
				Default struct {
					Color struct {
						Light string `json:"light"`
						Dark  string `json:"dark"`
					} `json:"color"`
					Opacity  float64 `json:"opacity"`
					LineType struct {
						Type        string  `json:"type"`
						StrokeWidth float64 `json:"stroke_width"`
						CapStyle    string  `json:"cap_style"`
					} `json:"line_type"`
				} `json:"default"`
				Active struct {
					Color struct {
						Light string `json:"light"`
						Dark  string `json:"dark"`
					} `json:"color"`
					Opacity  float64 `json:"opacity"`
					LineType struct {
						Type        string  `json:"type"`
						StrokeWidth float64 `json:"stroke_width"`
						CapStyle    string  `json:"cap_style"`
					} `json:"line_type"`
				} `json:"active"`
				Inactive struct {
					Color struct {
						Light string `json:"light"`
						Dark  string `json:"dark"`
					} `json:"color"`
					Opacity  float64 `json:"opacity"`
					LineType struct {
						Type        string  `json:"type"`
						StrokeWidth float64 `json:"stroke_width"`
						CapStyle    string  `json:"cap_style"`
					} `json:"line_type"`
				} `json:"inactive"`
			} `json:"styles"`
		} `json:"segments"`
		Direction string `json:"direction"`
		IsPrimary bool   `json:"is_primary"`
	} `json:"lines"`
	XAxis      interface{} `json:"x_axis"`
	YAxis      interface{} `json:"y_axis"`
	LegendData struct {
	} `json:"legend_data"`
	Fills          []interface{} `json:"fills"`
	ID             string        `json:"id"`
	DefaultDisplay struct {
		Label struct {
			Value string `json:"value"`
			Color struct {
				Light string `json:"light"`
				Dark  string `json:"dark"`
			} `json:"color"`
		} `json:"label"`
		PrimaryValue struct {
			Value string `json:"value"`
			Color struct {
				Light string `json:"light"`
				Dark  string `json:"dark"`
			} `json:"color"`
		} `json:"primary_value"`
		SecondaryValue struct {
			Main struct {
				Value string `json:"value"`
				Color struct {
					Light string `json:"light"`
					Dark  string `json:"dark"`
				} `json:"color"`
				Icon string `json:"icon"`
			} `json:"main"`
			StringFormat string `json:"string_format"`
			Description  struct {
				Value string `json:"value"`
				Color struct {
					Light string `json:"light"`
					Dark  string `json:"dark"`
				} `json:"color"`
			} `json:"description"`
		} `json:"secondary_value"`
		TertiaryValue  interface{} `json:"tertiary_value"`
		PriceChartData struct {
			DollarValue struct {
				CurrencyCode string `json:"currency_code"`
				CurrencyID   string `json:"currency_id"`
				Amount       string `json:"amount"`
			} `json:"dollar_value"`
			DollarValueForReturn struct {
				CurrencyCode string `json:"currency_code"`
				CurrencyID   string `json:"currency_id"`
				Amount       string `json:"amount"`
			} `json:"dollar_value_for_return"`
			DollarValueForRateOfReturn struct {
				CurrencyCode string `json:"currency_code"`
				CurrencyID   string `json:"currency_id"`
				Amount       string `json:"amount"`
			} `json:"dollar_value_for_rate_of_return"`
		} `json:"price_chart_data"`
	} `json:"default_display"`
	DisplaySpan   string `json:"display_span"`
	PageDirection string `json:"page_direction"`
	Overlays      []struct {
		SduiComponentType      string        `json:"sdui_component_type"`
		CurrentPlatform        interface{}   `json:"current_platform"`
		SkipCompatibilityCheck interface{}   `json:"skip_compatibility_check"`
		Content                []interface{} `json:"content"`
	} `json:"overlays"`
}

// Request account history data
func (c *Client) GetHistory(ctx context.Context, timeframe string) (*HistoryResponse, error) {

	// https: //bonfire.robinhood.com/portfolio/###/historical-chart/?display_span=###

	rsp := HistoryResponse{}

	url := EPHistPortfolio
	url = url + fmt.Sprintf("%s/historical-chart/?display_span=%s", c.Account.AccountNumber, timeframe)

	err := c.GetAndDecode(ctx, url, &rsp)
	return &rsp, err

}
