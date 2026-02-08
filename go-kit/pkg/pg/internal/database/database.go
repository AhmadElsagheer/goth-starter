package database

import (
	"context"
	"errors"

	"{{BACKEND_MODULE}}/pkg/pg/internal"

	"github.com/bokwoon95/sq"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

const dialect = sq.DialectPostgres

var (
	ErrNoRows       = errors.New("no rows found")
	ErrMultipleRows = errors.New("multiple rows found")
)

type database struct {
	client          internal.Client
	defaultExecutor executor
}

func NewDatabase(client internal.Client) internal.Database {
	return &database{
		client: client,
		defaultExecutor: newExecutor(func(ctx context.Context) (executorConnection, error) {
			return client.AcquireConnection(ctx)
		}),
	}
}

func (db *database) Query(ctx context.Context, dest any, query sq.Query) error {
	return db.executor(ctx).Query(ctx, dest, query)
}

func (db *database) QueryOne(ctx context.Context, dest any, query sq.Query) error {
	return db.executor(ctx).QueryOne(ctx, dest, query)
}

func (db *database) Exec(ctx context.Context, query sq.Query) (pgconn.CommandTag, error) {
	return db.executor(ctx).Exec(ctx, query)
}

func (db *database) Transaction(ctx context.Context, txFunc func(context.Context) error, opts ...internal.TxOption) error {
	conn, err := db.client.AcquireConnection(ctx)
	if err != nil {
		return err
	}

	defer conn.Release()

	txOpts := pgx.TxOptions{IsoLevel: pgx.Serializable}
	for _, opt := range opts {
		opt(&txOpts)
	}

	tx, err := conn.BeginTx(ctx, txOpts)
	if err != nil {
		return err
	}

	defer func() {
		// Ensure to rollback if there is a panic
		err := recover()
		if err == nil {
			return
		}
		rerr := tx.Rollback(ctx)
		if rerr != nil {
			db.client.Logger().Error("transaction rollback failed", zap.Error(rerr))
		}
		panic(err)
	}()

	e := newExecutor(func(ctx context.Context) (executorConnection, error) {
		return &txWrapper{Tx: tx}, nil
	})

	tctx := contextWithTx(ctx, e)

	err = txFunc(tctx)

	if err != nil {
		rerr := tx.Rollback(ctx)
		if rerr != nil {
			db.client.Logger().Error("transaction rollback failed", zap.Error(rerr))
		}
		return err
	}

	return tx.Commit(ctx)
}

func (db *database) executor(ctx context.Context) executor {
	tx := txFromContext(ctx)
	if tx != nil {
		return tx
	}
	return db.defaultExecutor
}

type txWrapper struct {
	pgx.Tx
}

func (tx *txWrapper) Release() {
	// no-op
}
