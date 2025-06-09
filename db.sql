create table accounts (
    id primary key,
    balance NUMERIC(20, 5) NOT NULL DEFAULT 0.00000,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE transactions (
    id serial primary key,
    source_id integer not null references accounts(id),
    destination_id integer not null references accounts(id),
    amount NUMERIC(20, 5) NOT NULL DEFAULT 0.00000,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
