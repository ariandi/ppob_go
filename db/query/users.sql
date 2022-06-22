-- name: CreateUser :one
INSERT INTO users (
    name, email, username, password, balance, phone, identity_number, created_by
) values (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: ListUser :many
SELECT * FROM users
ORDER BY name
LIMIT $1
OFFSET $2;

-- name: UpdateUser :one
UPDATE
    users
SET
    name = $2,
    password = $3,
    updated_by = $4,
    updated_at = now()
WHERE
    id = $1
RETURNING *;

-- name: UpdateInactiveUser :one
UPDATE users SET deleted_by = $2, deleted_at = now() WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;