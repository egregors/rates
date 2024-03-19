package server

//go:generate mockgen -source=deps.go -destination=../../mocks/deps.go -package=mocks

type Logger interface {
	Printf(format string, v ...interface{})
}

type Converter interface {
	Conv(amount float64, from, to string) (float64, error)
}
