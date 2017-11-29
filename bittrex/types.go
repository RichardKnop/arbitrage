package bittrex

// GetMarketsResponse ...
type GetMarketsResponse struct {
	Success bool      `json:"success"`
	Message string    `json:"message"`
	Result  []*Market `json:"result"`
}

// Market ...
type Market struct {
	MarketCurrency     string
	BaseCurrency       string
	MarketCurrencyLong string
	BaseCurrencyLong   string
	MinTradeSize       float64
	MarketName         string
	IsActive           bool
	Created            string
}

// GetCurrenciesResponse ...
type GetCurrenciesResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Result  []*Currency `json:"result"`
}

// Currency ...
type Currency struct {
	Currency        string
	CurrencyLong    string
	MinConfirmation int
	TxFee           float64
	IsActive        bool
	CoinType        string
	BaseAddress     *string
}

// GetTickerResponse ...
type GetTickerResponse struct {
	Success bool    `json:"success"`
	Message string  `json:"message"`
	Result  *Ticker `json:"result"`
}

// Ticker ...
type Ticker struct {
	Bid  float64
	Ask  float64
	Last float64
}
