package bittrex

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

const (
	// APIHost is the domain name used for API endpoints
	APIHost = "https://bittrex.com/api/v1.1"
	// GetMarketsEndpoint is public endpoint to get open markets
	GetMarketsEndpoint = "/public/getmarkets"
	// GetCurrenciesEndpoint is a public endpoint to get traded currencies
	GetCurrenciesEndpoint = "/public/getcurrencies"
	// GetTickerEndpoint is a public endpoint to get tickers
	GetTickerEndpoint = "/public/getticker"
)

// GetMarkets ...
func (e *Exchange) GetMarkets() ([]*Market, error) {
	data, err := e.makeGetRequest(GetMarketsEndpoint)
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
	data, err := e.makeGetRequest(GetCurrenciesEndpoint)
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
	data, err := e.makeGetRequest(fmt.Sprintf("%s?market=%s", GetTickerEndpoint, market))
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

func (e *Exchange) makeGetRequest(path string) ([]byte, error) {
	resp, err := e.client.Get(e.cnf.Host + path)
	if err != nil {
		return []byte{}, err
	}

	return ioutil.ReadAll(resp.Body)
}
