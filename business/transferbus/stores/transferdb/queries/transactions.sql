-- name: GetBalance :one
SELECT COALESCE(SUM(amount), 0)::NUMERIC AS balance 
FROM transactions WHERE account_id = @account_id;

-- name: CreateTransaction :exec
INSERT INTO transactions (account_id, amount, created_date)
VALUES (@account_id, @amount, @created_date);