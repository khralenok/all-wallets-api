package store

import (
	"fmt"
	"strings"

	"github.com/khralenok/all-wallets-api/internal/database"
	"github.com/khralenok/all-wallets-api/internal/models"
)

func AddUpdatedExchangeRates(exchangeRates []models.ExchangeRate) error {
	if err := deletePreviousRates(); err != nil {
		return err
	}

	query := "INSERT INTO exchange_rates(from_currency, to_currency, rate, fetched_at) VALUES "
	args := []any{}
	valueStrings := []string{}

	for i, value := range exchangeRates {
		placeholder := i*4 + 1

		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d)", placeholder, placeholder+1, placeholder+2, placeholder+3))
		args = append(args, value.FromCurrency, value.ToCurrency, value.Rate, value.FetchedAt)

	}

	query += strings.Join(valueStrings, ",")

	_, err := database.DB.Exec(query, args...)

	if err != nil {
		return err
	}

	return nil
}

func deletePreviousRates() error {
	query := "DELETE FROM exchange_rates"

	_, err := database.DB.Exec(query)

	if err != nil {
		return err
	}

	return nil
}
