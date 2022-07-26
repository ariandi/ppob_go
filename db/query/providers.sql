-- name: CreateProvider :one
INSERT INTO providers (
    name, "user", secret, add_info1, add_info2, valid_from, valid_to, base_url, "method", inq,
    pay, adv, cmt, rev, status, created_by, created_at
) values (
             $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, now()
         ) RETURNING *;

-- name: GetProvider :one
SELECT * FROM providers
WHERE id = $1 AND deleted_at is null LIMIT 1;

-- name: ListProvider :many
SELECT * FROM providers
WHERE deleted_at is null
ORDER BY name
LIMIT $1
OFFSET $2;

-- name: UpdateProvider :one
UPDATE providers
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
    base_url = CASE
                   WHEN sqlc.arg(set_base_url)::bool
                    THEN sqlc.arg(base_url)
                   ELSE base_url
        END,
    method = CASE
                   WHEN sqlc.arg(set_method)::bool
                    THEN sqlc.arg(method)
                   ELSE method
        END,
    inq = CASE
                 WHEN sqlc.arg(set_inq)::bool
                    THEN sqlc.arg(inq)
                 ELSE inq
        END,
    pay = CASE
              WHEN sqlc.arg(set_pay)::bool
                    THEN sqlc.arg(pay)
              ELSE pay
        END,
    adv = CASE
              WHEN sqlc.arg(set_adv)::bool
                    THEN sqlc.arg(adv)
              ELSE adv
        END,
    cmt = CASE
              WHEN sqlc.arg(set_cmt)::bool
                    THEN sqlc.arg(cmt)
              ELSE cmt
        END,
    rev = CASE
              WHEN sqlc.arg(set_rev)::bool
                    THEN sqlc.arg(rev)
              ELSE rev
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

-- name: UpdateInactiveProvider :one
UPDATE providers SET deleted_by = $2, deleted_at = now() WHERE id = $1
RETURNING *;

-- name: DeleteProvider :exec
DELETE FROM providers
WHERE id = $1;