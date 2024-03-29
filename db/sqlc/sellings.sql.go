// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
// source: sellings.sql

package db

import (
	"context"
	"database/sql"
)

const createSelling = `-- name: CreateSelling :one
INSERT INTO sellings (
    partner_id, category_id, amount, created_by
) values (
             $1, $2, $3, $4
         ) RETURNING id, partner_id, category_id, amount, created_at, updated_at, deleted_at, created_by, updated_by, deleted_by
`

type CreateSellingParams struct {
	PartnerID  sql.NullInt64  `json:"partner_id"`
	CategoryID sql.NullInt64  `json:"category_id"`
	Amount     sql.NullString `json:"amount"`
	CreatedBy  sql.NullInt64  `json:"created_by"`
}

func (q *Queries) CreateSelling(ctx context.Context, arg CreateSellingParams) (Selling, error) {
	row := q.db.QueryRowContext(ctx, createSelling,
		arg.PartnerID,
		arg.CategoryID,
		arg.Amount,
		arg.CreatedBy,
	)
	var i Selling
	err := row.Scan(
		&i.ID,
		&i.PartnerID,
		&i.CategoryID,
		&i.Amount,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.CreatedBy,
		&i.UpdatedBy,
		&i.DeletedBy,
	)
	return i, err
}

const deleteSelling = `-- name: DeleteSelling :exec
DELETE FROM sellings
WHERE id = $1
`

func (q *Queries) DeleteSelling(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteSelling, id)
	return err
}

const getSelling = `-- name: GetSelling :one
SELECT id, partner_id, category_id, amount, created_at, updated_at, deleted_at, created_by, updated_by, deleted_by FROM sellings
WHERE id = $1 AND deleted_at is null LIMIT 1
`

func (q *Queries) GetSelling(ctx context.Context, id int64) (Selling, error) {
	row := q.db.QueryRowContext(ctx, getSelling, id)
	var i Selling
	err := row.Scan(
		&i.ID,
		&i.PartnerID,
		&i.CategoryID,
		&i.Amount,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.CreatedBy,
		&i.UpdatedBy,
		&i.DeletedBy,
	)
	return i, err
}

const listSelling = `-- name: ListSelling :many
SELECT id, partner_id, category_id, amount, created_at, updated_at, deleted_at, created_by, updated_by, deleted_by FROM sellings
WHERE deleted_at is null
ORDER BY partner_id, category_id
LIMIT $1
OFFSET $2
`

type ListSellingParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListSelling(ctx context.Context, arg ListSellingParams) ([]Selling, error) {
	rows, err := q.db.QueryContext(ctx, listSelling, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Selling{}
	for rows.Next() {
		var i Selling
		if err := rows.Scan(
			&i.ID,
			&i.PartnerID,
			&i.CategoryID,
			&i.Amount,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
			&i.CreatedBy,
			&i.UpdatedBy,
			&i.DeletedBy,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listSellingByParams = `-- name: ListSellingByParams :many
SELECT id, partner_id, category_id, amount, created_at, updated_at, deleted_at, created_by, updated_by, deleted_by FROM sellings
WHERE deleted_at is null
AND (CASE WHEN $3::bool THEN partner_id = $4 ELSE TRUE END)
AND (CASE WHEN $5::bool THEN category_id = $6 ELSE TRUE END)
ORDER BY partner_id, category_id
LIMIT $1
OFFSET $2
`

type ListSellingByParamsParams struct {
	Limit      int32         `json:"limit"`
	Offset     int32         `json:"offset"`
	IsPartner  bool          `json:"is_partner"`
	PartnerID  sql.NullInt64 `json:"partner_id"`
	IsCategory bool          `json:"is_category"`
	CategoryID sql.NullInt64 `json:"category_id"`
}

func (q *Queries) ListSellingByParams(ctx context.Context, arg ListSellingByParamsParams) ([]Selling, error) {
	rows, err := q.db.QueryContext(ctx, listSellingByParams,
		arg.Limit,
		arg.Offset,
		arg.IsPartner,
		arg.PartnerID,
		arg.IsCategory,
		arg.CategoryID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Selling{}
	for rows.Next() {
		var i Selling
		if err := rows.Scan(
			&i.ID,
			&i.PartnerID,
			&i.CategoryID,
			&i.Amount,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
			&i.CreatedBy,
			&i.UpdatedBy,
			&i.DeletedBy,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateInactiveSelling = `-- name: UpdateInactiveSelling :one
UPDATE sellings SET deleted_by = $2, deleted_at = now() WHERE id = $1
RETURNING id, partner_id, category_id, amount, created_at, updated_at, deleted_at, created_by, updated_by, deleted_by
`

type UpdateInactiveSellingParams struct {
	ID        int64         `json:"id"`
	DeletedBy sql.NullInt64 `json:"deleted_by"`
}

func (q *Queries) UpdateInactiveSelling(ctx context.Context, arg UpdateInactiveSellingParams) (Selling, error) {
	row := q.db.QueryRowContext(ctx, updateInactiveSelling, arg.ID, arg.DeletedBy)
	var i Selling
	err := row.Scan(
		&i.ID,
		&i.PartnerID,
		&i.CategoryID,
		&i.Amount,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.CreatedBy,
		&i.UpdatedBy,
		&i.DeletedBy,
	)
	return i, err
}

const updateSelling = `-- name: UpdateSelling :one
UPDATE sellings
SET
    partner_id = CASE
               WHEN $1::bool
                THEN $2
               ELSE partner_id
        END,
    category_id = CASE
                 WHEN $3::bool
                THEN $4
                 ELSE category_id
        END,
    amount = CASE
                WHEN $5::bool
                THEN $6
                ELSE amount
        END,
    updated_by = $7,
    updated_at = now()
WHERE
id = $8
RETURNING id, partner_id, category_id, amount, created_at, updated_at, deleted_at, created_by, updated_by, deleted_by
`

type UpdateSellingParams struct {
	SetPartnerID  bool           `json:"set_partner_id"`
	PartnerID     sql.NullInt64  `json:"partner_id"`
	SetCategoryID bool           `json:"set_category_id"`
	CategoryID    sql.NullInt64  `json:"category_id"`
	SetAmount     bool           `json:"set_amount"`
	Amount        sql.NullString `json:"amount"`
	UpdatedBy     sql.NullInt64  `json:"updated_by"`
	ID            int64          `json:"id"`
}

func (q *Queries) UpdateSelling(ctx context.Context, arg UpdateSellingParams) (Selling, error) {
	row := q.db.QueryRowContext(ctx, updateSelling,
		arg.SetPartnerID,
		arg.PartnerID,
		arg.SetCategoryID,
		arg.CategoryID,
		arg.SetAmount,
		arg.Amount,
		arg.UpdatedBy,
		arg.ID,
	)
	var i Selling
	err := row.Scan(
		&i.ID,
		&i.PartnerID,
		&i.CategoryID,
		&i.Amount,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.CreatedBy,
		&i.UpdatedBy,
		&i.DeletedBy,
	)
	return i, err
}
