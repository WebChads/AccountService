-- Up migration: creates the accounts table
CREATE TABLE accounts (
    id SERIAL PRIMARY KEY,
    firstname VARCHAR(255) NOT NULL,
    surname VARCHAR(255) NOT NULL,
    patronymic VARCHAR(255),
    gender VARCHAR(1),
    birthdate DATE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create an index on surname for better query performance
CREATE INDEX idx_accounts_surname ON accounts(surname);