// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
// source: users.sql

package db

import (
	"context"
	"database/sql"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
    name, email, username, password, balance, phone, identity_number, created_by, bank_code
) values (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING id, name, email, username, bank_code, password, balance, phone, identity_number, verified_at, created_at, updated_at, deleted_at, created_by, updated_by, deleted_by
`

type CreateUserParams struct {
	Name           string         `json:"name"`
	Email          string         `json:"email"`
	Username       string         `json:"username"`
	Password       sql.NullString `json:"password"`
	Balance        sql.NullString `json:"balance"`
	Phone          string         `json:"phone"`
	IdentityNumber string         `json:"identity_number"`
	CreatedBy      sql.NullInt64  `json:"created_by"`
	BankCode       sql.NullInt64  `json:"bank_code"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.Name,
		arg.Email,
		arg.Username,
		arg.Password,
		arg.Balance,
		arg.Phone,
		arg.IdentityNumber,
		arg.CreatedBy,
		arg.BankCode,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Username,
		&i.BankCode,
		&i.Password,
		&i.Balance,
		&i.Phone,
		&i.IdentityNumber,
		&i.VerifiedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.CreatedBy,
		&i.UpdatedBy,
		&i.DeletedBy,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteUser, id)
	return err
}

const getUser = `-- name: GetUser :one
SELECT users.id, users.name, users.email, users.username, users.bank_code, users.password, users.balance, users.phone, users.identity_number, users.verified_at, users.created_at, users.updated_at, users.deleted_at, users.created_by, users.updated_by, users.deleted_by, roles.id AS role_id, roles.name FROM users
LEFT JOIN role_users on role_users.user_id = users.id
LEFT JOIN roles on roles.id = role_users.role_id
WHERE users.id = $1 AND users.deleted_at is null LIMIT 1
`

type GetUserRow struct {
	ID             int64          `json:"id"`
	Name           string         `json:"name"`
	Email          string         `json:"email"`
	Username       string         `json:"username"`
	BankCode       sql.NullInt64  `json:"bank_code"`
	Password       sql.NullString `json:"password"`
	Balance        sql.NullString `json:"balance"`
	Phone          string         `json:"phone"`
	IdentityNumber string         `json:"identity_number"`
	VerifiedAt     sql.NullTime   `json:"verified_at"`
	CreatedAt      sql.NullTime   `json:"created_at"`
	UpdatedAt      sql.NullTime   `json:"updated_at"`
	DeletedAt      sql.NullTime   `json:"deleted_at"`
	CreatedBy      sql.NullInt64  `json:"created_by"`
	UpdatedBy      sql.NullInt64  `json:"updated_by"`
	DeletedBy      sql.NullInt64  `json:"deleted_by"`
	RoleID         sql.NullInt64  `json:"role_id"`
	Name_2         sql.NullString `json:"name_2"`
}

func (q *Queries) GetUser(ctx context.Context, id int64) (GetUserRow, error) {
	row := q.db.QueryRowContext(ctx, getUser, id)
	var i GetUserRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Username,
		&i.BankCode,
		&i.Password,
		&i.Balance,
		&i.Phone,
		&i.IdentityNumber,
		&i.VerifiedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.CreatedBy,
		&i.UpdatedBy,
		&i.DeletedBy,
		&i.RoleID,
		&i.Name_2,
	)
	return i, err
}

const getUserByUsername = `-- name: GetUserByUsername :one
SELECT users.id, users.name, users.email, users.username, users.bank_code, users.password, users.balance, users.phone, users.identity_number, users.verified_at, users.created_at, users.updated_at, users.deleted_at, users.created_by, users.updated_by, users.deleted_by, roles.id AS role_id, roles.name FROM users
LEFT JOIN role_users on role_users.user_id = users.id
LEFT JOIN roles on roles.id = role_users.role_id
WHERE username = $1 AND users.deleted_at is null LIMIT 1
`

