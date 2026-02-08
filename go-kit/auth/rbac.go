package auth

import (
	"net/http"

	"{{BACKEND_MODULE}}/pkg/xerrors"

	"github.com/labstack/echo/v4"
)

type (
	Role          string
	PermissionSet map[string]struct{}
)

const (
	RoleCustomer Role = "customer"
	RoleAdmin    Role = "admin"
)

var permissionsByRole = map[Role]PermissionSet{
	RoleCustomer: set([]string{}),
	RoleAdmin: set([]string{
		"users:read",
	}),
}

func (r Role) HasPermission(p string) bool {
	_, ok := permissionsByRole[r][p]
	return ok
}

func WithPermission(p string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userAuthData, ok := c.Get(string(UserAuthDataCtxKey)).(UserAuthData)
			if !ok {
				return echo.NewHTTPError(http.StatusInternalServerError, "User auth data not in echo context")
			}

			if userAuthData.AuthError != nil {
				return userAuthData.AuthError
			}

			for _, role := range userAuthData.Roles {
				if role.HasPermission(p) {
					return next(c)
				}
			}
			return xerrors.NewUnauthorized(xerrors.InsufficientPermissionsOrScope)
		}
	}
}
