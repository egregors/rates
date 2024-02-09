package server

type Logger interface {
	Printf(format string, v ...interface{})
}

type RateProvider interface {
	GetRate(from, to string) (float64, error)
	GetCurrencyList() (map[string]string, error)
}
