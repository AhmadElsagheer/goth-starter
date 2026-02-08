package xerrors

type Code int

const (
	codeInternal Code = iota
	codeInvalidArgument
	codeNotFound
	codeCanceled
	codeNotImplemented
	codeUnavailable
	codeUnauthorized
)

func NewInternal(message string, opts ...option) error {
	return newError(codeInternal, message, opts...)
}

func NewInvalidArgument(message string, opts ...option) error {
	return newError(codeInvalidArgument, message, opts...)
}
func NewNotFound(message string, opts ...option) error {
	return newError(codeNotFound, message, opts...)
}
func NewCanceled(message string) error {
	return newError(codeCanceled, message)
}
func NewNotImplemented(message string) error {
	return newError(codeNotImplemented, message)
}
func NewUnavailable(message string) error {
	return newError(codeUnavailable, message)
}

func NewUnauthorized(message unauthorizedMessage, opts ...option) error {
	return newError(codeUnauthorized, string(message), opts...)
}

func IsNotFound(err error) bool {
	return isCode(err, codeNotFound)
}

func IsInvalidArgument(err error) bool {
	return isCode(err, codeInvalidArgument)
}
