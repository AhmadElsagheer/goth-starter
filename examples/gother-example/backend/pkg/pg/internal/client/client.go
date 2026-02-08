package client

import (
	"context"
	"fmt"
	"net/url"

	"github.com/ahmad/gother-example/pkg/pg/internal"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"go.uber.org/zap"
)

type client struct {
	connConfig *pgx.ConnConfig
	pool       *pgxpool.Pool
	logger     *zap.Logger
}

func New(conf internal.PostgresConfig, log *zap.Logger) (internal.Client, error) {
	log = log.With(zap.String("component", "pg_client"))

	poolConfig, err := pgxpool.ParseConfig(connectionString(conf))
	if err != nil {
		return nil, err
	}
	poolConfig.ConnConfig.ConnectTimeout = conf.Config.Timeout
	poolConfig.MaxConns = conf.Config.Pool.MaxConns
	poolConfig.MinConns = conf.Config.Pool.MinConns
	poolConfig.MaxConnIdleTime = conf.Config.Pool.MaxConnIdleTime
	poolConfig.MaxConnLifetime = conf.Config.Pool.MaxConnLifetime
	poolConfig.ConnConfig.Tracer = &tracelog.TraceLog{
		Logger:   &logger{logger: log.Sugar()},
		LogLevel: tracelog.LogLevelInfo,
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, err
	}

	return &client{pool: pool, connConfig: poolConfig.ConnConfig, logger: log}, nil
}

func (c *client) AcquireConnection(ctx context.Context) (internal.Connection, error) {
	return c.pool.Acquire(ctx)
}

func (c *client) ConnectionConfig() pgx.ConnConfig {
	return *c.connConfig
}

func (c *client) Logger() *zap.Logger {
	return c.logger
}

func connectionString(conf internal.PostgresConfig) string {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		conf.Config.Host,
		conf.Config.Port,
		url.QueryEscape(conf.Config.Auth.Username),
		url.QueryEscape(conf.Config.Auth.Password),
		url.QueryEscape(conf.Database))
	if !conf.Config.SSL.Enabled {
		connStr += " sslmode=" + "disable"
	}
	return connStr
}
