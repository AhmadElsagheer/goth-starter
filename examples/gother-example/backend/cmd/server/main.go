package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/ahmad/gother-example/pkg/authprovider"
	"github.com/ahmad/gother-example/pkg/pg"
	"github.com/ahmad/gother-example/pkg/xerrors"
	"github.com/ahmad/gother-example/server/modules/auth"
	authrepo "github.com/ahmad/gother-example/server/modules/auth/repo"
	authservice "github.com/ahmad/gother-example/server/modules/auth/service"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

const (
	serverPort = 8081
	jwksURL    = "http://localhost:8787/api/auth/jwks"
	pgHost     = "localhost"
	pgPort     = 5432
	pgUser     = "gother-example"
	pgPass     = "gother-example"
	pgDB       = "gother-example"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	logger, _ := zap.NewDevelopment()

	// 1. Initialize Postgres
	db, err := initDB(logger)
	if err != nil {
		return err
	}

	// 2. Initialize Auth Service
	svc, err := initAuthService(db, logger)
	if err != nil {
		return err
	}

	// 3. Start HTTP Server with auth middleware
	e := echo.New()
	e.HTTPErrorHandler = errorHandler(logger)
	e.Use(svc.Middleware())
	e.GET("/ping", func(c echo.Context) error {
		user, err := auth.GetUserFromEcho(c)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusOK, user)
	})
	e.GET("/admin/ping", func(c echo.Context) error {
		user, err := auth.GetUserFromEcho(c)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusOK, user)
	}, auth.WithPermission("users:read"))

	fmt.Printf("Starting server on :%d\n", serverPort)
	return e.Start(fmt.Sprintf(":%d", serverPort))
}

func initDB(logger *zap.Logger) (pg.Database, error) {
	pgConfig := pg.PostgresConfig{
		Database: pgDB,
		Config: pg.DatabaseConfig{
			Host: pgHost,
			Port: pgPort,
			Auth: pg.DatabaseAuth{
				Username: pgUser,
				Password: pgPass,
			},
			Pool: pg.DatabasePool{
				MinConns: 2,
				MaxConns: 10,
			},
		},
	}

	pgClient, err := pg.NewClient(pgConfig, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create pg client: %w", err)
	}
	return pg.NewDatabase(pgClient), nil
}

func initAuthService(db pg.Database, logger *zap.Logger) (auth.Service, error) {
	repo := authrepo.New(db)
	authProvider, err := authprovider.NewJWKSProvider(jwksURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create auth provider: %w", err)
	}
	svc := authservice.New(authservice.NewServiceArgs{
		Repo:         repo,
		AuthProvider: authProvider,
		Logger:       logger,
	})
	return svc, nil
}

func errorHandler(logger *zap.Logger) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		var (
			status  int
			message string
		)
		switch {
		case errors.Is(err, echo.ErrNotFound):
			status, message = http.StatusNotFound, fmt.Sprintf("Path not found: %s", c.Request().URL)
		case errors.Is(err, echo.ErrMethodNotAllowed):
			status, message = http.StatusMethodNotAllowed, fmt.Sprintf("Method %s not allowed for path : %s", c.Request().Method, c.Request().URL)
		default:
			coder, ok := err.(xerrors.HTTPCoder)
			switch {
			case !ok:
				status, message = 500, "Internal Server Error"
			case coder.HTTPStatus() >= 500:
				status, message = coder.HTTPStatus(), "Internal Server Error"
			default:
				status, message = coder.HTTPStatus(), err.Error()
			}
		}

		logger.Error("http error", zap.Error(err))
		err = c.JSON(status, map[string]string{"message": message})
		if err != nil {
			logger.Error("Unable to send error response", zap.Error(err))
		}
	}
}
