package rates

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/egregors/rates/pkg/cache"
)

type Conv struct {
	pool       []Source
	strategy   Strategy
	currencies map[string]string

	l     Logger
	cache Cache[map[string]Decimal]
}

func New(providers []Source, opts ...Options) *Conv {
	r := &Conv{
		pool: providers,
	}

	defaultOpts := []Options{
		WithLogger(log.New(os.Stdout, "", log.LstdFlags)),
		WithCache(cache.NewInMem[map[string]Decimal](6 * time.Hour)),
		WithStrategy(Failover),
	}

	// init default options
	for _, o := range defaultOpts {
		o(r)
	}

	// override by custom options
	for _, o := range opts {
		o(r)
	}

	// init currencies list
	for _, s := range r.pool {
		if cs, err := s.Currencies(); err == nil {
			r.currencies = cs
			break
		}
	}

	return r
}

func (c *Conv) Conv(amount Decimal, from, to string) (Decimal, error) {
	// get from cache
	rFrom, hasFrom := c.cache.Get(from)
	if hasFrom {
		rTo, hasTo := rFrom[to]
		if hasTo {
			c.l.Printf("[INFO] hit rates cache: %s -> %s = %s", from, to, rTo.String())
			return amount.Mul(rTo), nil
		}
	}

	if !hasFrom {
		rFrom = make(map[string]Decimal)
	}

	// get from source pool
	switch c.strategy {
	// start from first and try to get rate from each source in pool until success
	case Failover:
		for _, s := range c.pool {
			r, err := s.Rate(from, to)
			if err != nil {
				c.l.Printf("[ERROR] can't get rate: %v", err)
				continue
			}

			rFrom[to] = r
			c.cache.Set(from, rFrom)

			c.l.Printf("[INFO] call rates source api: %s -> %s = %s", from, to, r.String())

			return amount.Mul(r), nil
		}

		return NewDecimal(0), fmt.Errorf("failed to get rate: %w", fmt.Errorf("all sources failed"))
	// just take first source from pool
	default:
		r, err := c.pool[0].Rate(from, to)
		if err != nil {
			c.l.Printf("[ERROR] can't get rate: %v", err)
			return NewDecimal(0), fmt.Errorf("failed to get rate: %w", err)
		}

		rFrom[to] = r
		c.cache.Set(from, rFrom)

		c.cache.Set(from, map[string]Decimal{to: r})
		c.l.Printf("[INFO] call rates source api: %s -> %s = %s", from, to, r.String())

		return amount.Mul(r), nil
	}
}

func (c *Conv) Currencies() map[string]string {
	return c.currencies
}
