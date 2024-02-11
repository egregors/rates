package backends

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	CurrencyPairURL = "https://cdn.jsdelivr.net/gh/fawazahmed0/currency-api@1/latest/currencies/%s/%s.json"
)

type CurrencyAPI struct{}

func NewCurrencyAPI() *CurrencyAPI {
	return &CurrencyAPI{}
}

func (c *CurrencyAPI) Rate(from, to string) (float64, error) {
	url := fmt.Sprintf(CurrencyPairURL, from, to)
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return 0, fmt.Errorf("currency not found: %s -> %s", from, to)
		}
		return 0, fmt.Errorf("bad response: %s", resp.Status)
	}

	defer func(Body io.ReadCloser) { _ = Body.Close() }(resp.Body)

	var payload map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return 0, fmt.Errorf("failed to decode response: %w", err)
	}

	if _, ok := payload[to]; !ok {
		return 0, fmt.Errorf("rate not found: %s -> %s", from, to)
	}
	rate := map[string]float64{to: payload[to].(float64)}

	return rate[to], nil
}
