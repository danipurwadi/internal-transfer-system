-- name: CreateAccount :exec
INSERT INTO accounts (account_id, status, created_date, last_modified_date)
VALUES (@account_id=account_id, @status=status, @created_date=created_date, @last_modified_date=last_modified_date); 