package bittrex

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

// GetMarkets ...
func (e *Exchange) GetMarkets() ([]*Market, error) {
	resp, err := e.client.Get(e.host + GetMarketsEndpoint)
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
	resp, err := e.client.Get(e.host + GetCurrenciesEndpoint)
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
	resp, err := e.client.Get(fmt.Sprintf("%s?market=%s", e.host+GetTickerEndpoint, market))
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
