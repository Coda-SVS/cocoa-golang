package log

// TODO: 비동기 로깅 구현

type Logger struct {
	option    LoggerOption
	parent    *Logger
	logWriter *LogWriter
}

func newLogger(parentLogger *Logger, option LoggerOption, logWriter *LogWriter) *Logger {
	return &Logger{
		option:    option,
		parent:    parentLogger,
		logWriter: logWriter,
	}
}

func (l *Logger) LogLevel() LogLevel {
	return l.option.LogLevel
}

func (l *Logger) Parent() *Logger {
	return l.parent
}

func (l *Logger) NewLogger(option LoggerOption, logWriter *LogWriter) *Logger {
	return newLogger(l, option, logWriter)
}

func (l *Logger) NewSimpleLogger(option LoggerOption) *Logger {
	logWriter := NewLogWriter(nil, nil, nil, nil, nil, nil, nil, nil)
	return newLogger(l, option, logWriter)
}

func (l *Logger) trace(box *LogBox) {
	if l.option.LogLevel <= TRACE {
		l.messagePrefixBuild(box)

		logMutex.Lock()
		if l.logWriter.consoleLogWriter.loggerTrace != nil {
			l.logWriter.consoleLogWriter.loggerTrace.Println(box.BuildMessage(true))
		}
		if l.logWriter.fileLogWriter.loggerTrace != nil {
			l.logWriter.fileLogWriter.loggerTrace.Println(box.BuildMessage(false))
		}
		logMutex.Unlock()

		if l.parent != nil && l.option.PassToParent {
			l.parent.trace(box)
		}
	}
}

func (l *Logger) info(box *LogBox) {
	if l.option.LogLevel <= INFO {
		l.messagePrefixBuild(box)

		logMutex.Lock()
		if l.logWriter.consoleLogWriter.loggerInfo != nil {
			l.logWriter.consoleLogWriter.loggerInfo.Println(box.BuildMessage(true))
		}
		if l.logWriter.fileLogWriter.loggerInfo != nil {
			l.logWriter.fileLogWriter.loggerInfo.Println(box.BuildMessage(false))
		}
		logMutex.Unlock()

		if l.parent != nil && l.option.PassToParent {
			l.parent.info(box)
		}
	}
}

func (l *Logger) warning(box *LogBox) {
	if l.option.LogLevel <= WARNING {
		l.messagePrefixBuild(box)

		logMutex.Lock()
		if l.logWriter.consoleLogWriter.loggerWarning != nil {
			l.logWriter.consoleLogWriter.loggerWarning.Println(box.BuildMessage(true))
		}
		if l.logWriter.fileLogWriter.loggerWarning != nil {
			l.logWriter.fileLogWriter.loggerWarning.Println(box.BuildMessage(false))
		}
		logMutex.Unlock()

		if l.parent != nil && l.option.PassToParent {
			l.parent.warning(box)
		}
	}
}

func (l *Logger) error(box *LogBox) {
	if l.option.LogLevel <= ERROR {
		l.messagePrefixBuild(box)

		logMutex.Lock()
		if l.logWriter.consoleLogWriter.loggerError != nil {
			l.logWriter.consoleLogWriter.loggerError.Println(box.BuildMessage(true))
		}
		if l.logWriter.fileLogWriter.loggerError != nil {
			l.logWriter.fileLogWriter.loggerError.Println(box.BuildMessage(false))
		}
		logMutex.Unlock()

		if l.parent != nil && l.option.PassToParent {
			l.parent.error(box)
		}
	}
}

func (l *Logger) messagePrefixBuild(box *LogBox) {
	// box.logBoxInternal.joinedPrefix = stringConcat("", l.option.Prefix, box.logBoxInternal.joinedPrefix)

	if box.logBoxInternal.joinedPrefix == "" {
		box.logBoxInternal.joinedPrefix = l.option.Prefix
	}
}
