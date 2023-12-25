package log

import (
	"fmt"
)

func (l *Logger) Direct(box *LogBox) {
	l.logBoxHandler(box, true)
}

func (l *Logger) Trace(v ...any) {
	box := NewLogBox(TRACE)
	box.message = fmt.Sprint(v...)
	box.AddCallStack(1)
	l.logBoxHandler(box, true)
}

func (l *Logger) Info(v ...any) {
	box := NewLogBox(INFO)
	box.message = fmt.Sprint(v...)
	box.AddCallStack(1)
	l.logBoxHandler(box, true)
}

func (l *Logger) Warning(v ...any) {
	box := NewLogBox(WARNING)
	box.message = fmt.Sprint(v...)
	box.AddCallStack(1)
	l.logBoxHandler(box, true)
}

func (l *Logger) Error(v ...any) {
	box := NewLogBox(ERROR)
	box.message = fmt.Sprint(v...)
	box.AddCallStack(1)
	l.logBoxHandler(box, true)
}

func (l *Logger) Tracef(msg string, v ...any) {
	box := NewLogBox(TRACE)
	box.message = fmt.Sprintf(msg, v...)
	box.AddCallStack(1)
	l.logBoxHandler(box, true)
}

func (l *Logger) Infof(msg string, v ...any) {
	box := NewLogBox(INFO)
	box.message = fmt.Sprintf(msg, v...)
	box.AddCallStack(1)
	l.logBoxHandler(box, true)
}

func (l *Logger) Warningf(msg string, v ...any) {
	box := NewLogBox(WARNING)
	box.message = fmt.Sprintf(msg, v...)
	box.AddCallStack(1)
	l.logBoxHandler(box, true)
}

func (l *Logger) Errorf(msg string, v ...any) {
	box := NewLogBox(ERROR)
	box.message = fmt.Sprintf(msg, v...)
	box.AddCallStack(1)
	l.logBoxHandler(box, true)
}
