-- name: CreatePartner :one
INSERT INTO "partners" (
    name, "user", secret, add_info1, add_info2, valid_from, valid_to, payment_type, status,
    created_by, created_at
) values (
             $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, now()
         ) RETURNING *;

-- name: GetPartner :one
SELECT * FROM "partners"
WHERE id = $1 AND deleted_at is null LIMIT 1;

-- name: GetPartnerByParams :one
SELECT * FROM "partners"
WHERE deleted_at is null
AND (CASE WHEN @is_name::bool THEN name = @name ELSE TRUE END)
AND (CASE WHEN @is_user::bool THEN user = @user_params ELSE TRUE END)
LIMIT 1;

-- name: ListPartner :many
SELECT * FROM "partners"
WHERE deleted_at is null
ORDER BY name
LIMIT $1
OFFSET $2;

-- name: UpdatePartner :one
UPDATE "partners"
SET
    "name" = CASE
            WHEN sqlc.arg(set_name)::bool
                THEN sqlc.arg(name)
            ELSE name
            END,
    "user" = CASE
               WHEN sqlc.arg(set_user)::bool
                THEN sqlc.arg(user_params)
               ELSE "user"
        END,
    secret = CASE
                 WHEN sqlc.arg(set_secret)::bool
                    THEN sqlc.arg(secret)
                 ELSE secret
        END,
    add_info1 = CASE
                 WHEN sqlc.arg(set_add_info1)::bool
                    THEN sqlc.arg(add_info1)
                 ELSE add_info1
        END,
    add_info2 = CASE
                    WHEN sqlc.arg(set_add_info2)::bool
                    THEN sqlc.arg(add_info2)
                    ELSE add_info2
        END,
    valid_from = CASE
                    WHEN sqlc.arg(set_valid_from)::bool
                    THEN sqlc.arg(valid_from)
                    ELSE valid_from
        END,
    valid_to = CASE
                     WHEN sqlc.arg(set_valid_to)::bool
                    THEN sqlc.arg(valid_to)
                     ELSE valid_to
        END,
    payment_type = CASE
                   WHEN sqlc.arg(set_payment_type)::bool
                    THEN sqlc.arg(payment_type)
                   ELSE payment_type
        END,
    status = CASE
                      WHEN sqlc.arg(set_status)::bool
                    THEN sqlc.arg(status)
                      ELSE status
        END,
    updated_by = sqlc.arg(updated_by),
    updated_at = now()
WHERE
    id = sqlc.arg(id)
RETURNING *;

-- name: UpdateInactivePartner :one
UPDATE "partners" SET deleted_by = $2, deleted_at = now() WHERE id = $1
RETURNING *;

-- name: DeletePartner :exec
DELETE FROM "partners"
WHERE id = $1;