package storage

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

// Client wraps the storage client (simplified for MVP)
type Client struct {
	bucketName string
	accountID  string
	apiToken   string
}

// Config holds storage configuration
type Config struct {
	AccountID  string
	APIToken   string
	BucketName string
	Region     string
}

// NewClient creates a new storage client for Cloudflare R2
func NewClient(cfg Config) (*Client, error) {
	return &Client{
		bucketName: cfg.BucketName,
		accountID:  cfg.AccountID,
		apiToken:   cfg.APIToken,
	}, nil
}

// UploadFile uploads a file to storage (mock implementation for MVP)
func (c *Client) UploadFile(ctx context.Context, content []byte, filename string, contentType string) (string, error) {
	// For MVP, we'll save files locally and return a mock URL
	// In production, implement actual R2 upload using their API

	// Create uploads directory if it doesn't exist
	uploadsDir := "./uploads"
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create uploads directory: %w", err)
	}

	// Generate unique filename
	fileID := uuid.New().String()
	localPath := filepath.Join(uploadsDir, fileID+"_"+filename)

	// Save file locally
	if err := os.WriteFile(localPath, content, 0644); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	// Return mock URL (in production, return actual R2 URL)
	return fmt.Sprintf("http://localhost:8080/uploads/%s_%s", fileID, filename), nil
}

// DownloadFile downloads a file from storage (mock implementation)
func (c *Client) DownloadFile(ctx context.Context, url string) ([]byte, error) {
	// For MVP, read from local storage
	// Extract filename from URL and read local file
	filename := filepath.Base(url)
	localPath := filepath.Join("./uploads", filename)

	content, err := os.ReadFile(localPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return content, nil
}

// DeleteFile deletes a file from storage (mock implementation)
func (c *Client) DeleteFile(ctx context.Context, url string) error {
	// For MVP, delete from local storage
	filename := filepath.Base(url)
	localPath := filepath.Join("./uploads", filename)

	if err := os.Remove(localPath); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

// GetPresignedURL generates a presigned URL for direct file access (mock implementation)
func (c *Client) GetPresignedURL(ctx context.Context, url string, expiration time.Duration) (string, error) {
	// For MVP, return the same URL (since we're serving locally)
	// In production, generate actual presigned R2 URL
	return url, nil
}
