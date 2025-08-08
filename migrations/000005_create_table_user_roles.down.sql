DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP TRIGGER IF EXISTS update_user_profiles_updated_at ON user_profiles;
DROP TRIGGER IF EXISTS update_email_verification_tokens_updated_at ON email_verification_tokens;
DROP TRIGGER IF EXISTS update_roles_updated_at ON roles;
DROP TRIGGER IF EXISTS update_user_roles_updated_at ON user_roles;

DROP FUNCTION IF EXISTS update_updated_at_column();

DROP INDEX IF EXISTS idx_user_roles_user_id;
DROP INDEX IF EXISTS idx_user_roles_role_id;
DROP INDEX IF EXISTS idx_user_roles_created_at;

ALTER TABLE user_roles DROP CONSTRAINT IF EXISTS fk_user_roles_role;
ALTER TABLE user_roles DROP CONSTRAINT IF EXISTS fk_user_roles_user;

DROP TABLE IF EXISTS user_roles;