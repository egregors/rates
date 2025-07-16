package backends

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/egregors/rates"
)

type CurrencyAPI struct {
	url string
}

func (c *CurrencyAPI) Currencies() (map[string]string, error) {
	url := fmt.Sprintf("%ss.min.json", c.url[:len(c.url)-2])
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return nil, fmt.Errorf("worng URL")
		}
		return nil, fmt.Errorf("bad response: %s", resp.Status)
	}

	defer func(Body io.ReadCloser) { _ = Body.Close() }(resp.Body)

	var payload map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	res := make(map[string]string, len(payload))
	for k, v := range payload {
		res[k] = v.(string)
	}

	return res, nil
}

func NewCurrencyAPI(url string) *CurrencyAPI {
	return &CurrencyAPI{url: url}
}

func (c *CurrencyAPI) Rate(from, to string) (rates.Decimal, error) {
	url := fmt.Sprintf("%s%s.json", c.url, from)
	resp, err := http.Get(url)
	if err != nil {
		return rates.NewDecimal(0), err
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return rates.NewDecimal(0), fmt.Errorf("currency not found: %s -> %s", from, to)
		}
		return rates.NewDecimal(0), fmt.Errorf("bad response: %s", resp.Status)
	}

	defer func(Body io.ReadCloser) { _ = Body.Close() }(resp.Body)

	var payload map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return rates.NewDecimal(0), fmt.Errorf("failed to decode response: %w", err)
	}

	ratesData, ok := payload[from]
	if !ok {
		return rates.NewDecimal(0), fmt.Errorf("rate not found: %s -> %s", from, to)
	}

	rate, ok := ratesData.(map[string]any)[to]
	if !ok {
		return rates.NewDecimal(0), fmt.Errorf("rate not found: %s -> %s", from, to)
	}

	return rates.NewDecimal(rate.(float64)), nil
}
