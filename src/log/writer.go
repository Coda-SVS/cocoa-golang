package log

import (
	"log"
)

type LogWriter struct {
	loggerTrace   *log.Logger
	loggerInfo    *log.Logger
	loggerWarning *log.Logger
	loggerError   *log.Logger
}

func NewLogWriter(
	loggerTrace *log.Logger,
	loggerInfo *log.Logger,
	loggerWarning *log.Logger,
	loggerError *log.Logger,
) LogWriter {
	return LogWriter{
		loggerTrace:   loggerTrace,
		loggerInfo:    loggerInfo,
		loggerWarning: loggerWarning,
		loggerError:   loggerError,
	}
}
