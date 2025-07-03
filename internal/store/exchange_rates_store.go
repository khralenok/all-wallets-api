package store

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/khralenok/all-wallets-api/internal/database"
	"github.com/khralenok/all-wallets-api/internal/models"
)

func AddUpdatedExchangeRates(exchangeRates []models.ExchangeRate) error {
	if err := removePreviousRates(); err != nil {
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

// Return rate for specified currencies pair
func GetRate(from, to string, context *gin.Context) (float64, error) {
	var rate float64

	query := "SELECT rate FROM exchange_rates WHERE from_currency = $1 AND to_currency = $2"

	err := database.DB.QueryRow(query, from, to).Scan(&rate)

	if err != nil {
		return 0.0, err
	}

	return rate, nil
}

// Clean up exchange_rates table before it's update
func removePreviousRates() error {
	query := "DELETE FROM exchange_rates"

	_, err := database.DB.Exec(query)

	if err != nil {
		return err
	}

	return nil
}
