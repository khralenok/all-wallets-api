package models

type CurrencyMetadata struct {
	Code          string `json:"code"`
	Name          string `json:"name"`
	Type          string `json:"type"`
	DecimalPlaces int    `json:"decimal_places"`
	Symbol        string `json:"symbol"`
}
