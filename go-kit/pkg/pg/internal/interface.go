package internal

import (
	"context"

	"github.com/bokwoon95/sq"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

//go:generate mockgen -destination ./mockpg/pgx_tx.go -package mockpg github.com/jackc/pgx/v5 Tx
//go:generate mockgen -destination ./mockpg/pgx_rows.go -package mockpg github.com/jackc/pgx/v5 Rows
//go:generate mockgen -source ./interface.go -destination ./mockpg/pg.go -package mockpg

type Client interface {
	AcquireConnection(context.Context) (Connection, error)
	ConnectionConfig() pgx.ConnConfig
	Logger() *zap.Logger
}

type Connection interface {
	Release()

	Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	BeginTx(context.Context, pgx.TxOptions) (pgx.Tx, error)
}

type Database interface {
	// Query executes the provided read query and populates dest with the result.
	// The query may return zero or multiple rows. `dest` must be a pointer to a slice.
	Query(ctx context.Context, dest any, query sq.Query) error

	// QueryOne executes the provided read query and populates dest with the result.
	// QueryOne fails if the query returns zero or multiple rows.
	QueryOne(ctx context.Context, dest any, query sq.Query) error

	// Exec executes the provided mutation query
	Exec(ctx context.Context, query sq.Query) (pgconn.CommandTag, error)

	Transaction(ctx context.Context, txFunc func(tctx context.Context) error, opts ...TxOption) error
}

type Transaction interface {
	// Query executes the provided read query and populates dest with the result.
	// The query may return zero or multiple rows. `dest` must be a pointer to a slice.
	Query(ctx context.Context, dest any, query sq.Query) error

	// QueryOne executes the provided read query and populates dest with the result.
	// QueryOne fails if the query returns zero or multiple rows.
	QueryOne(ctx context.Context, dest any, query sq.Query) error

	// Exec executes the provided mutation query
	Exec(ctx context.Context, query sq.Query) (pgconn.CommandTag, error)
}

// TODO(sagheer): use composition for both interfaces
