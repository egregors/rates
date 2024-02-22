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
	Rate(from, to string) (float64, error)
}
