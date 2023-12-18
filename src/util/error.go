package util

import (
	"fmt"
)

type ErrorW struct {
	error
	critical bool
	callFile string
}

func NewErrorW(errorObj error, critical bool, callStackSkip int) *ErrorW {
	ew := &ErrorW{
		error:    errorObj,
		critical: critical,
	}
	ew.AddCallStack(callStackSkip+1, true)
	return ew
}

func (ew *ErrorW) Critical() bool {
	return ew.critical
}

func (ew *ErrorW) CallFile() string {
	return ew.callFile
}

func (ew *ErrorW) AddCallStack(skip int, force bool) bool {
	if !force && ew.callFile != "" {
		return true
	}

	callFile, ok := GetCallFileFromCallStack(skip + 1)

	if !ok {
		return false
	} else {
		ew.callFile = callFile
		return true
	}
}

func PanicToErrorW(f func()) func() *ErrorW {
	return func() (err *ErrorW) {
		defer func() {
			r := recover()
			if r != nil {
				callStackSkip := 1
				errw, ok := r.(ErrorW)
				if ok {
					errw.AddCallStack(callStackSkip, false)
					err = &errw
					return
				}

				_err, ok := r.(error)
				if ok {
					err = NewErrorW(_err, true, callStackSkip)
					return
				}

				err = NewErrorW(fmt.Errorf("%v", r), true, callStackSkip)
			}
		}()

		f()

		return err
	}
}
