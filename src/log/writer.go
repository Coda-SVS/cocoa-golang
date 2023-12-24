package log

import (
	"log"
)

type LowLogWriter struct {
	loggerTrace   *log.Logger
	loggerInfo    *log.Logger
	loggerWarning *log.Logger
	loggerError   *log.Logger
}

type LogWriter struct {
	fileLogWriter    *LowLogWriter
	consoleLogWriter *LowLogWriter
}

func NewLogWriter(
	fileLoggerTrace *log.Logger,
	fileLoggerInfo *log.Logger,
	fileLoggerWarning *log.Logger,
	fileLoggerError *log.Logger,
	consoleLoggerTrace *log.Logger,
	consoleLoggerInfo *log.Logger,
	consoleLoggerWarning *log.Logger,
	consoleLoggerError *log.Logger,
) *LogWriter {
	lw := &LogWriter{}

	lw.fileLogWriter = &LowLogWriter{
		loggerTrace:   fileLoggerTrace,
		loggerInfo:    fileLoggerInfo,
		loggerWarning: fileLoggerWarning,
		loggerError:   fileLoggerError,
	}

	lw.consoleLogWriter = &LowLogWriter{
		loggerTrace:   consoleLoggerTrace,
		loggerInfo:    consoleLoggerInfo,
		loggerWarning: consoleLoggerWarning,
		loggerError:   consoleLoggerError,
	}

	return lw
}

func NewFileLogWriter(
	fileLoggerTrace *log.Logger,
	fileLoggerInfo *log.Logger,
	fileLoggerWarning *log.Logger,
	fileLoggerError *log.Logger,
) *LogWriter {
	lw := &LogWriter{}

	lw.fileLogWriter = &LowLogWriter{
		loggerTrace:   fileLoggerTrace,
		loggerInfo:    fileLoggerInfo,
		loggerWarning: fileLoggerWarning,
		loggerError:   fileLoggerError,
	}

	return lw
}

func NewConsoleLogWriter(
	consoleLoggerTrace *log.Logger,
	consoleLoggerInfo *log.Logger,
	consoleLoggerWarning *log.Logger,
	consoleLoggerError *log.Logger,
) *LogWriter {
	lw := &LogWriter{}

	lw.consoleLogWriter = &LowLogWriter{
		loggerTrace:   consoleLoggerTrace,
		loggerInfo:    consoleLoggerInfo,
		loggerWarning: consoleLoggerWarning,
		loggerError:   consoleLoggerError,
	}

	return lw
}
