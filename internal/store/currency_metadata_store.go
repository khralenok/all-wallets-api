package store

import (
	"github.com/khralenok/all-wallets-api/internal/database"
	"github.com/khralenok/all-wallets-api/internal/models"
)

func GetAvailableCurrencies() ([]models.CurrencyMetadata, error) {
	var availableCurrencies []models.CurrencyMetadata

	query := "SELECT * FROM currency_metadata"

	rows, err := database.DB.Query(query)

	if err != nil {
		return []models.CurrencyMetadata{}, err
	}

	for rows.Next() {
		var nextCurrency models.CurrencyMetadata
		err := rows.Scan(&nextCurrency.Code, &nextCurrency.Name, &nextCurrency.Type, &nextCurrency.DecimalPlaces, &nextCurrency.Symbol)

		if err != nil {
			return []models.CurrencyMetadata{}, err
		}

		availableCurrencies = append(availableCurrencies, nextCurrency)
	}

	return availableCurrencies, nil
}

// Return decimal places for specific currency. Don't use it for definining wallet specific decimal places — instead use specialized function GetWalletDecimalPlaces.
func GetCurrencyDecimalPlaces(currency string) (int, error) {
	var decimalPlaces int

	query := "SELECT decimal_places FROM currency_metadata WHERE code = $1"

	err := database.DB.QueryRow(query, currency).Scan(&decimalPlaces)

	if err != nil {
		return -1, err
	}

	return decimalPlaces, nil
}
