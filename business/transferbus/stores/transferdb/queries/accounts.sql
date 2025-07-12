-- name: CreateAccount :exec
INSERT INTO accounts (account_id, created_date, last_modified_date)
VALUES (@account_id, @created_date, @last_modified_date);

-- name: GetAccount :one
SELECT * FROM accounts WHERE account_id = @account_id;

-- name: GetAccounts :many
SELECT * FROM accounts where account_id = any(@account_ids::bigint[]);