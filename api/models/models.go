package models

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID                   string    `json:"id" db:"id"`
	Email                string    `json:"email" db:"email"`
	Password             string    `json:"-" db:"password_hash"` // Hidden from JSON response
	FirebaseUID          *string   `json:"firebase_uid,omitempty" db:"firebase_uid"`
	DisplayName          string    `json:"display_name" db:"display_name"`
	Credits              int       `json:"credits" db:"credits"`
	Plan                 string    `json:"plan" db:"plan"`
	MonthlyJobsUsed      int       `json:"monthly_jobs_used" db:"monthly_jobs_used"`
	MonthlyJobsResetDate time.Time `json:"monthly_jobs_reset_date" db:"monthly_jobs_reset_date"`
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time `json:"updated_at" db:"updated_at"`
}

// Job represents a mod processing job (alias for ModJob for handler compatibility)
type Job struct {
	ID               string    `json:"id" db:"id"`
	UserID           string    `json:"user_id" db:"user_id"`
	Status           string    `json:"status" db:"status"`
	ModType          string    `json:"mod_type" db:"game_type"` // Map to game_type in DB
	OriginalFilename *string   `json:"original_filename,omitempty" db:"original_filename"`
	OriginalFileSize *int64    `json:"original_file_size,omitempty" db:"original_file_size"`
	OriginalURL      string    `json:"original_url" db:"original_file_url"`
	ProcessedURL     *string   `json:"processed_url,omitempty" db:"processed_file_url"`
	PresetType       *string   `json:"preset_type,omitempty" db:"preset_type"`
	AIPrompt         *string   `json:"ai_prompt,omitempty" db:"ai_prompt"`
	AIResponse       *string   `json:"ai_response,omitempty" db:"ai_response"`
	Changelog        *string   `json:"changelog,omitempty" db:"changelog"`
	TokensUsed       *int      `json:"tokens_used,omitempty" db:"tokens_used"`
	CreditsUsed      *int      `json:"credits_used,omitempty" db:"credits_used"`
	ErrorMessage     *string   `json:"error_message,omitempty" db:"error_message"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

// ModJob represents a mod processing job (legacy name, keeping for compatibility)
type ModJob = Job

// ModPreset represents a predefined AI transformation preset
type ModPreset struct {
	ID             string    `json:"id" db:"id"`
	Name           string    `json:"name" db:"name"`
	Description    string    `json:"description" db:"description"`
	GameType       string    `json:"game_type" db:"game_type"`
	PromptTemplate string    `json:"prompt_template" db:"prompt_template"`
	CreditCost     int       `json:"credit_cost" db:"credit_cost"`
	IsActive       bool      `json:"is_active" db:"is_active"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

// UserSession represents a user authentication session
type UserSession struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	Token     string    `json:"token" db:"token"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Job status constants
const (
	JobStatusPending    = "pending"
	JobStatusProcessing = "processing"
	JobStatusCompleted  = "completed"
	JobStatusFailed     = "failed"
)

// Game type constants
const (
	GameTypeMinecraft = "minecraft"
	GameTypeSkyrim    = "skyrim"
	GameTypeLua       = "lua"
)

// Plan constants
const (
	PlanFree = "free"
	PlanPro  = "pro"
)
