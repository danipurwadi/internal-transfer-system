-- name: CreateAccount :exec
INSERT INTO accounts (account_id, balance, created_date, last_modified_date)
VALUES (@account_id, @balance, @created_date, @last_modified_date);

-- name: GetBalance :one
SELECT balance FROM accounts WHERE account_id = @account_id;

-- name: GetAccount :one
SELECT * FROM accounts WHERE account_id = @account_id;

-- name: GetAccounts :many
SELECT * FROM accounts where account_id = any(@account_ids::bigint[]);

-- name: DebitAccount :execresult
UPDATE accounts
SET
    balance = balance - @amount,
    last_modified_date = NOW()
WHERE
    account_id = @account_id AND balance >= @amount;

-- name: CreditAccount :execresult
UPDATE accounts
SET
    balance = balance + @amount,
    last_modified_date = NOW()
WHERE
    account_id = @account_id;
