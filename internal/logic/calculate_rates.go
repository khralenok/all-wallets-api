package logic

import (
	"github.com/khralenok/all-wallets-api/internal/models"
)

// STEP 3. Calculate X to USD exchange rates
// STEP 5. Calculate X to Y exchanhe rates (via USD)

func CalcExchangeRates(rates map[string]float64, availableCurrencies []models.CurrencyMetadata) []models.ExchangeRate {
	var calculatedRates []models.ExchangeRate

	for _, fromValue := range availableCurrencies {
		for _, toValue := range availableCurrencies {
			if toValue.Code == fromValue.Code {
				continue
			}

			rate := rates[toValue.Code] / rates[fromValue.Code]

			if fromValue.Code == "USD" {
				rate = rates[toValue.Code]
			}

			newExchangeRate := models.NewExchangeRate(fromValue.Code, toValue.Code, rate)
			calculatedRates = append(calculatedRates, newExchangeRate)
		}
	}

	return calculatedRates
}
