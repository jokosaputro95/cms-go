DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP TRIGGER IF EXISTS update_user_profiles_updated_at ON user_profiles;
DROP TRIGGER IF EXISTS update_email_verification_tokens_updated_at ON email_verification_tokens;
DROP TRIGGER IF EXISTS update_roles_updated_at ON roles;
DROP TRIGGER IF EXISTS update_user_roles_updated_at ON user_roles;

DROP FUNCTION IF EXISTS update_updated_at_column();

DROP INDEX IF EXISTS idx_revoked_tokens_expires_at;
DROP INDEX IF EXISTS idx_revoked_tokens_revoked_at;
DROP INDEX IF EXISTS idx_revoked_tokens_token;


DROP TABLE IF EXISTS revoked_tokens;