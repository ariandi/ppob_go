// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
// source: categories.sql

package db

import (
	"context"
	"database/sql"
)

const createCategory = `-- name: CreateCategory :one
INSERT INTO categories (
    name, parent, created_by, up_selling
) values (
             $1, 0, $2, $3
         ) RETURNING id, name, up_selling, parent, created_at, updated_at, deleted_at, created_by, updated_by, deleted_by
`

type CreateCategoryParams struct {
	Name      string         `json:"name"`
	CreatedBy sql.NullInt64  `json:"created_by"`
	UpSelling sql.NullString `json:"up_selling"`
}

func (q *Queries) CreateCategory(ctx context.Context, arg CreateCategoryParams) (Category, error) {
	row := q.db.QueryRowContext(ctx, createCategory, arg.Name, arg.CreatedBy, arg.UpSelling)
	var i Category
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.UpSelling,
		&i.Parent,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.CreatedBy,
		&i.UpdatedBy,
		&i.DeletedBy,
	)
	return i, err
}

const deleteCategories = `-- name: DeleteCategories :exec
DELETE FROM categories
WHERE id = $1
`

func (q *Queries) DeleteCategories(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteCategories, id)
	return err
}

const getCategory = `-- name: GetCategory :one
SELECT id, name, up_selling, parent, created_at, updated_at, deleted_at, created_by, updated_by, deleted_by FROM categories
WHERE id = $1 AND deleted_at is null LIMIT 1
`

func (q *Queries) GetCategory(ctx context.Context, id int64) (Category, error) {
	row := q.db.QueryRowContext(ctx, getCategory, id)
	var i Category
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.UpSelling,
		&i.Parent,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.CreatedBy,
		&i.UpdatedBy,
		&i.DeletedBy,
	)
	return i, err
}

const listCategory = `-- name: ListCategory :many
SELECT id, name, up_selling, parent, created_at, updated_at, deleted_at, created_by, updated_by, deleted_by FROM categories
WHERE deleted_at is null
ORDER BY name
LIMIT $1
OFFSET $2
`

type ListCategoryParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListCategory(ctx context.Context, arg ListCategoryParams) ([]Category, error) {
	rows, err := q.db.QueryContext(ctx, listCategory, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Category{}
	for rows.Next() {
		var i Category
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.UpSelling,
			&i.Parent,
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

const updateCategory = `-- name: UpdateCategory :one
UPDATE categories
SET
    name = CASE
            WHEN $1::bool
                THEN $2
            ELSE name
            END,
    parent = CASE
              WHEN $3::bool
                THEN $4
              ELSE parent
              END,
    up_selling = CASE
                 WHEN $5::bool
                    THEN $6
                 ELSE up_selling
                END,
    updated_by = $7,
    updated_at = now()
WHERE
        id = $8
    RETURNING id, name, up_selling, parent, created_at, updated_at, deleted_at, created_by, updated_by, deleted_by
`

type UpdateCategoryParams struct {
	SetName      bool           `json:"set_name"`
	Name         string         `json:"name"`
	SetParent    bool           `json:"set_parent"`
	Parent       int64          `json:"parent"`
	SetUpSelling bool           `json:"set_up_selling"`
	UpSelling    sql.NullString `json:"up_selling"`
	UpdatedBy    sql.NullInt64  `json:"updated_by"`
	ID           int64          `json:"id"`
}

func (q *Queries) UpdateCategory(ctx context.Context, arg UpdateCategoryParams) (Category, error) {
	row := q.db.QueryRowContext(ctx, updateCategory,
		arg.SetName,
		arg.Name,
		arg.SetParent,
		arg.Parent,
		arg.SetUpSelling,
		arg.UpSelling,
		arg.UpdatedBy,
		arg.ID,
	)
	var i Category
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.UpSelling,
		&i.Parent,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.CreatedBy,
		&i.UpdatedBy,
		&i.DeletedBy,
	)
	return i, err
}

const updateInactiveCategory = `-- name: UpdateInactiveCategory :one
UPDATE categories SET deleted_by = $2, deleted_at = now() WHERE id = $1
    RETURNING id, name, up_selling, parent, created_at, updated_at, deleted_at, created_by, updated_by, deleted_by
`

type UpdateInactiveCategoryParams struct {
	ID        int64         `json:"id"`
	DeletedBy sql.NullInt64 `json:"deleted_by"`
}

func (q *Queries) UpdateInactiveCategory(ctx context.Context, arg UpdateInactiveCategoryParams) (Category, error) {
	row := q.db.QueryRowContext(ctx, updateInactiveCategory, arg.ID, arg.DeletedBy)
	var i Category
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.UpSelling,
		&i.Parent,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.CreatedBy,
		&i.UpdatedBy,
		&i.DeletedBy,
	)
	return i, err
}
