package bittrex

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

const (
	// BaseHost will be used by default
	BaseHost = "https://bittrex.com/api/v1.1"
	// GetMarketsEndpoint is public endpoint to get open markets
	GetMarketsEndpoint = "/public/getmarkets"
	// GetCurrenciesEndpoint is a public endpoint to get traded currencies
	GetCurrenciesEndpoint = "/public/getcurrencies"
	// GetTickerEndpoint is a public endpoint to get tickers
	GetTickerEndpoint = "/public/getticker"
)

// GetMarkets ...
func (e *Exchange) GetMarkets() ([]*Market, error) {
	resp, err := e.client.Get(e.cnf.Host + GetMarketsEndpoint)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	response := new(GetMarketsResponse)
	if err := json.Unmarshal(data, response); err != nil {
		return nil, err
	}

	if !response.Success {
		return nil, errors.New(response.Message)
	}

	return response.Result, nil
}

// GetCurrencies ...
func (e *Exchange) GetCurrencies() ([]*Currency, error) {
	resp, err := e.client.Get(e.cnf.Host + GetCurrenciesEndpoint)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	response := new(GetCurrenciesResponse)
	if err := json.Unmarshal(data, response); err != nil {
		return nil, err
	}

	if !response.Success {
		return nil, errors.New(response.Message)
	}

	return response.Result, nil
}

// GetTicker ...
func (e *Exchange) GetTicker(market string) (*Ticker, error) {
	resp, err := e.client.Get(fmt.Sprintf("%s?market=%s", e.cnf.Host+GetTickerEndpoint, market))
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	response := new(GetTickerResponse)
	if err := json.Unmarshal(data, response); err != nil {
		return nil, err
	}

	if !response.Success {
		return nil, errors.New(response.Message)
	}

	if response.Result == nil {
		return nil, ErrEmptyResult
	}

	return response.Result, nil
}
