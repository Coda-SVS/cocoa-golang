package util

import (
	"fmt"

	"github.com/Kor-SVS/cocoa/src/core"
)

func PanicToErrorW(f func()) func() *core.ErrorW {
	return func() (err *core.ErrorW) {
		defer func() {
			r := recover()
			if r != nil {
				errw, ok := r.(core.ErrorW)
				if ok {
					err = &errw
					return
				}

				_err, ok := r.(error)
				if ok {
					err = core.NewErrorW(_err, true)
					return
				}

				err = core.NewErrorW(fmt.Errorf("%v", r), true)
			}
		}()

		f()

		return err
	}
}
