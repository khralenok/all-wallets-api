package logic

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

// Call OpenExchangeRates API to get exchange rates for USD. Return map with currencies and corresponding rates.
func FetchExchangeRates() (map[string]float64, error) {
	rates := make(map[string]float64)

	url := fmt.Sprintf("https://openexchangerates.org/api/latest.json?app_id=%s", os.Getenv("OXR_APP_ID"))

	data, err := getDataFromAPI(url)

	if err != nil {
		return map[string]float64{}, err
	}

	rawRates, ok := data["rates"].(map[string]any) // We need only rates

	if !ok {
		return map[string]float64{}, err
	}

	for key, value := range rawRates {
		assertedValue, ok := value.(float64)

		if !ok {
			continue
		}

		rates[key] = assertedValue
	}

	return rates, nil
}

func getDataFromAPI(url string) (map[string]any, error) {
	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, errors.New("can't read the data from API response")
	}

	var data map[string]any

	err = json.Unmarshal(body, &data)

	if err != nil {
		return nil, errors.New("can't parse response JSON")
	}

	return data, nil
}
