package log

import (
	"io"
	"log"
	"os"
	"path"

	"gopkg.in/natefinch/lumberjack.v2"
)

type Level int

type Logger struct {
	LogLevel LogLevel
}

var (
	loggerTrace *log.Logger
	loggerInfo  *log.Logger
	loggerWarn  *log.Logger
	loggerError *log.Logger
)

func NewLogger() *Logger {
	return &Logger{
		LogLevel: TRACE,
	}
}

func (logger Logger) Trace(v ...any) {
	if logger.LogLevel <= TRACE {
		loggerTrace.Println(v...)
	}
}

func (logger Logger) Info(v ...any) {
	if logger.LogLevel <= INFO {
		loggerInfo.Println(v...)
	}
}

func (logger Logger) Warn(v ...any) {
	if logger.LogLevel <= WARNING {
		loggerWarn.Println(v...)
	}
}

func (logger Logger) Error(v ...any) {
	if logger.LogLevel <= ERROR {
		loggerError.Println(v...)
	}
}

func (logger Logger) Tracef(msg string, v ...any) {
	if logger.LogLevel <= TRACE {
		loggerTrace.Printf(msg, v...)
	}
}

func (logger Logger) Infof(msg string, v ...any) {
	if logger.LogLevel <= INFO {
		loggerInfo.Printf(msg, v...)
	}
}

func (logger Logger) Warnf(msg string, v ...any) {
	if logger.LogLevel <= WARNING {
		loggerWarn.Printf(msg, v...)
	}
}

func (logger Logger) Errorf(msg string, v ...any) {
	if logger.LogLevel <= ERROR {
		loggerError.Printf(msg, v...)
	}
}

func init() {
	basePath := getExecutablePath()

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

	loggerTrace = log.New(infoWriter, "[TRACE] ", log.Ldate|log.Ltime|log.Lshortfile)
	loggerInfo = log.New(infoWriter, "[INFO ] ", log.Ldate|log.Ltime|log.Lshortfile)
	loggerWarn = log.New(errorWriter, "[WARN ] ", log.Ldate|log.Ltime|log.Lshortfile)
	loggerError = log.New(errorWriter, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
}

// 프로그램 작업 디렉토리 가져오기
func getExecutablePath() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return path.Dir(ex)
}
