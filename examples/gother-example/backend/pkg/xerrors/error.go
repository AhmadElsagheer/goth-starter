package xerrors

import (
	"fmt"
	"io"
	"runtime/debug"
)

type xerror struct {
	code    Code
	message string
	cause   error
	stack   []byte
}

func newError(code Code, message string, opts ...option) error {
	err := &xerror{code: code, message: message, stack: debug.Stack()}
	for _, opt := range opts {
		opt(err)
	}
	return err
}

func (err *xerror) Error() string { return err.message }

func (err *xerror) Unwrap() error { return err.cause }

func (err *xerror) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		fmt.Fprintf(s, "%+v", err.Error())
		// Add cause if exists
		if err.cause != nil {
			fmt.Fprintf(s, ": %v", err.cause)
		}
		// Add error stack
		if s.Flag('+') {
			fmt.Fprintf(s, "\n%s", string(err.stack))
		}
	case 's':
		_, _ = io.WriteString(s, err.Error())
	case 'q':
		fmt.Fprintf(s, "%q", err.Error())
	}
}

func (err *xerror) HTTPStatus() int {
	return _codeToHTTPStatus[err.code]
}

func isCode(err error, code Code) bool {
	xerr, ok := err.(*xerror)
	if !ok {
		return false
	}
	return xerr.code == code
}

type HTTPCoder interface {
	HTTPStatus() int
}
