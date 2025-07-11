CREATE TABLE
    if not EXISTS transactions (
        account_id BIGINT NOT NULL,
        amount NUMERIC(19, 4) NOT NULL,
        created_date TIMESTAMPTZ NOT NULL
    );