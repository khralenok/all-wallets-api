package models

import "time"

type ExchangeRate struct {
	FromCurrency string    `json:"from_currency"`
	ToCurrency   string    `json:"to_currency"`
	Rate         float64   `json:"rate"`
	FetchedAt    time.Time `json:"fetched_at"`
}

// Constructor function for exchange rate
func NewExchangeRate(from, to string, rate float64) ExchangeRate {
	return ExchangeRate{
		FromCurrency: from,
		ToCurrency:   to,
		Rate:         rate,
		FetchedAt:    time.Now(),
	}
}
