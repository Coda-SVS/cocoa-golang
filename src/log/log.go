package log

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"

	"github.com/Kor-SVS/cocoa/src/core"
	"github.com/Kor-SVS/cocoa/src/util"
	"gopkg.in/natefinch/lumberjack.v2"
)

func init() {
	basePath := util.GetExecutablePath()

	infoFileWriter := &lumberjack.Logger{
		Filename:   path.Join(basePath, "log", "syntool.log"),
		MaxSize:    100, // megabytes
		MaxBackups: 3,
		MaxAge:     28,   //days
		Compress:   true, // disabled by default
	}

	errorFileWriter := &lumberjack.Logger{
		Filename:   path.Join(basePath, "log", "error_syntool.log"),
		MaxSize:    100, // megabytes
		MaxBackups: 3,
		MaxAge:     28,   //days
		Compress:   true, // disabled by default
	}

	infoWriter := io.MultiWriter(infoFileWriter, os.Stdout)
	errorWriter := io.MultiWriter(infoFileWriter, errorFileWriter, os.Stderr)

	loggerTrace := log.New(infoWriter, "[TRACE] ", log.Ldate|log.Ltime)
	loggerInfo := log.New(infoWriter, "[INFO ] ", log.Ldate|log.Ltime)
	loggerWarning := log.New(errorWriter, "[WARN ] ", log.Ldate|log.Ltime)
	loggerError := log.New(errorWriter, "[ERROR] ", log.Ldate|log.Ltime)

	option := NewLoggerOption()
	option.Prefix = "[root]"

	logWriter := NewLogWriter(loggerTrace, loggerInfo, loggerWarning, loggerError)

	RootLogger = newLogger(nil, option, logWriter)
}

var (
	RootLogger *Logger
)

func NewLogger(option LoggerOption, logWriter LogWriter) *Logger {
	return newLogger(RootLogger, option, logWriter)
}

func PanicLogHandler(logger *Logger, f func() *core.ErrorW) func() {
	return func() {
		err := f()

		if err != nil {
			box := NewLogBox()
			box.message = fmt.Sprintf("<PanicLogHandler> 처리되지 않은 오류 발생 (Critical=%v, err=%v)", err.Critical, err.Error())
			box.AddCallStack(1)

			logger.Errorb(box)

			if err.Critical {
				panic(err)
			}
		}
	}
}
