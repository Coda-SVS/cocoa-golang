package log

import (
	"fmt"
)

func (l *Logger) Trace(v ...any) {
	box := NewLogBox()
	box.message = fmt.Sprint(v...)
	box.AddCallStack(2)
	l.trace(box)
}

func (l *Logger) Info(v ...any) {
	box := NewLogBox()
	box.message = fmt.Sprint(v...)
	box.AddCallStack(2)
	l.info(box)
}

func (l *Logger) Warning(v ...any) {
	box := NewLogBox()
	box.message = fmt.Sprint(v...)
	box.AddCallStack(2)
	l.warning(box)
}

func (l *Logger) Error(v ...any) {
	box := NewLogBox()
	box.message = fmt.Sprint(v...)
	box.AddCallStack(2)
	l.error(box)
}

func (l *Logger) Tracef(msg string, v ...any) {
	box := NewLogBox()
	box.message = fmt.Sprintf(msg, v...)
	box.AddCallStack(2)
	l.trace(box)
}

func (l *Logger) Infof(msg string, v ...any) {
	box := NewLogBox()
	box.message = fmt.Sprintf(msg, v...)
	box.AddCallStack(2)
	l.info(box)
}

func (l *Logger) Warningf(msg string, v ...any) {
	box := NewLogBox()
	box.message = fmt.Sprintf(msg, v...)
	box.AddCallStack(2)
	l.warning(box)
}

func (l *Logger) Errorf(msg string, v ...any) {
	box := NewLogBox()
	box.message = fmt.Sprintf(msg, v...)
	box.AddCallStack(2)
	l.error(box)
}
