package server

type Logger interface {
	Printf(format string, v ...interface{})
}

type Converter interface {
	Conv(amount float64, from, to string) (float64, error)
}
