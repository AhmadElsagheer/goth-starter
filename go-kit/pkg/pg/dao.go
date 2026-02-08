package pg

import (
	"context"
	"errors"

	"{{BACKEND_MODULE}}/pkg/pg/internal"

	"github.com/bokwoon95/sq"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	errMissingPredicates  = errors.New("dao: operation must contain 'where' predicates")
	errMissingAssignments = errors.New("dao: operation must contain 'set' assignments")
)

type dao[T sq.Table, E any] struct {
	db    internal.Database
	table T
}

func NewDao[T sq.Table, E any](db internal.Database) Dao[T, E] {
	return &dao[T, E]{
		db:    db,
		table: sq.New[T](""),
	}
}

func (d *dao[T, E]) Table() T {
	return d.table
}

func (d *dao[T, E]) Insert(ctx context.Context, colmapper func(col *sq.Column)) (pgconn.CommandTag, error) {
	q := sq.InsertInto(d.table).ColumnValues(colmapper)

	return d.db.Exec(ctx, q)
}

func (d *dao[T, E]) Upsert(
	ctx context.Context,
	colmapper func(col *sq.Column),
	conflictField sq.Field,
	assignmentsOnConflict ...sq.Assignment,
) (pgconn.CommandTag, error) {
	q := sq.Postgres.InsertInto(d.table).ColumnValues(colmapper).
		OnConflict(conflictField).
		DoUpdateSet(assignmentsOnConflict...)

	return d.db.Exec(ctx, q)
}

func (d *dao[T, E]) Update(ctx context.Context, opts ...internal.Option) (pgconn.CommandTag, error) {
	options := internal.BuildOptions(opts...)

	if len(options.Predicates) == 0 {
		return pgconn.CommandTag{}, errMissingPredicates
	}

	if len(options.Assignments) == 0 {
		return pgconn.CommandTag{}, errMissingAssignments
	}

	q := sq.Update(d.table).Set(options.Assignments...).Where(options.Predicates...)

	return d.db.Exec(ctx, q)
}

func (d *dao[T, E]) List(ctx context.Context, opts ...internal.Option) ([]E, error) {
	options := internal.BuildOptions(opts...)

	q := sq.Select(AllFields(d.table)...).From(d.table)

	if len(options.Predicates) > 0 {
		q = q.Where(options.Predicates...)
	}

	if options.OrderBy != nil {
		q = q.OrderBy(options.OrderBy...)
	}

	if options.Limit != nil {
		q = q.Limit(*options.Limit)
	}

	var list []E
	err := d.db.Query(ctx, &list, q)
	if err != nil {
		return nil, err
	}

	return list, err
}

func (d *dao[T, E]) Get(ctx context.Context, opts ...internal.Option) (E, error) {
	var entity E

	options := internal.BuildOptions(opts...)

	if len(options.Predicates) == 0 {
		return entity, errMissingPredicates
	}

	q := sq.Postgres.Select(AllFields(d.table)...).From(d.table).Where(options.Predicates...)

	if options.OrderBy != nil {
		q = q.OrderBy(options.OrderBy...)
	}

	if options.Limit != nil {
		q = q.Limit(*options.Limit)
	}

	if options.Lock != nil {
		q = q.LockRows(*options.Lock)
	}

	err := d.db.QueryOne(ctx, &entity, q)

	return entity, err
}
