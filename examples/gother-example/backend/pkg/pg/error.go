package pg

import "github.com/ahmad/gother-example/pkg/pg/internal/database"

var (
	ErrNoRows       = database.ErrNoRows
	ErrMultipleRows = database.ErrMultipleRows
)
