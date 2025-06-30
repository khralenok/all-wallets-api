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
