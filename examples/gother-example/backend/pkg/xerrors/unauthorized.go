package xerrors

type unauthorizedMessage string

const (
	MissingAuthorizationHeader     = "Missing authorization header"
	InvalidAuthorizationHeader     = "Invalid authorization header"
	InsufficientPermissionsOrScope = "Insufficient permissions or scope"
)
