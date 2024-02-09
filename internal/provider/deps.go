package provider

type Logger interface {
	Printf(format string, v ...interface{})
}
