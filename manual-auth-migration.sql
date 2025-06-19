-- Manual migration to add auth fields to existing production database
-- Run this manually in Railway PostgreSQL console

-- Add password_hash column to users table (if it doesn't exist)
DO $$ 
BEGIN
    BEGIN
        ALTER TABLE users ADD COLUMN password_hash TEXT;
    EXCEPTION
        WHEN duplicate_column THEN NULL;
    END;
END $$;

-- Make firebase_uid nullable (for email/password auth)
ALTER TABLE users ALTER COLUMN firebase_uid DROP NOT NULL;

-- Drop and recreate user_sessions table with correct structure
DROP TABLE IF EXISTS user_sessions;

CREATE TABLE user_sessions (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    token TEXT UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_user_sessions_token ON user_sessions(token);
CREATE INDEX IF NOT EXISTS idx_user_sessions_user_id ON user_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
