-- +migrate Up
-- Create wallet table
CREATE TABLE IF NOT EXISTS wallets (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    balance DECIMAL(15,2) NOT NULL DEFAULT 0,
    currency VARCHAR(3) NOT NULL DEFAULT 'IDR',
    is_active BOOLEAN NOT NULL DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    UNIQUE KEY unique_user_wallet (user_id),
    KEY idx_user_id (user_id),
    KEY idx_deleted_at (deleted_at),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create wallet transactions table
CREATE TABLE IF NOT EXISTS wallet_transactions (
    id VARCHAR(36) PRIMARY KEY,
    wallet_id VARCHAR(36) NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    type ENUM('WITHDRAWAL', 'DEPOSIT') NOT NULL,
    status ENUM('PENDING', 'COMPLETED', 'FAILED') NOT NULL DEFAULT 'PENDING',
    description TEXT,
    metadata JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    KEY idx_wallet_id (wallet_id),
    KEY idx_type (type),
    KEY idx_status (status),
    KEY idx_created_at (created_at),
    KEY idx_deleted_at (deleted_at),
    FOREIGN KEY (wallet_id) REFERENCES wallets(id) ON DELETE CASCADE
);

-- Create index for transaction history queries
CREATE INDEX idx_wallet_created_at ON wallet_transactions(wallet_id, created_at DESC);

-- +migrate Down
DROP TABLE IF EXISTS wallet_transactions;
DROP TABLE IF EXISTS wallets;
