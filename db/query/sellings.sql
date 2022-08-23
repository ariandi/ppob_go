-- name: CreateSelling :one
INSERT INTO sellings (
    partner_id, category_id, amount, created_by
) values (
             $1, $2, $3, $4
         ) RETURNING *;

-- name: GetSelling :one
SELECT * FROM sellings
WHERE id = $1 AND deleted_at is null LIMIT 1;

-- name: ListSellingByParams :many
SELECT * FROM sellings
WHERE deleted_at is null
AND (CASE WHEN @is_partner::bool THEN partner_id = @partner_id ELSE TRUE END)
AND (CASE WHEN @is_category::bool THEN category_id = @category_id ELSE TRUE END)
ORDER BY partner_id, category_id
LIMIT $1
OFFSET $2;

-- name: ListSelling :many
SELECT * FROM sellings
WHERE deleted_at is null
ORDER BY partner_id, category_id
LIMIT $1
OFFSET $2;

-- name: UpdateSelling :one
UPDATE sellings
SET
    partner_id = CASE
               WHEN sqlc.arg(set_partner_id)::bool
                THEN sqlc.arg(partner_id)
               ELSE partner_id
        END,
    category_id = CASE
                 WHEN sqlc.arg(set_category_id)::bool
                THEN sqlc.arg(category_id)
                 ELSE category_id
        END,
    amount = CASE
                WHEN sqlc.arg(set_amount)::bool
                THEN sqlc.arg(amount)
                ELSE amount
        END,
    updated_by = sqlc.arg(updated_by),
    updated_at = now()
WHERE
id = sqlc.arg(id)
RETURNING *;

-- name: UpdateInactiveSelling :one
UPDATE sellings SET deleted_by = $2, deleted_at = now() WHERE id = $1
RETURNING *;

-- name: DeleteSelling :exec
DELETE FROM sellings
WHERE id = $1;