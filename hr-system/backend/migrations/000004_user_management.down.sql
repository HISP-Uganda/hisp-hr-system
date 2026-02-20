DROP INDEX IF EXISTS idx_users_username;

ALTER TABLE users
    DROP COLUMN IF EXISTS last_login_at;
