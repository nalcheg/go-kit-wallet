CREATE TABLE payments
(
    id UUID NOT NULL
        CONSTRAINT payments_pk
        PRIMARY KEY,
    amount BIGINT NOT NULL,
    account VARCHAR(255) NOT NULL
        CONSTRAINT payments_accounts_id_fk
            REFERENCES accounts,
    from_account VARCHAR(255) NULL
        CONSTRAINT payments_accounts_id_fk_2
            REFERENCES accounts,
    to_account VARCHAR(255) NULL
        CONSTRAINT payments_accounts_id_fk_3
            REFERENCES accounts
);

CREATE INDEX payments_account_index
    ON payments (account);

CREATE INDEX payments_from_account_index
    ON payments (from_account);

CREATE INDEX payments_to_account_index
    ON payments (to_account);

---- create above / drop below ----

DROP TABLE payments;
