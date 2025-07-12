-- name: CreateTransaction :exec
INSERT INTO transactions (account_id, amount, created_date)
VALUES (@account_id, @amount, @created_date);