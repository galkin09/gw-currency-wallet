DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS wallets;

CREATE TABLE wallets
(
    id         INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    uuid       varchar NOT NULL UNIQUE,
    balanceRUB float DEFAULT 0 CHECK ( balanceRUB >= 0 ),
    balanceUSD float DEFAULT 0 CHECK ( balanceUSD >= 0 ),
    balanceEUR float DEFAULT 0 CHECK ( balanceEUR >= 0 )
);

CREATE TABLE users
(
    id serial PRIMARY KEY,
    username varchar NOT NULl UNIQUE ,
    password varchar NOT NULL,
    email varchar NOT NULL UNIQUE,
    wallet_id varchar,
    CONSTRAINT fk_users
        FOREIGN KEY(wallet_id) REFERENCES wallets(uuid)
);