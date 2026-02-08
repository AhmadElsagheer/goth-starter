package internal

import (
	"github.com/bokwoon95/sq"
	"github.com/jackc/pgx/v5"
)

type Options struct {
	Predicates  []sq.Predicate
	Assignments []sq.Assignment
	OrderBy     []sq.Field
	Limit       *int
	Lock        *string
}

type Option func(*Options)

func BuildOptions(opts ...Option) Options {
	var out Options
	for _, opt := range opts {
		opt(&out)
	}
	return out
}

func Where(predicates ...sq.Predicate) Option {
	return func(qo *Options) {
		qo.Predicates = append(qo.Predicates, predicates...)
	}
}

func Set(assignments ...sq.Assignment) Option {
	return func(qo *Options) {
		qo.Assignments = append(qo.Assignments, assignments...)
	}
}

func OrderBy(fields ...sq.Field) Option {
	return func(qo *Options) {
		qo.OrderBy = append(qo.OrderBy, fields...)
	}
}

func Limit(limit int) Option {
	return func(qo *Options) {
		qo.Limit = &limit
	}
}

func LockForUpdate() Option {
	return func(qo *Options) {
		s := "FOR UPDATE"
		qo.Lock = &s
	}
}

type TxOption func(*pgx.TxOptions)

func TxReadCommitted() TxOption {
	return func(txOpts *pgx.TxOptions) {
		txOpts.IsoLevel = pgx.ReadCommitted
	}
}
