DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP TRIGGER IF EXISTS update_user_profiles_updated_at ON user_profiles;
DROP TRIGGER IF EXISTS update_email_verification_tokens_updated_at ON email_verification_tokens;
DROP FUNCTION IF EXISTS update_updated_at_column();

DROP INDEX IF EXISTS idx_email_verfication_expires_at;
DROP INDEX IF EXISTS idx_email_verfication_email;
DROP INDEX IF EXISTS idx_email_verfication_user_id;
DROP INDEX IF EXISTS idx_email_verfication_token;
ALTER TABLE email_verification_tokens DROP CONSTRAINT IF EXISTS fk_email_verification_user;
DROP TABLE IF EXISTS email_verification_tokens;