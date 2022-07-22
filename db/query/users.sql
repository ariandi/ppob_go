-- name: CreateUser :one
INSERT INTO users (
    name, email, username, password, balance, phone, identity_number, created_by
) values (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: GetUser :one
SELECT users.*, roles.id AS role_id, roles.name FROM users
LEFT JOIN role_users on role_users.user_id = users.id
LEFT JOIN roles on roles.id = role_users.role_id
WHERE users.id = $1 LIMIT 1;

-- name: GetUserByUsername :one
SELECT users.*, roles.id AS role_id, roles.name FROM users
LEFT JOIN role_users on role_users.user_id = users.id
LEFT JOIN roles on roles.id = role_users.role_id
WHERE username = $1 LIMIT 1;

-- name: ListUser :many
SELECT users.*, roles.id AS role_id, roles.name
FROM users
LEFT JOIN role_users on role_users.user_id = users.id
LEFT JOIN roles on roles.id = role_users.role_id
ORDER BY users.name
LIMIT $1
OFFSET $2;

-- name: UpdateUser :one
UPDATE
    users
SET
    name = CASE
                   WHEN sqlc.arg(set_name)::bool
                    THEN sqlc.arg(name)
                   ELSE name
            END,
    phone = CASE
                 WHEN sqlc.arg(set_phone)::bool
                    THEN sqlc.arg(phone)
                 ELSE phone
            END,
    identity_number = CASE
                WHEN sqlc.arg(set_identity_number)::bool
                    THEN sqlc.arg(identity_number)
                ELSE identity_number
            END,
    password = CASE
                WHEN sqlc.arg(set_password)::bool
                    THEN sqlc.arg(password)
                ELSE password
            END,
    balance = CASE
               WHEN sqlc.arg(set_balance)::bool
                THEN sqlc.arg(balance)
               ELSE balance
            END,
    updated_by = sqlc.arg(updated_by),
    updated_at = now()
WHERE
    id = sqlc.arg(id)
RETURNING *;

-- name: UpdateInactiveUser :one
UPDATE users SET deleted_by = $2, deleted_at = now() WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;