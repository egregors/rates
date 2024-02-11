package conv

type Options func(*Converter)

// WithLogger sets the logger for the rates conv
func WithLogger(l Logger) Options {
	return func(r *Converter) {
		r.l = l
	}
}

// WithCache sets the cache for the rates conv
func WithCache(c Cache[map[string]float64]) Options {
	return func(r *Converter) {
		r.cache = c
	}
}
