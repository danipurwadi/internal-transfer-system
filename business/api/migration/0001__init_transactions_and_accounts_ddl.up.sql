CREATE TABLE
    IF NOT EXISTS accounts (
        account_id BIGINT NOT NULL,
        status VARCHAR(100) NOT NULL,
        created_date TIMESTAMPTZ NOT NULL,
        last_modified_date TIMESTAMPTZ NOT NULL
    );

CREATE TABLE
    if not EXISTS transactions (
        account_id BIGINT NOT NULL,
        amount NUMERIC(19, 4) NOT NULL,
        created_date TIMESTAMPTZ NOT NULL
    );