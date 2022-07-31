-- name: CreateRole :one
INSERT INTO roles (
    name, level, created_by
) values (
    $1, $2, $3
) RETURNING *;

-- name: GetRole :one
SELECT * FROM roles
WHERE id = $1 AND deleted_at is null LIMIT 1;

-- name: ListRoleWithDelete :many
SELECT * FROM roles
ORDER BY name
LIMIT $1
OFFSET $2;

-- name: ListRole :many
SELECT * FROM roles
WHERE deleted_at is null
ORDER BY name
LIMIT $1
OFFSET $2;

-- name: UpdateRole :one
UPDATE roles
SET
    name = $2,
    level = $3,
    updated_by = $4,
    updated_at = now()
WHERE
    id = $1
RETURNING *;

-- name: UpdateInactiveRole :one
UPDATE roles SET deleted_by = $2, deleted_at = now() WHERE id = $1
RETURNING *;

-- name: DeleteRole :exec
DELETE FROM roles
WHERE id = $1;