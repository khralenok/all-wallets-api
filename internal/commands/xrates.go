package commands

import (
	"github.com/khralenok/all-wallets-api/internal/logic"
	"github.com/khralenok/all-wallets-api/internal/store"
)

// Worker function for updating data about current exchange rates in DB
func UpdateExchangeRates() error {

	rates, err := logic.FetchExchangeRates()

	if err != nil {
		return err
	}

	availableCurrencies, err := store.GetAvailableCurrencies()

	if err != nil {
		return err
	}

	updatedRates := logic.CalcExchangeRates(rates, availableCurrencies)

	err = store.AddUpdatedExchangeRates(updatedRates)

	if err != nil {
		return err
	}

	return nil
}
