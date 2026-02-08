package pg

import "{{BACKEND_MODULE}}/pkg/pg/internal/database"

var (
	ErrNoRows       = database.ErrNoRows
	ErrMultipleRows = database.ErrMultipleRows
)