type GetUserByUsernameRow struct {
	ID             int64          `json:"id"`
	Name           string         `json:"name"`
	Email          string         `json:"email"`
	Username       string         `json:"username"`
	BankCode       sql.NullInt64  `json:"bank_code"`
	Password       sql.NullString `json:"password"`
	Balance        sql.NullString `json:"balance"`
	Phone          string         `json:"phone"`
	IdentityNumber string         `json:"identity_number"`
	VerifiedAt     sql.NullTime   `json:"verified_at"`
	CreatedAt      sql.NullTime   `json:"created_at"`
	UpdatedAt      sql.NullTime   `json:"updated_at"`
	DeletedAt      sql.NullTime   `json:"deleted_at"`
	CreatedBy      sql.NullInt64  `json:"created_by"`
	UpdatedBy      sql.NullInt64  `json:"updated_by"`
	DeletedBy      sql.NullInt64  `json:"deleted_by"`
	RoleID         sql.NullInt64  `json:"role_id"`
	Name_2         sql.NullString `json:"name_2"`
}

func (q *Queries) GetUserByUsername(ctx context.Context, username string) (GetUserByUsernameRow, error) {
	row := q.db.QueryRowContext(ctx, getUserByUsername, username)
	var i GetUserByUsernameRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Username,
		&i.BankCode,
		&i.Password,
		&i.Balance,
		&i.Phone,
		&i.IdentityNumber,
		&i.VerifiedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.CreatedBy,
		&i.UpdatedBy,
		&i.DeletedBy,
		&i.RoleID,
		&i.Name_2,
	)
	return i, err
}

const listUser = `-- name: ListUser :many
SELECT users.id, users.name, users.email, users.username, users.bank_code, users.password, users.balance, users.phone, users.identity_number, users.verified_at, users.created_at, users.updated_at, users.deleted_at, users.created_by, users.updated_by, users.deleted_by, roles.id AS role_id, roles.name
FROM users
LEFT JOIN role_users on role_users.user_id = users.id
LEFT JOIN roles on roles.id = role_users.role_id
WHERE users.deleted_at is null
AND (CASE WHEN $3::bool THEN role_users.role_id = $4 ELSE TRUE END)
ORDER BY users.name
LIMIT $1
OFFSET $2
`

type ListUserParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
	IsRole bool  `json:"is_role"`
	RoleID int64 `json:"role_id"`
}

type ListUserRow struct {
	ID             int64          `json:"id"`
	Name           string         `json:"name"`
	Email          string         `json:"email"`
	Username       string         `json:"username"`
	BankCode       sql.NullInt64  `json:"bank_code"`
	Password       sql.NullString `json:"password"`
	Balance        sql.NullString `json:"balance"`
	Phone          string         `json:"phone"`
	IdentityNumber string         `json:"identity_number"`
	VerifiedAt     sql.NullTime   `json:"verified_at"`
	CreatedAt      sql.NullTime   `json:"created_at"`
	UpdatedAt      sql.NullTime   `json:"updated_at"`
	DeletedAt      sql.NullTime   `json:"deleted_at"`
	CreatedBy      sql.NullInt64  `json:"created_by"`
	UpdatedBy      sql.NullInt64  `json:"updated_by"`
	DeletedBy      sql.NullInt64  `json:"deleted_by"`
	RoleID         sql.NullInt64  `json:"role_id"`
	Name_2         sql.NullString `json:"name_2"`
}

