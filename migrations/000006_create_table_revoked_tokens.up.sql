CREATE TABLE revoked_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    token TEXT NOT NULL UNIQUE,
    revoked_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL
);

-- Indexes untuk performance
CREATE INDEX IF NOT EXISTS idx_revoked_tokens_token ON revoked_tokens(token);
CREATE INDEX IF NOT EXISTS idx_revoked_tokens_revoked_at ON revoked_tokens(revoked_at);
CREATE INDEX IF NOT EXISTS idx_revoked_tokens_expires_at ON revoked_tokens(expires_at);