package database

import (
	"context"

	"github.com/bokwoon95/sq"
	"github.com/georgysavva/scany/v2/dbscan"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type executorConnection interface {
	Release()

	Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
}

type executor interface {
	Query(ctx context.Context, dest any, query sq.Query) error
	QueryOne(ctx context.Context, dest any, query sq.Query) error
	Exec(ctx context.Context, query sq.Query) (pgconn.CommandTag, error)
}

type executorImpl struct {
	acquireConnection func(context.Context) (executorConnection, error)
	scanner           *pgxscan.API
}

func newExecutor(connectionFunc func(context.Context) (executorConnection, error)) executor {
	return &executorImpl{
		acquireConnection: connectionFunc,
		scanner:           newScanner(),
	}
}

func (e *executorImpl) Query(ctx context.Context, dest any, query sq.Query) error {
	q, args, err := compileQuery(query)
	if err != nil {
		return err
	}

	conn, err := e.acquireConnection(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, q, args...)
	if err != nil {
		return err
	}

	return e.scanner.ScanAll(dest, rows)
}

func (e *executorImpl) QueryOne(ctx context.Context, dest any, query sq.Query) error {
	q, args, err := compileQuery(query)
	if err != nil {
		return err
	}

	conn, err := e.acquireConnection(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, q, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	if !rows.Next() {
		return ErrNoRows
	}

	err = e.scanner.ScanRow(dest, rows)
	if err != nil {
		return err
	}

	if rows.Next() {
		return ErrMultipleRows
	}

	return rows.Err()
}

func (e *executorImpl) Exec(ctx context.Context, query sq.Query) (pgconn.CommandTag, error) {
	q, args, err := compileQuery(query)
	if err != nil {
		return pgconn.CommandTag{}, err
	}

	conn, err := e.acquireConnection(ctx)
	if err != nil {
		return pgconn.CommandTag{}, err
	}
	defer conn.Release()

	return conn.Exec(ctx, q, args...)
}

func compileQuery(query sq.Query) (string, []any, error) {
	return sq.ToSQL(dialect, query, nil)
}

func newScanner() *pgxscan.API {
	// allowUnknownColumns to true is necessary to ensure that old versions can handle
	// unknown columns from the database without errors. This is for example useful when
	// deploying a new version of the application that has new columns in the database,
	// but old versions don't know about them. Also useful for rolling back, etc.
	//
	// TODO(sagheer): Ensure that don't use "select *" so the below option is unnecessary.
	pgxscanOptions := dbscan.WithAllowUnknownColumns(true)
	pgxscanDbAPI, err := pgxscan.NewDBScanAPI(pgxscanOptions)
	if err != nil {
		panic(err)
	}
	scanner, err := pgxscan.NewAPI(pgxscanDbAPI)
	if err != nil {
		panic(err)
	}
	return scanner
}
