package rates

type Logger interface {
	Printf(format string, v ...interface{})
}

type Cache[T any] interface {
	Get(key string) (T, bool)
	Set(key string, value T)
	Len() int
	ToMap() map[string]T
}

type Source interface {
	Rate(from, to string) (Decimal, error)
	Currencies() (map[string]string, error)
}

type Converter interface {
	Conv(amount Decimal, from, to string) (Decimal, error)
	Currencies() map[string]string
}
