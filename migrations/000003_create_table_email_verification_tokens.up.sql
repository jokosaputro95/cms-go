CREATE TABLE IF NOT EXISTS email_verification_tokens (
    id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    token_type VARCHAR(255) NOT NULL,
    token VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    used_at TIMESTAMP WITH TIME ZONE,         -- NULL jika belum dipakai, terisi waktu saat dipakai
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,


    CONSTRAINT fk_email_verification_user
        FOREIGN KEY(user_id)
            REFERENCES users(id)
            ON DELETE CASCADE
);

-- Indexes untuk performance
CREATE INDEX IF NOT EXISTS idx_email_verfication_token ON email_verification_tokens(token);
CREATE INDEX IF NOT EXISTS idx_email_verfication_user_id ON email_verification_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_email_verfication_email ON email_verification_tokens(email);
CREATE INDEX IF NOT EXISTS idx_email_verfication_expires_at ON email_verification_tokens(expires_at);

-- Trigger untuk auto-update updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = NOW();
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_email_verification_tokens_updated_at
    BEFORE UPDATE ON email_verification_tokens
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();