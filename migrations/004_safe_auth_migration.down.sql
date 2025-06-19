-- Rollback safe auth migration
ALTER TABLE users DROP COLUMN IF EXISTS password_hash;

-- Restore firebase_uid as NOT NULL (this might fail if there are NULL values)
-- ALTER TABLE users ALTER COLUMN firebase_uid SET NOT NULL;

-- Drop user_sessions table
DROP TABLE IF EXISTS user_sessions;
