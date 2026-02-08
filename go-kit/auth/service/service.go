package service

import (
	"context"
	"strings"

	"{{AUTH_MODULE}}/repo"
	"{{BACKEND_MODULE}}/pkg/authprovider"
	"{{BACKEND_MODULE}}/pkg/uuid"
	"{{BACKEND_MODULE}}/pkg/xerrors"
	"{{BACKEND_MODULE}}/server/modules/auth"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type service struct {
	repo         repo.Repository
	authProvider authprovider.Provider
	logger       *zap.Logger
}

type NewServiceArgs struct {
	fx.In

	Repo         repo.Repository
	AuthProvider authprovider.Provider
	Logger       *zap.Logger
}

func New(args NewServiceArgs) auth.Service {
	return &service{
		repo:         args.Repo,
		authProvider: args.AuthProvider,
		logger:       args.Logger,
	}
}

func (s *service) GetUserByID(ctx context.Context, id uuid.UUID) (auth.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			defer func() {
				if err != nil {
					// If there is error, teleport it to permission middleware
					c.Set(string(auth.UserAuthDataCtxKey), auth.UserAuthData{AuthError: err})
					err = next(c)
				}
			}()

			req := c.Request()

			token := req.Header.Get("Authorization")
			if token == "" {
				return xerrors.NewUnauthorized(xerrors.MissingAuthorizationHeader)
			}
			token = strings.TrimPrefix(token, "Bearer ")

			decodedToken, err := s.authProvider.VerifyIDToken(req.Context(), token)
			if err != nil {
				return xerrors.NewUnauthorized(xerrors.InvalidAuthorizationHeader, xerrors.Cause(err))
			}

			userID, err := uuid.Parse(decodedToken.ID)
			if err != nil {
				return xerrors.NewInternal("Invalid user id in token: "+decodedToken.ID, xerrors.Cause(err))
			}

			user, err := s.repo.GetByID(req.Context(), userID)
			if err != nil {
				return xerrors.NewInternal("Failed to get user "+userID.String(), xerrors.Cause(err))
			}

			authData := auth.UserAuthData{
				ID:    user.ID,
				Name:  user.Name,
				Roles: user.Roles,
			}

			c.Set(string(auth.UserAuthDataCtxKey), authData)
			c.SetRequest(req.WithContext(auth.WithUser(req.Context(), authData)))

			return next(c)
		}
	}
}
