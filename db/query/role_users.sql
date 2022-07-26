-- name: CreateRoleUser :one
INSERT INTO role_users (
    role_id, user_id, created_by
) values (
             $1, $2, $3
         ) RETURNING *;

-- name: GetRoleUserByID :one
SELECT * FROM role_users
WHERE id = $1 AND deleted_at is null LIMIT 1;

-- name: GetRoleUserByUserID :many
SELECT * FROM role_users
WHERE
    user_id = $1 AND deleted_at is null
LIMIT $2
OFFSET $3;

-- name: GetRoleUserByRoleID :many
SELECT * FROM role_users
WHERE
    role_id = $1 AND deleted_at is null
LIMIT $2
OFFSET $3;

-- name: ListRoleUser :many
SELECT * FROM role_users
WHERE deleted_at is null
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateRoleUser :one
UPDATE role_users
SET
    user_id = $2,
    role_id = $3,
    updated_by = $4,
    updated_at = now()
WHERE
    id = $1
    RETURNING *;

-- name: UpdateInactiveRoleUser :one
UPDATE role_users SET deleted_by = $2, deleted_at = now() WHERE id = $1
RETURNING *;

-- name: DeleteRoleUser :exec
DELETE FROM role_users
WHERE id = $1;