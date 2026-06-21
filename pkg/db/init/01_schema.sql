CREATE TABLE IF NOT EXISTS users (
    id BIGINT PRIMARY KEY,
    full_name VARCHAR,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_users_full_name ON users(full_name);

CREATE TABLE IF NOT EXISTS wallets (
    id BIGINT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    balance BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_wallets_user_id FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT balance_check CHECK (balance >= 0)
);

CREATE INDEX IF NOT EXISTS idx_wallets_user_id ON wallets(user_id);

CREATE TABLE IF NOT EXISTS transaction_logs (
    id BIGINT PRIMARY KEY,
    user_id BIGINT NOT NULL, -- not fk since this is table for logging that persist even if user deleted 
    user_full_name VARCHAR NOT NULL, 
    action TEXT,
    amount BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_transaction_logs_action ON transaction_logs(action);
CREATE INDEX IF NOT EXISTS idx_transaction_logs_user_id ON transaction_logs(user_id);