CREATE TABLE
    IF NOT EXISTS accounts (
        account_id BIGINT PRIMARY KEY,
        created_date TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        last_modified_date TIMESTAMPTZ NOT NULL DEFAULT NOW ()
    );

CREATE TABLE
    if not EXISTS transactions (
        account_id BIGINT NOT NULL REFERENCES accounts (account_id) ON DELETE RESTRICT,
        amount NUMERIC(19, 5) NOT NULL,
        created_date TIMESTAMPTZ NOT NULL DEFAULT NOW ()
    );