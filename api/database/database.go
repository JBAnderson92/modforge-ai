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
	_ "github.com/lib/pq"          // PostgreSQL driver
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
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// GetUserByID retrieves a user by ID
func (db *DB) GetUserByID(id string) (*models.User, error) {
	query := `
		SELECT id, email, firebase_uid, display_name, credits, plan, 
		       monthly_jobs_used, monthly_jobs_reset_date, created_at, updated_at
		FROM users WHERE id = ?
	`

	user := &models.User{}
	err := db.QueryRow(query, id).Scan(
		&user.ID, &user.Email, &user.FirebaseUID, &user.DisplayName,
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
			email = ?, display_name = ?, credits = ?, plan = ?,
			monthly_jobs_used = ?, monthly_jobs_reset_date = ?, updated_at = ?
		WHERE id = ?
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

// CreateJob creates a new job record
func (db *DB) CreateJob(job *models.Job) error {
	query := `
		INSERT INTO mod_jobs (id, user_id, status, game_type, original_filename, original_file_size, original_file_url, preset_type, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
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
		FROM mod_jobs WHERE id = ?
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
			status = ?, processed_file_url = ?, preset_type = ?, ai_prompt = ?,
			ai_response = ?, changelog = ?, tokens_used = ?, credits_used = ?,
			error_message = ?, updated_at = ?
		WHERE id = ?
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

// GetUserJobs retrieves jobs for a user with pagination
func (db *DB) GetUserJobs(userID string, page, limit int, status string) ([]*models.Job, error) {
	offset := (page - 1) * limit

	query := `
		SELECT id, user_id, status, game_type, original_filename, original_file_size,
		       original_file_url, processed_file_url, preset_type, ai_prompt,
		       ai_response, changelog, tokens_used, credits_used, error_message,
		       created_at, updated_at
		FROM mod_jobs 
		WHERE user_id = ?
	`
	args := []interface{}{userID}

	if status != "" {
		query += " AND status = ?"
		args = append(args, status)
	}

	query += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
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
		FROM mod_presets WHERE game_type = ? AND is_active = true
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
