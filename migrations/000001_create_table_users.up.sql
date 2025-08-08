CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(255) PRIMARY KEY,
    username VARCHAR(100) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255), -- NULL if register or login by OAuth
    registration_method VARCHAR(20) NOT NULL, -- manual, oauth, admin
    oauth_provider VARCHAR(50), -- NULL if not OAuth
    status VARCHAR(50) NOT NULL DEFAULT 'pending', -- pending, active, suspend, banned, locked, inactive
    email_verified BOOLEAN NOT NULL DEFAULT FALSE,
    email_verified_at TIMESTAMP WITH TIME ZONE,
    last_action_by VARCHAR(255),
    issued_reason VARCHAR(255),
    issued_at TIMESTAMP WITH TIME ZONE,
    current_login_at TIMESTAMP WITH TIME ZONE,
    current_login_ip VARCHAR(45), -- IPv4 or IPv6
    failed_login_attempts INT DEFAULT 0,
    locked_until TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes untuk performance
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_status ON users(status);

-- Trigger untuk auto-update updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = NOW();
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();