package pg

import (
	"context"

	"github.com/ahmad/gother-example/pkg/pg/internal"
	"github.com/ahmad/gother-example/pkg/pg/internal/client"
	"github.com/ahmad/gother-example/pkg/pg/internal/database"

	"github.com/bokwoon95/sq"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

// 1. Interfaces

// Client is the lowest level component to interact with the database.
// It manages a connection pool for a single database and can lease connections.
type Client = internal.Client

// Connection is a single db connection from the pool managed by the [Client].
// It can do Exec/Query operations on any table in the database.
// Whoever acquires the connection must eventually release it.
type Connection = internal.Connection

// Database abstracts connections to simplify the logic of acquisition & release.
// It provides simple operations that can be used directly by the caller.
type Database = internal.Database

// Transaction: TODO
type Transaction = internal.Transaction

// Dao is an abstracttion over [Database] providing additional functions to read or mutate
// a single table (data access object).
type Dao[T sq.Table, E any] interface {
	// Table returns the typed table of this dao. The table is named after
	// the underlying table name (i.e. no alias). If you need to alias, create
	// a table from scratch.
	Table() T

	// Insert inserts a single record in the table. The `colmapper` function must use
	// `col` to define the assignments of each table column.
	//
	// Example:
	//  var myEntity Entity // holds the values to be used in insertion
	//
	// 	table := myDao.Table()
	// 	myDao.Insert(ctx, func(col *sq.Column) {
	// 		col.SetUUID(table.ID, myEntity.ID)
	// 		col.SetString(table.NAME, myEntity.NAME)
	// 		col.SetJSON(table.DETAILS, myEntity.Details)
	// 	})
	Insert(ctx context.Context, colmapper func(col *sq.Column)) (pgconn.CommandTag, error)

	// Update updates one or more records in the table.
	// Must pass `Where` and `Set`` options.
	//
	// Example:
	// 	table := myDao.Table()
	//	myDao.Update(ctx,
	// 		pg.Where(table.ID.Eq(myEntity.ID)),
	// 		pg.Set(table.DETAILS.Set(myEntity.Details)),
	// 	)
	Update(ctx context.Context, opts ...internal.Option) (pgconn.CommandTag, error)

	// Upsert upserts a single record in the table. Works similar to `Insert` but uses
	// `conflictField` to detect an insert conflict and `assignmentsOnConflict` to update
	// fields similar to `Set()` option in `Update` operation.
	Upsert(ctx context.Context, colmapper func(col *sq.Column), conflictField sq.Field, assignmentsOnConflict ...sq.Assignment) (pgconn.CommandTag, error)

	// List returns all records in the table.
	// `Where` option can be used to filter records.
	List(ctx context.Context, opts ...internal.Option) ([]E, error)

	// Get returns a single record in the table.
	// Must pass `Where` option. If zero or multiple rows match
	// the query, an error is returned.
	Get(ctx context.Context, opts ...internal.Option) (E, error)
}

// 2. Options
type Option = internal.Option

func Where(predicates ...sq.Predicate) internal.Option {
	return internal.Where(predicates...)
}

func Set(assignments ...sq.Assignment) internal.Option {
	return internal.Set(assignments...)
}

func LockForUpdate() internal.Option {
	return internal.LockForUpdate()
}

func OrderBy(fields ...sq.Field) internal.Option {
	return internal.OrderBy(fields...)
}

func Limit(limit int) internal.Option {
	return internal.Limit(limit)
}

type TxOption = internal.TxOption

func TxReadCommitted() TxOption {
	return internal.TxReadCommitted()
}

// 3. Constructors

func NewClient(conf PostgresConfig, log *zap.Logger) (Client, error) {
	return client.New(conf, log)
}

func NewDatabase(client Client) Database {
	return database.NewDatabase(client)
}
