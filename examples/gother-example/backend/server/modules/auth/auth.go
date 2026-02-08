package auth

import (
	"context"

	"github.com/ahmad/gother-example/pkg/phone"
	"github.com/ahmad/gother-example/pkg/uuid"
	"github.com/ahmad/gother-example/pkg/xerrors"

	"github.com/labstack/echo/v4"
)

type Service interface {
	Middleware() echo.MiddlewareFunc
	GetUserByID(context.Context, uuid.UUID) (User, error)
}

type User struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Roles []Role    `json:"roles"`

	Email       string             `json:"email"`
	PhoneNumber *phone.PhoneNumber `json:"phoneNumber"`
}

type userAuthDataCtx string

const UserAuthDataCtxKey userAuthDataCtx = "user"

type UserAuthData struct {
	ID    uuid.UUID
	Name  string
	Roles []Role

	// Teleporting auth error from auth middleware to permission middleware
	// If the route has no permission, then auth error is ignored
	// If the route has permission and there is auth error, then it is returned as http erro
	AuthError error
}

func GetUserFromEcho(c echo.Context) (UserAuthData, error) {
	userAuthData, ok := c.Get(string(UserAuthDataCtxKey)).(UserAuthData)
	if !ok {
		return UserAuthData{}, xerrors.NewInternal("User auth data not in echo context. This might be because 'GetUserFromEcho' is called on an auth-free endpoint")
	}
	return userAuthData, nil
}

func WithUser(ctx context.Context, user UserAuthData) context.Context {
	return context.WithValue(ctx, UserAuthDataCtxKey, user)
}

func GetUser(ctx context.Context) (UserAuthData, error) {
	user, ok := ctx.Value(UserAuthDataCtxKey).(UserAuthData)
	if !ok {
		return UserAuthData{}, xerrors.NewInternal("User auth data not in echo context. This might be because 'GetUser' is called on an auth-free endpoint")
	}
	return user, nil
}

func set(in []string) map[string]struct{} {
	out := make(map[string]struct{}, len(in))
	for _, e := range in {
		out[e] = struct{}{}
	}
	return out
}
