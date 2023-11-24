package log

type LoggerOption struct {
	Prefix       string
	LogLevel     LogLevel
	PassToParent bool
}

func NewLoggerOption() LoggerOption {
	return LoggerOption{
		Prefix:       "",
		LogLevel:     TRACE,
		PassToParent: true,
	}
}