func (q *Queries) ListUser(ctx context.Context, arg ListUserParams) ([]ListUserRow, error) {
	rows, err := q.db.QueryContext(ctx, listUser,
		arg.Limit,
		arg.Offset,
		arg.IsRole,
		arg.RoleID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListUserRow{}
	for rows.Next() {
		var i ListUserRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Email,
			&i.Username,
			&i.BankCode,
			&i.Password,
			&i.Balance,
			&i.Phone,
			&i.IdentityNumber,
			&i.VerifiedAt,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
			&i.CreatedBy,
			&i.UpdatedBy,
			&i.DeletedBy,
			&i.RoleID,
			&i.Name_2,
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

const listUserCount = `-- name: ListUserCount :one
SELECT COUNT(id),
CASE WHEN COUNT(id) > 0 THEN SUM(balance) ELSE 0 END AS sum
FROM users
WHERE deleted_at is null
`

type ListUserCountRow struct {
	Count int64  `json:"count"`
	Sum   string `json:"sum"`
}

func (q *Queries) ListUserCount(ctx context.Context) (ListUserCountRow, error) {
	row := q.db.QueryRowContext(ctx, listUserCount)
	var i ListUserCountRow
	err := row.Scan(&i.Count, &i.Sum)
	return i, err
}

const updateInactiveUser = `-- name: UpdateInactiveUser :one
UPDATE users SET deleted_by = $2, deleted_at = now() WHERE id = $1
RETURNING id, name, email, username, bank_code, password, balance, phone, identity_number, verified_at, created_at, updated_at, deleted_at, created_by, updated_by, deleted_by
`

type UpdateInactiveUserParams struct {
	ID        int64         `json:"id"`
	DeletedBy sql.NullInt64 `json:"deleted_by"`
}

func (q *Queries) UpdateInactiveUser(ctx context.Context, arg UpdateInactiveUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateInactiveUser, arg.ID, arg.DeletedBy)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Username,
		&i.BankCode,
		&i.Password,
		&i.Balance,
		&i.Phone,
		&i.IdentityNumber,
		&i.VerifiedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.CreatedBy,
		&i.UpdatedBy,
		&i.DeletedBy,
	)
	return i, err
}

const updateUser = `-- name: UpdateUser :one
UPDATE
    users
SET
    name = CASE
                   WHEN $1::bool
                    THEN $2
                   ELSE name
            END,
    phone = CASE
                 WHEN $3::bool
                    THEN $4
                 ELSE phone
            END,
    identity_number = CASE
                WHEN $5::bool
                    THEN $6
                ELSE identity_number
            END,
    password = CASE
                WHEN $7::bool
                    THEN $8
                ELSE password
            END,
    balance = CASE
               WHEN $9::bool
                THEN $10
               ELSE balance
            END,
    bank_code = CASE
                  WHEN $11::bool
                THEN $12
                  ELSE bank_code
        END,
    email = CASE
                    WHEN $13::bool
                THEN $14
                    ELSE email
        END,
    updated_by = $15,
    updated_at = now()
WHERE
    id = $16
RETURNING id, name, email, username, bank_code, password, balance, phone, identity_number, verified_at, created_at, updated_at, deleted_at, created_by, updated_by, deleted_by
`

type UpdateUserParams struct {
	SetName           bool           `json:"set_name"`
	Name              string         `json:"name"`
	SetPhone          bool           `json:"set_phone"`
	Phone             string         `json:"phone"`
	SetIdentityNumber bool           `json:"set_identity_number"`
	IdentityNumber    string         `json:"identity_number"`
	SetPassword       bool           `json:"set_password"`
	Password          sql.NullString `json:"password"`
	SetBalance        bool           `json:"set_balance"`
	Balance           sql.NullString `json:"balance"`
	SetBankCode       bool           `json:"set_bank_code"`
	BankCode          sql.NullInt64  `json:"bank_code"`
	SetEmail          bool           `json:"set_email"`
	Email             string         `json:"email"`
	UpdatedBy         sql.NullInt64  `json:"updated_by"`
	ID                int64          `json:"id"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUser,
		arg.SetName,
		arg.Name,
		arg.SetPhone,
		arg.Phone,
		arg.SetIdentityNumber,
		arg.IdentityNumber,
		arg.SetPassword,
		arg.Password,
		arg.SetBalance,
		arg.Balance,
		arg.SetBankCode,
		arg.BankCode,
		arg.SetEmail,
		arg.Email,
		arg.UpdatedBy,
		arg.ID,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Username,
		&i.BankCode,
		&i.Password,
		&i.Balance,
		&i.Phone,
		&i.IdentityNumber,
		&i.VerifiedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.CreatedBy,
		&i.UpdatedBy,
		&i.DeletedBy,
	)
	return i, err
}
