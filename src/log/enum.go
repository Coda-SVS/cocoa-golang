package log

type LogLevel int

const (
	TRACE LogLevel = iota + 1
	INFO
	WARNING
	ERROR
)

func (d LogLevel) String() string {
	var logLevels = [...]string{
		"TRACE",
		"INFO",
		"WARNING",
		"ERROR",
	}

	return logLevels[int(d)%len(logLevels)]
}
