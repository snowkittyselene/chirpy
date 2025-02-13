// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: auth.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const updateUserCredentials = `-- name: UpdateUserCredentials :one
UPDATE users
SET email=$2,
    hashed_password=$3,
    updated_at=NOW()
WHERE id=$1
RETURNING id, created_at, updated_at, email, is_chirpy_red
`

type UpdateUserCredentialsParams struct {
	ID             uuid.UUID
	Email          string
	HashedPassword string
}

type UpdateUserCredentialsRow struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Email       string
	IsChirpyRed bool
}

func (q *Queries) UpdateUserCredentials(ctx context.Context, arg UpdateUserCredentialsParams) (UpdateUserCredentialsRow, error) {
	row := q.db.QueryRowContext(ctx, updateUserCredentials, arg.ID, arg.Email, arg.HashedPassword)
	var i UpdateUserCredentialsRow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.IsChirpyRed,
	)
	return i, err
}
