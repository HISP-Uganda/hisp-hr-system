ALTER TABLE users
    ADD COLUMN IF NOT EXISTS last_login_at TIMESTAMPTZ;

CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
