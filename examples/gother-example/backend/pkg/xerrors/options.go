package xerrors

import "fmt"

type option func(*xerror)

func Cause(cause error) option {
	return func(err *xerror) {
		err.cause = cause
	}
}

func Args(args ...any) option {
	return func(err *xerror) {
		err.message = fmt.Sprintf(err.message, args...)
	}
}
