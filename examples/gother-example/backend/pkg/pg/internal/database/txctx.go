package database

import (
	"context"

	"github.com/ahmad/gother-example/pkg/pg/internal"
)

type txCtxKey int

const ctxTransactionKey txCtxKey = iota

func txFromContext(ctx context.Context) internal.Transaction {
	tx, ok := ctx.Value(ctxTransactionKey).(internal.Transaction)
	if !ok {
		return nil
	}
	return tx
}

func contextWithTx(ctx context.Context, tx internal.Transaction) context.Context {
	return context.WithValue(ctx, ctxTransactionKey, tx)
}
