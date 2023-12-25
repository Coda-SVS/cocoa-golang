package log

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"

	"github.com/Kor-SVS/cocoa/src/util"
	"github.com/fatih/color"
	"gopkg.in/natefinch/lumberjack.v2"
)

func init() {
	basePath := util.GetExecutablePath()

	infoFileWriter := &lumberjack.Logger{
		Filename:   path.Join(basePath, "log", "syntool.log"),
		MaxSize:    20, // megabytes
		MaxBackups: 3,
		MaxAge:     28,   //days
		Compress:   true, // disabled by default
	}

	errorFileWriter := &lumberjack.Logger{
		Filename:   path.Join(basePath, "log", "error_syntool.log"),
		MaxSize:    50, // megabytes
		MaxBackups: 3,
		MaxAge:     28,   //days
		Compress:   true, // disabled by default
	}

	infoFileLogWriter := infoFileWriter
	errorFileLogWriter := io.MultiWriter(infoFileWriter, errorFileWriter)
	infoConsoleLogWriter := os.Stdout
	errorConsoleLogWriter := os.Stderr

	fileLoggerTrace := log.New(infoFileLogWriter, "[TRACE] ", log.Ldate|log.Ltime)
	fileLoggerInfo := log.New(infoFileLogWriter, "[INFO ] ", log.Ldate|log.Ltime)
	fileLoggerWarning := log.New(errorFileLogWriter, "[WARN ] ", log.Ldate|log.Ltime)
	fileLoggerError := log.New(errorFileLogWriter, "[ERROR] ", log.Ldate|log.Ltime)

	ConsoleLoggerTrace := log.New(infoConsoleLogWriter, color.New(color.FgHiBlue).Sprint("[TRACE] "), log.Ldate|log.Ltime)
	ConsoleLoggerInfo := log.New(infoConsoleLogWriter, color.New(color.FgHiGreen).Sprint("[INFO ] "), log.Ldate|log.Ltime)
	ConsoleLoggerWarning := log.New(errorConsoleLogWriter, color.New(color.FgHiYellow).Sprint("[WARN ] "), log.Ldate|log.Ltime)
	ConsoleLoggerError := log.New(errorConsoleLogWriter, color.New(color.FgHiRed).Sprint("[ERROR] "), log.Ldate|log.Ltime)

	option := NewLoggerOption()
	option.Prefix = "[root]"

	logWriter := NewLogWriter(
		fileLoggerTrace,
		fileLoggerInfo,
		fileLoggerWarning,
		fileLoggerError,
		ConsoleLoggerTrace,
		ConsoleLoggerInfo,
		ConsoleLoggerWarning,
		ConsoleLoggerError,
	)

	rootLogger = newLogger(nil, option, logWriter)
}

var (
	rootLogger *Logger
)

func RootLogger() *Logger {
	return rootLogger
}

func PanicLogHandler(logger *Logger, f func() *util.ErrorW) func() {
	return func() {
		err := f()

		if err != nil {
			box := NewLogBox(ERROR)
			box.message = fmt.Sprintf(
				"<PanicLogHandler> 처리되지 않은 오류 발생 (Critical=%v, err=%v)",
				err.Critical(),
				err.Error(),
			)
			box.AddCallStackFromError(err, 1)

			logger.Direct(box)

			if err.Critical() {
				panic(err)
			}
		}
	}
}
