-- name: CreateProduct :one
INSERT INTO products (
    cat_id, name, amount, provider_id, provider_code, status, parent, created_by
) values (
             $1, $2, $3, $4, $5, $6, $7, $8
         ) RETURNING *;

-- name: GetProduct :one
SELECT * FROM products
WHERE id = $1 AND deleted_at is null LIMIT 1;

-- name: ListProduct :many
SELECT * FROM products
WHERE deleted_at is null
AND (CASE WHEN @is_cat::bool THEN cat_id = @cat_id ELSE TRUE END)
AND (CASE WHEN @is_prov::bool THEN provider_id = @provider_id ELSE TRUE END)
ORDER BY name
LIMIT $1
OFFSET $2;

-- name: ListProductByCatID :many
SELECT * FROM products
WHERE cat_id = $1
ORDER BY name
LIMIT $1
OFFSET $2;

-- name: UpdateProduct :one
UPDATE products
SET
    name = CASE
            WHEN sqlc.arg(set_name)::bool
                THEN sqlc.arg(name)
            ELSE name
            END,
    cat_id = CASE
               WHEN sqlc.arg(set_cat)::bool
                THEN sqlc.arg(cat_id)
               ELSE cat_id
                END,
    "amount" = CASE
                 WHEN sqlc.arg(set_amount)::bool
                    THEN sqlc.arg(amount)
                 ELSE amount::DECIMAL
                END,
    provider_id = CASE
                 WHEN sqlc.arg(set_provider)::bool
                    THEN sqlc.arg(provider_id)
                 ELSE provider_id
                END,
    provider_code = CASE
                 WHEN sqlc.arg(set_provider_code)::bool
                    THEN sqlc.arg(provider_code)
                 ELSE provider_code
                END,
    status = CASE
                      WHEN sqlc.arg(set_status)::bool
                    THEN sqlc.arg(status)
                      ELSE status
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

-- name: UpdateInactiveProduct :one
UPDATE products SET deleted_by = $2, deleted_at = now() WHERE id = $1
RETURNING *;

-- name: DeleteProduct :exec
DELETE FROM products
WHERE id = $1;