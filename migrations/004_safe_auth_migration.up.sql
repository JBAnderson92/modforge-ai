-- Safe migration that adds auth fields only if they don't exist
DO $$ 
BEGIN
    -- Add password_hash column if it doesn't exist
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_name = 'users' AND column_name = 'password_hash') THEN
        ALTER TABLE users ADD COLUMN password_hash TEXT;
    END IF;
    
    -- Make firebase_uid nullable
    BEGIN
        ALTER TABLE users ALTER COLUMN firebase_uid DROP NOT NULL;
    EXCEPTION
        WHEN OTHERS THEN NULL; -- Ignore if already nullable
    END;
END $$;

-- Recreate user_sessions table with new structure
DROP TABLE IF EXISTS user_sessions;

CREATE TABLE user_sessions (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    token TEXT UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_user_sessions_token ON user_sessions(token);
CREATE INDEX IF NOT EXISTS idx_user_sessions_user_id ON user_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
