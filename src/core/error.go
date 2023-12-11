package core

type ErrorW struct {
	error
	Critical bool
}

func NewErrorW(errorObj error, critical bool) *ErrorW {
	return &ErrorW{
		error:    errorObj,
		Critical: critical,
	}
}
