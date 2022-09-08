-- name: CreateMediaStorage :one
INSERT INTO madiastorages (
    sec_id, tab_id, name, type, content, created_by
) values (
             $1, $2, $3, $4, $5, $6
         ) RETURNING *;

-- name: GetMediaStorage :one
SELECT * FROM madiastorages
WHERE deleted_at is null
AND (CASE WHEN @is_id::bool THEN id = @id ELSE TRUE END)
AND (CASE WHEN @is_sec::bool THEN sec_id = @sec_id ELSE TRUE END)
AND (CASE WHEN @is_tab::bool THEN tab_id = @tab_id ELSE TRUE END)
LIMIT 1;

-- name: ListMediaStorage :many
SELECT * FROM madiastorages
WHERE deleted_at is null
AND (CASE WHEN @is_sec::bool THEN sec_id = @sec_id ELSE TRUE END)
AND (CASE WHEN @is_tab::bool THEN tab_id = @tab_id ELSE TRUE END)
ORDER BY name
LIMIT $1
OFFSET $2;

-- name: UpdateMediaStorage :one
UPDATE madiastorages
SET
    name = CASE
            WHEN sqlc.arg(set_name)::bool
                THEN sqlc.arg(name)
            ELSE name
            END,
    type = CASE
              WHEN sqlc.arg(set_type)::bool
                THEN sqlc.arg(type)
              ELSE type
              END,
    updated_by = sqlc.arg(updated_by),
    updated_at = now()
WHERE
id = sqlc.arg(id)
RETURNING *;

-- name: UpdateInactiveMediaStorage :one
UPDATE madiastorages SET deleted_by = $2, deleted_at = now() WHERE id = $1
    RETURNING *;

-- name: DeleteMediaStorage :exec
DELETE FROM madiastorages
WHERE id = $1;