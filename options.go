package rates

import (
	"github.com/egregors/rates/lib"
	"github.com/egregors/rates/lib/cache"
)

type Options func(*Converter)

// WithLogger sets the logger for the rates conv
func WithLogger(l Logger) Options {
	// if no logger is provided, use a noop logger
	if l == nil {
		l = &lib.NoopLogger{}
	}

	return func(r *Converter) {
		r.l = l
	}
}

// WithCache sets the cache for the rates conv
func WithCache(c Cache[map[string]float64]) Options {
	// if no cache is provided, use a noop cache
	if c == nil {
		c = &cache.Noop[map[string]float64]{}
	}

	return func(r *Converter) {
		r.cache = c
	}
}
