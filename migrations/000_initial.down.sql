-- Drop all tables in reverse order to handle foreign key dependencies

-- Drop indexes first
DROP INDEX IF EXISTS idx_mod_jobs_status;
DROP INDEX IF EXISTS idx_mod_jobs_user_id;
DROP INDEX IF EXISTS idx_mod_jobs_created_at;
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_created_at;

-- Drop tables in reverse order
DROP TABLE IF EXISTS user_sessions;
DROP TABLE IF EXISTS mod_presets;
DROP TABLE IF EXISTS mod_jobs;
DROP TABLE IF EXISTS users;
