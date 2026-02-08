package xerrors

var _codeToHTTPStatus = map[Code]int{
	codeInvalidArgument: 400,
	codeUnauthorized:    403,
	codeNotFound:        404,
	codeCanceled:        499,
	codeInternal:        500,
	codeNotImplemented:  501,
	codeUnavailable:     503,
}
