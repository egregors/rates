package conv

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/egregors/rates/lib/cache"
)

type Converter struct {
	rp RatesSource

	l     Logger
	cache Cache[map[string]float64]
}

func New(provider RatesSource, opts ...Options) *Converter {
	r := &Converter{
		rp: provider,
	}

	defaultOpts := []Options{
		WithLogger(log.New(os.Stdout, "", log.LstdFlags)),
		WithCache(cache.NewInMem[map[string]float64](6 * time.Hour)),
	}

	// init default options
	for _, o := range defaultOpts {
		o(r)
	}

	// override by custom options
	for _, o := range opts {
		o(r)
	}

	return r
}

func (c *Converter) Conv(amount float64, from, to string) (float64, error) {
	var cacheHasFrom bool
	if r, ok := c.cache.Get(from); ok {
		cacheHasFrom = true
		if r, ok := r[to]; ok {
			c.l.Printf("[INFO] hit rates cache: %s -> %s = %f", from, to, r)
			return amount * r, nil
		}
	}

	r, err := c.rp.Rate(from, to)
	if err != nil {
		c.l.Printf("[ERROR] can't get rate: %v", err)
		return 0, fmt.Errorf("failed to get rate: %w", err)
	}

	if cacheHasFrom {
		curr, _ := c.cache.Get(from)
		curr[to] = r
		c.cache.Set(from, curr)
	}

	c.cache.Set(from, map[string]float64{to: r})
	c.l.Printf("[INFO] call rates source api: %s -> %s = %f", from, to, r)

	return amount * r, nil
}
