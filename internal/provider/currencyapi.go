package provider

import (
	"encoding/json"
	"fmt"
	"github.com/egregors/rates/lib/cache"
	"io"
	"net/http"
)

const (
	CurrencyPairURL = "https://cdn.jsdelivr.net/gh/fawazahmed0/currency-api@1/latest/currencies/%s/%s.json"
	CurrencyListURL = "https://cdn.jsdelivr.net/gh/fawazahmed0/currency-api@1/latest/currencies.min.json"
)

// TODO: add TTL
// TODO: add cache invalidation
// TODO: add check if the rate is outdated
// TODO: add integration tests

type CurrencyAPI struct {
	inMemRates        Cache[map[string]float64]
	inMemCurrencyList Cache[string]

	l Logger
}

func NewCurrencyAPI(l Logger) *CurrencyAPI {
	return &CurrencyAPI{
		inMemRates:        cache.NewInMem[map[string]float64](),
		inMemCurrencyList: cache.NewInMem[string](),
		l:                 l,
	}
}

func (c *CurrencyAPI) GetRate(from, to string) (float64, error) {
	var cacheHasFrom bool
	if rate, ok := c.inMemRates.Get(from); ok {
		cacheHasFrom = true
		if rate, ok := rate[to]; ok {
			c.l.Printf("[INFO] from cache: %s -> %s = %f", from, to, rate)
			return rate, nil
		}
	}

	url := fmt.Sprintf(CurrencyPairURL, from, to)
	resp, err := http.Get(url)
	if err != nil {
		c.l.Printf("[ERROR] failed to get rate: %v", err)
		return 0, fmt.Errorf("failed to get rate: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		c.l.Printf("[ERROR] failed to get rate: %s", resp.Status)
		return 0, fmt.Errorf("failed to get rate: %s", resp.Status)
	}

	defer func(Body io.ReadCloser) { _ = Body.Close() }(resp.Body)

	var payload map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		c.l.Printf("[ERROR] failed to decode response: %v", err)
		return 0, fmt.Errorf("failed to decode response: %w", err)
	}

	if _, ok := payload[to]; !ok {
		c.l.Printf("[ERROR] rate not found: %s -> %s", from, to)
		return 0, fmt.Errorf("rate not found: %s -> %s", from, to)
	}
	rate := map[string]float64{to: payload[to].(float64)}

	if cacheHasFrom {
		curr, _ := c.inMemRates.Get(from)
		curr[to] = rate[to]
		rate = curr
	}

	c.inMemRates.Set(from, rate)
	c.l.Printf("[INFO] from api: %s -> %s = %f", from, to, rate[to])

	return rate[to], nil
}

func (c *CurrencyAPI) GetCurrencyList() (map[string]string, error) {
	if c.inMemCurrencyList.Len() > 0 {
		c.l.Printf("[INFO] from cache: currency list")
		return c.inMemCurrencyList.ToMap(), nil
	}

	resp, err := http.Get(CurrencyListURL)
	if err != nil {
		c.l.Printf("[ERROR] failed to get currency list: %v", err)
		return nil, fmt.Errorf("failed to get currency list: %w", err)
	}

	defer func(Body io.ReadCloser) { _ = Body.Close() }(resp.Body)

	var currencyList map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&currencyList); err != nil {
		c.l.Printf("[ERROR] failed to decode response: %v", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	for code, title := range currencyList {
		c.inMemCurrencyList.Set(code, title)
	}

	c.l.Printf("[INFO] from api: currency list")

	return c.inMemCurrencyList.ToMap(), nil
}
