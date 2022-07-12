-- name: CreateCategory :one
INSERT INTO categories (
    name, parent, created_by
) values (
             $1, 0, $2
         ) RETURNING *;

-- name: GetCategory :one
SELECT * FROM categories
WHERE id = $1 LIMIT 1;

-- name: ListCategory :many
SELECT * FROM categories
ORDER BY name
LIMIT $1
OFFSET $2;

-- name: UpdateCategory :one
UPDATE categories
SET
    name = CASE
            WHEN sqlc.arg(set_name)::bool
                THEN sqlc.arg(name)
            ELSE name
            END,
    parent = CASE
              WHEN sqlc.arg(set_parent)::bool
                THEN sqlc.arg(parent)
              ELSE parent
              END,
    updated_by = sqlc.arg(updated_by),
    updated_at = now()
WHERE
        id = sqlc.arg(id)
    RETURNING *;

-- name: UpdateInactiveCategory :one
UPDATE categories SET deleted_by = $2, deleted_at = now() WHERE id = $1
    RETURNING *;

-- name: DeleteCategories :exec
DELETE FROM categories
WHERE id = $1;