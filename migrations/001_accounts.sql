CREATE TABLE accounts
(
    id VARCHAR(255) NOT NULL
        CONSTRAINT accounts_pk
        PRIMARY KEY,
    balance BIGINT NOT NULL,
    currency VARCHAR(255) NOT NULL
);

---- create above / drop below ----

DROP TABLE accounts;
