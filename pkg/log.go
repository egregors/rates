package pkg

// NoopLogger is a logger that does nothing.
type NoopLogger struct{}

func (l NoopLogger) Printf(string, ...interface{}) {}
