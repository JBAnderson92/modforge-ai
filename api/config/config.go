package config

import (
	"os"
	"strconv"
)

// Config holds all configuration for the application
type Config struct {
	DatabaseURL    string
	RedisURL       string
	OpenAIAPIKey   string
	FirebaseConfig string
	CloudflareR2   CloudflareR2Config
	VirusTotalKey  string
	AllowedOrigins string
	RateLimit      RateLimitConfig
}

// CloudflareR2Config holds Cloudflare R2 storage configuration
type CloudflareR2Config struct {
	AccountID  string
	APIToken   string
	BucketName string
	Region     string
}

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	RequestsPerMinute int
	FreeMonthlyJobs   int
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		DatabaseURL:    getEnv("DATABASE_URL", "sqlite://./modforge.db"),
		RedisURL:       getEnv("REDIS_URL", "redis://localhost:6379"),
		OpenAIAPIKey:   getEnv("OPENAI_API_KEY", ""),
		FirebaseConfig: getEnv("FIREBASE_CONFIG", ""),
		CloudflareR2: CloudflareR2Config{
			AccountID:  getEnv("CLOUDFLARE_R2_ACCOUNT_ID", ""),
			APIToken:   getEnv("CLOUDFLARE_R2_API_TOKEN", ""),
			BucketName: getEnv("CLOUDFLARE_R2_BUCKET_NAME", "modforge-files"),
			Region:     getEnv("CLOUDFLARE_R2_REGION", "auto"),
		},
		VirusTotalKey:  getEnv("VIRUSTOTAL_API_KEY", ""),
		AllowedOrigins: getEnv("ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:5173"),
		RateLimit: RateLimitConfig{
			RequestsPerMinute: getEnvAsInt("RATE_LIMIT_RPM", 5),
			FreeMonthlyJobs:   getEnvAsInt("FREE_MONTHLY_JOBS", 3),
		},
	}
}

// getEnv gets an environment variable with a fallback value
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// getEnvAsInt gets an environment variable as an integer with a fallback value
func getEnvAsInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return fallback
}
