-- name: CreateTransaction :one
INSERT INTO "transactions" (
    tx_id, "bill_id", cust_name, amount, "admin", tot_amount, fee_partner, fee_ppob, valid_from, valid_to,
    cat_id, cat_name, prod_id, prod_name, partner_id, partner_name, provider_id, provider_name,
    status, req_inq_params, res_inq_params, req_pay_params, res_pay_params,
    req_cmt_params, res_cmt_params, req_adv_params, res_adv_params, req_rev_params, res_rev_params,
    created_by, ref_id, created_at, first_balance, last_balance, payment_type
) values (
             $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16,
            $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, now(), $32, $33, $34
         ) RETURNING *;

-- name: GetTransaction :one
SELECT * FROM "transactions"
WHERE id = $1 AND deleted_at is null LIMIT 1;

-- name: GetTransactionByTxID :one
SELECT * FROM "transactions"
WHERE tx_id = $1 AND deleted_at is null LIMIT 1;

-- name: GetTransactionByRefID :one
SELECT * FROM "transactions"
WHERE ref_id = $1
AND status = '0'
AND partner_id = $2
AND to_char(created_at,'YYYY-MM-DD') = to_char(now(),'YYYY-MM-DD')
AND deleted_at is null
LIMIT 1;

-- name: GetTransactionPending :one
SELECT * FROM "transactions"
WHERE bill_id = $1
  AND status = '4'
  AND to_char(created_at,'YYYY-MM-DD') = to_char(now(),'YYYY-MM-DD')
  AND deleted_at is null
    LIMIT 1;

-- name: ListTransaction :many
SELECT * FROM "transactions"
WHERE deleted_at is null
AND DATE(created_at) >= to_date(@from_date,'YYYY-MM-DD')
AND DATE(created_at) <= to_date(@to_date,'YYYY-MM-DD')
AND (CASE WHEN @is_status::bool THEN status = @status ELSE TRUE END)
AND (CASE WHEN @is_cat::bool THEN cat_id = @cat_id ELSE TRUE END)
AND (CASE WHEN @is_partner::bool THEN partner_id = @partner_id ELSE TRUE END)
AND (CASE WHEN @is_created::bool THEN created_by = @created_by ELSE TRUE END)
AND (CASE WHEN @is_type::bool THEN payment_type = @payment_type ELSE TRUE END)
ORDER BY created_at DESC
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
    "first_balance" = CASE
                   WHEN sqlc.arg(set_first_balance)::bool
                THEN sqlc.arg(first_balance)
                   ELSE first_balance
        END,
    "last_balance" = CASE
                   WHEN sqlc.arg(set_last_balance)::bool
                THEN sqlc.arg(last_balance)
                   ELSE last_balance
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
    sn = CASE
            WHEN sqlc.arg(set_sn)::bool
                THEN sqlc.arg(sn)
            ELSE sn
        END,
    add_info1 = CASE
                         WHEN sqlc.arg(set_info1)::bool
                    THEN sqlc.arg(add_info1)
                         ELSE add_info1
        END,
    add_info2 = CASE
                         WHEN sqlc.arg(set_info2)::bool
                    THEN sqlc.arg(add_info2)
                         ELSE add_info2
        END,
    add_info3 = CASE
                         WHEN sqlc.arg(set_info3)::bool
                    THEN sqlc.arg(add_info3)
                         ELSE add_info3
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