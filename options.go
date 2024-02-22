package rates

type Options func(*Converter)

// WithLogger sets the logger for the rates conv
func WithLogger(l Logger) Options {
	// if no logger is provided, use a noop logger, same as the no logger at all
	if l == nil {
		l = noopLogger{}
	}

	return func(r *Converter) {
		r.l = l
	}
}

type noopLogger struct{}

func (l noopLogger) Printf(string, ...interface{}) {}

// WithCache sets the cache for the rates conv
func WithCache(c Cache[map[string]float64]) Options {
	return func(r *Converter) {
		r.cache = c
	}
}
