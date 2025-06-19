package database

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"modforge.ai/api/models"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"           // PostgreSQL driver
	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// DB wraps the database connection
type DB struct {
	*sql.DB
}

// Initialize creates a new database connection
func Initialize(databaseURL string) (*DB, error) {
	var sqlDB *sql.DB
	var err error

	if strings.HasPrefix(databaseURL, "postgresql://") || strings.HasPrefix(databaseURL, "postgres://") {
		// PostgreSQL
		sqlDB, err = sql.Open("postgres", databaseURL)
	} else {
		// SQLite (default)
		dbPath := strings.TrimPrefix(databaseURL, "sqlite://")
		sqlDB, err = sql.Open("sqlite3", dbPath)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{sqlDB}, nil
}

// RunMigrations runs database migrations
func RunMigrations(databaseURL string) error {
	// For production with existing database, apply manual schema updates
	if strings.HasPrefix(databaseURL, "postgresql://") || strings.HasPrefix(databaseURL, "postgres://") {
		log.Println("Production PostgreSQL detected - applying manual schema updates...")

		db, err := sql.Open("postgres", databaseURL)
		if err != nil {
			return fmt.Errorf("failed to open database: %w", err)
		}
		defer db.Close()

		// Check if users table exists
		var exists bool
		err = db.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'users')").Scan(&exists)
		if err != nil {
			return fmt.Errorf("failed to check for users table: %w", err)
		}

		if exists {
			log.Println("Applying auth schema updates to existing database...")

			// Add password_hash column if it doesn't exist
			_, err = db.Exec(`
				DO $$ 
				BEGIN
					IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
								   WHERE table_name = 'users' AND column_name = 'password_hash') THEN
						ALTER TABLE users ADD COLUMN password_hash TEXT;
					END IF;
				END $$;
			`)
			if err != nil {
				log.Printf("Warning: Failed to add password_hash column: %v", err)
			}

			// Make firebase_uid nullable
			_, err = db.Exec(`
				DO $$ 
				BEGIN
					ALTER TABLE users ALTER COLUMN firebase_uid DROP NOT NULL;
				EXCEPTION
					WHEN OTHERS THEN NULL;
				END $$;
			`)
			if err != nil {
				log.Printf("Warning: Failed to make firebase_uid nullable: %v", err)
			}

			// Drop and recreate user_sessions table
			_, err = db.Exec(`DROP TABLE IF EXISTS user_sessions;`)
			if err != nil {
				log.Printf("Warning: Failed to drop user_sessions: %v", err)
			}

			_, err = db.Exec(`
				CREATE TABLE user_sessions (
					id TEXT PRIMARY KEY,
					user_id TEXT NOT NULL,
					token TEXT UNIQUE NOT NULL,
					expires_at TIMESTAMP NOT NULL,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
				);
			`)
			if err != nil {
				log.Printf("Error creating user_sessions: %v", err)
				return err
			}

			// Create indexes
			_, err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_user_sessions_token ON user_sessions(token);`)
			if err != nil {
				log.Printf("Warning: Failed to create token index: %v", err)
			}

			log.Println("Schema updates completed successfully")
			return nil
		}
	}

	// Fall back to normal migration process for new databases
	return runStandardMigrations(databaseURL)
}

func runStandardMigrations(databaseURL string) error {
	var db *sql.DB
	var err error
	var driverName string

	if strings.HasPrefix(databaseURL, "postgresql://") || strings.HasPrefix(databaseURL, "postgres://") {
		// PostgreSQL
		db, err = sql.Open("postgres", databaseURL)
		driverName = "postgres"
	} else {
		// SQLite (default)
		dbPath := strings.TrimPrefix(databaseURL, "sqlite://")
		db, err = sql.Open("sqlite3", dbPath)
		driverName = "sqlite3"
	}

	if err != nil {
		return fmt.Errorf("failed to open database for migrations: %w", err)
	}
	defer db.Close()

	// For PostgreSQL, try to clean up any existing migration state
	if driverName == "postgres" {
		log.Println("Checking PostgreSQL migration state...")

		// Check if schema_migrations table exists
		var exists bool
		err := db.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'schema_migrations')").Scan(&exists)
		if err != nil {
			log.Printf("Error checking for schema_migrations table: %v", err)
		}

		if exists {
			log.Println("Found existing schema_migrations table, checking state...")

			// Check if we have a dirty state or version conflicts
			var version uint
			var dirty bool
			err := db.QueryRow("SELECT version, dirty FROM schema_migrations ORDER BY version DESC LIMIT 1").Scan(&version, &dirty)
			if err == nil && !dirty {
				log.Printf("Clean migration state found at version %d, proceeding normally", version)
			} else {
				log.Println("Migration state issues found, cleaning up...")
				// Clean up the migration table to start fresh
				if _, err := db.Exec("DELETE FROM schema_migrations"); err != nil {
					log.Printf("Warning: Failed to clean schema_migrations: %v", err)
				}
			}
		}
	}

	// Create driver for migrations
	var driver database.Driver
	if driverName == "postgres" {
		driver, err = postgres.WithInstance(db, &postgres.Config{})
	} else {
		driver, err = sqlite3.WithInstance(db, &sqlite3.Config{})
	}
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	// Create migration instance
	m, err := migrate.NewWithDatabaseInstance(
		"file://./migrations",
		driverName,
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	// Run migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		// Check if it's a dirty database error or missing migration file
		if strings.Contains(err.Error(), "Dirty database") || strings.Contains(err.Error(), "no migration found for version 0") {
			log.Printf("Migration issue detected: %v", err)
			log.Println("Attempting to reset migration state...")

			// Try to force to version 1 (our actual first migration)
			if err := m.Force(1); err != nil {
				log.Printf("Failed to force to version 1, trying version 0: %v", err)
				// If that fails, try forcing to 0
				if err := m.Force(0); err != nil {
					return fmt.Errorf("failed to force clean migration state: %w", err)
				}
			}

			// Try running migrations again
			if err := m.Up(); err != nil && err != migrate.ErrNoChange {
				return fmt.Errorf("failed to run migrations after force reset: %w", err)
			}
		} else {
			return fmt.Errorf("failed to run migrations: %w", err)
		}
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// GetUserByID retrieves a user by ID
func (db *DB) GetUserByID(id string) (*models.User, error) {
	query := `
		SELECT id, email, password_hash, firebase_uid, display_name, credits, plan,
		       monthly_jobs_used, monthly_jobs_reset_date, created_at, updated_at
		FROM users WHERE id = $1
	`

	user := &models.User{}
	err := db.QueryRow(query, id).Scan(
		&user.ID, &user.Email, &user.Password, &user.FirebaseUID, &user.DisplayName,
		&user.Credits, &user.Plan, &user.MonthlyJobsUsed,
		&user.MonthlyJobsResetDate, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// UpdateUser updates a user record
func (db *DB) UpdateUser(user *models.User) error {
	query := `
		UPDATE users SET 
			email = $1, display_name = $2, credits = $3, plan = $4,
			monthly_jobs_used = $5, monthly_jobs_reset_date = $6, updated_at = $7
		WHERE id = $8
	`

	user.UpdatedAt = time.Now()

	_, err := db.Exec(query,
		user.Email, user.DisplayName, user.Credits, user.Plan,
		user.MonthlyJobsUsed, user.MonthlyJobsResetDate, user.UpdatedAt,
		user.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// CreateUser creates a new user record
func (db *DB) CreateUser(user *models.User) error {
	query := `
		INSERT INTO users (id, email, password_hash, firebase_uid, display_name, credits, plan, monthly_jobs_used, monthly_jobs_reset_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $2)
	`

	_, err := db.Exec(query,
		user.ID, user.Email, user.Password, user.FirebaseUID, user.DisplayName,
		user.Credits, user.Plan, user.MonthlyJobsUsed, user.MonthlyJobsResetDate,
		user.CreatedAt, user.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetUserByEmail retrieves a user by email
func (db *DB) GetUserByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, email, password_hash, firebase_uid, display_name, credits, plan,
		       monthly_jobs_used, monthly_jobs_reset_date, created_at, updated_at
		FROM users WHERE email = $1
	`

	user := &models.User{}
	err := db.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.Password, &user.FirebaseUID, &user.DisplayName,
		&user.Credits, &user.Plan, &user.MonthlyJobsUsed,
		&user.MonthlyJobsResetDate, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// CreateJob creates a new job record
func (db *DB) CreateJob(job *models.Job) error {
	query := `
		INSERT INTO mod_jobs (id, user_id, status, game_type, original_filename, original_file_size, original_file_url, preset_type, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := db.Exec(query,
		job.ID, job.UserID, job.Status, job.ModType,
		job.OriginalFilename, job.OriginalFileSize, job.OriginalURL, job.PresetType,
		job.CreatedAt, job.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create job: %w", err)
	}

	return nil
}

// GetJobByID retrieves a job by ID
func (db *DB) GetJobByID(id string) (*models.Job, error) {
	query := `
		SELECT id, user_id, status, game_type, original_filename, original_file_size,
		       original_file_url, processed_file_url, preset_type, ai_prompt,
		       ai_response, changelog, tokens_used, credits_used, error_message,
		       created_at, updated_at
		FROM mod_jobs WHERE id = $1
	`

	job := &models.Job{}
	err := db.QueryRow(query, id).Scan(
		&job.ID, &job.UserID, &job.Status, &job.ModType,
		&job.OriginalFilename, &job.OriginalFileSize, &job.OriginalURL,
		&job.ProcessedURL, &job.PresetType, &job.AIPrompt,
		&job.AIResponse, &job.Changelog, &job.TokensUsed,
		&job.CreditsUsed, &job.ErrorMessage, &job.CreatedAt, &job.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("job not found")
		}
		return nil, fmt.Errorf("failed to get job: %w", err)
	}

	return job, nil
}

// UpdateJob updates a job record
func (db *DB) UpdateJob(job *models.Job) error {
	query := `
		UPDATE mod_jobs SET 
			status = $1, processed_file_url = $2, preset_type = $3, ai_prompt = $4,
			ai_response = $5, changelog = $6, tokens_used = $7, credits_used = $8,
			error_message = $9, updated_at = $10
		WHERE id = $11
	`

	job.UpdatedAt = time.Now()

	_, err := db.Exec(query,
		job.Status, job.ProcessedURL, job.PresetType, job.AIPrompt,
		job.AIResponse, job.Changelog, job.TokensUsed, job.CreditsUsed,
		job.ErrorMessage, job.UpdatedAt, job.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update job: %w", err)
	}

	return nil
}

// CreateSession creates a new user session
func (db *DB) CreateSession(session *models.UserSession) error {
	query := `
		INSERT INTO user_sessions (id, user_id, token, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := db.Exec(query,
		session.ID, session.UserID, session.Token, session.ExpiresAt, session.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	return nil
}

// GetSessionByToken retrieves a session by token
func (db *DB) GetSessionByToken(token string) (*models.UserSession, error) {
	query := `
		SELECT id, user_id, token, expires_at, created_at
		FROM user_sessions WHERE token = $1
	`

	session := &models.UserSession{}
	err := db.QueryRow(query, token).Scan(
		&session.ID, &session.UserID, &session.Token, &session.ExpiresAt, &session.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("session not found")
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	return session, nil
}

// DeleteSession deletes a user session
func (db *DB) DeleteSession(token string) error {
	query := `DELETE FROM user_sessions WHERE token = $1`

	_, err := db.Exec(query, token)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	return nil
}

// GetUserJobs retrieves jobs for a user with pagination
func (db *DB) GetUserJobs(userID string, page, limit int, status string) ([]*models.Job, error) {
	offset := (page - 1) * limit

	query := `
		SELECT id, user_id, status, game_type, original_filename, original_file_size,
		       original_file_url, processed_file_url, preset_type, ai_prompt,
		       ai_response, changelog, tokens_used, credits_used, error_message,
		       created_at, updated_at
		FROM mod_jobs 
		WHERE user_id = $1
	`
	args := []interface{}{userID}

	if status != "" {
		query += " AND status = $1"
		args = append(args, status)
	}

	query += " ORDER BY created_at DESC LIMIT $1 OFFSET $1"
	args = append(args, limit, offset)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query jobs: %w", err)
	}
	defer rows.Close()

	var jobs []*models.Job
	for rows.Next() {
		job := &models.Job{}
		err := rows.Scan(
			&job.ID, &job.UserID, &job.Status, &job.ModType,
			&job.OriginalFilename, &job.OriginalFileSize, &job.OriginalURL,
			&job.ProcessedURL, &job.PresetType, &job.AIPrompt,
			&job.AIResponse, &job.Changelog, &job.TokensUsed,
			&job.CreditsUsed, &job.ErrorMessage, &job.CreatedAt, &job.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan job: %w", err)
		}
		jobs = append(jobs, job)
	}

	return jobs, nil
}

// GetPresets retrieves all active presets
func (db *DB) GetPresets() ([]*models.ModPreset, error) {
	query := `
		SELECT id, name, description, game_type, prompt_template, credit_cost, is_active, created_at
		FROM mod_presets WHERE is_active = true
		ORDER BY name
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query presets: %w", err)
	}
	defer rows.Close()

	var presets []*models.ModPreset
	for rows.Next() {
		preset := &models.ModPreset{}
		err := rows.Scan(
			&preset.ID, &preset.Name, &preset.Description, &preset.GameType,
			&preset.PromptTemplate, &preset.CreditCost, &preset.IsActive, &preset.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan preset: %w", err)
		}
		presets = append(presets, preset)
	}

	return presets, nil
}

// GetPresetsByType retrieves presets filtered by game type
func (db *DB) GetPresetsByType(gameType string) ([]*models.ModPreset, error) {
	query := `
		SELECT id, name, description, game_type, prompt_template, credit_cost, is_active, created_at
		FROM mod_presets WHERE game_type = $1 AND is_active = true
		ORDER BY name
	`

	rows, err := db.Query(query, gameType)
	if err != nil {
		return nil, fmt.Errorf("failed to query presets: %w", err)
	}
	defer rows.Close()

	var presets []*models.ModPreset
	for rows.Next() {
		preset := &models.ModPreset{}
		err := rows.Scan(
			&preset.ID, &preset.Name, &preset.Description, &preset.GameType,
			&preset.PromptTemplate, &preset.CreditCost, &preset.IsActive, &preset.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan preset: %w", err)
		}
		presets = append(presets, preset)
	}

	return presets, nil
}
