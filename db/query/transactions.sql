-- name: CreateTransaction :one
INSERT INTO "transactions" (
    tx_id, "bill_id", cust_name, amount, "admin", tot_amount, fee_partner, fee_ppob, valid_from, valid_to,
    cat_id, cat_name, prod_id, prod_name, partner_id, partner_name, provider_id, provider_name,
    status, req_inq_params, res_inq_params, req_pay_params, res_pay_params,
    req_cmt_params, res_cmt_params, req_adv_params, res_adv_params, req_rev_params, res_rev_params,
    created_by, created_at
) values (
             $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16,
            $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, now()
         ) RETURNING *;

-- name: GetTransaction :one
SELECT * FROM "transactions"
WHERE id = $1 LIMIT 1;

-- name: GetTransactionByTxID :one
SELECT * FROM "transactions"
WHERE tx_id = $1 LIMIT 1;

-- name: ListTransaction :many
SELECT * FROM "transactions"
ORDER BY created_at
LIMIT $1
OFFSET $2;

-- name: UpdateTransaction :one
UPDATE "transactions"
SET
    "status" = CASE
            WHEN sqlc.arg(set_status)::bool
                THEN sqlc.arg(status)
            ELSE status
            END,
    res_inq_params = CASE
               WHEN sqlc.arg(set_res_inq_params)::bool
                THEN sqlc.arg(res_inq_params)
               ELSE res_inq_params
        END,
    req_pay_params = CASE
                 WHEN sqlc.arg(set_req_pay_params)::bool
                    THEN sqlc.arg(req_pay_params)
                 ELSE req_pay_params
        END,
    res_pay_params = CASE
                 WHEN sqlc.arg(set_res_pay_params)::bool
                    THEN sqlc.arg(res_pay_params)
                 ELSE res_pay_params
        END,
    req_cmt_params = CASE
                    WHEN sqlc.arg(set_req_cmt_params)::bool
                    THEN sqlc.arg(req_cmt_params)
                    ELSE req_cmt_params
        END,
    res_cmt_params = CASE
                    WHEN sqlc.arg(set_res_cmt_params)::bool
                    THEN sqlc.arg(res_cmt_params)
                    ELSE res_cmt_params
        END,
    req_adv_params = CASE
                     WHEN sqlc.arg(set_req_adv_params)::bool
                    THEN sqlc.arg(req_adv_params)
                     ELSE req_adv_params
        END,
    res_adv_params = CASE
                   WHEN sqlc.arg(set_res_adv_params)::bool
                    THEN sqlc.arg(res_adv_params)
                   ELSE res_adv_params
        END,
    req_rev_params = CASE
                   WHEN sqlc.arg(set_req_rev_params)::bool
                    THEN sqlc.arg(req_rev_params)
                   ELSE req_rev_params
        END,
    updated_by = sqlc.arg(updated_by),
    updated_at = now()
WHERE
    id = sqlc.arg(id)
RETURNING *;

-- name: UpdateInactiveTransaction :one
UPDATE "transactions" SET deleted_by = $2, deleted_at = now() WHERE id = $1
RETURNING *;

-- name: DeleteTransaction :exec
DELETE FROM "transactions"
WHERE id = $1;