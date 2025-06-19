package handlers

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"modforge.ai/ai"
	"modforge.ai/api/config"
	"modforge.ai/api/database"
	"modforge.ai/api/models"
	"modforge.ai/api/storage"
	"modforge.ai/mods"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Handlers contains all HTTP handlers
type Handlers struct {
	db       *database.DB
	cfg      *config.Config
	storage  *storage.Client
	aiClient *ai.Client
}

// New creates a new handlers instance
func New(db *database.DB, cfg *config.Config, storageClient *storage.Client, aiClient *ai.Client) *Handlers {
	return &Handlers{
		db:       db,
		cfg:      cfg,
		storage:  storageClient,
		aiClient: aiClient,
	}
}

// HealthCheck returns the health status of the API
func (h *Handlers) HealthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "ok",
		"service": "ModForge.ai API",
		"version": "1.0.0",
	})
}

// VerifyToken verifies Firebase auth token
func (h *Handlers) VerifyToken(c *fiber.Ctx) error {
	// TODO: Implement Firebase token verification
	return c.JSON(fiber.Map{"message": "Token verification endpoint"})
}

// GetUserProfile returns user profile information
func (h *Handlers) GetUserProfile(c *fiber.Ctx) error {
	userID := c.Params("id")
	if userID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "User ID is required"})
	}

	user, err := h.db.GetUserByID(userID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	return c.JSON(user)
}

// UpdateUserProfile updates user profile information
func (h *Handlers) UpdateUserProfile(c *fiber.Ctx) error {
	userID := c.Params("id")
	if userID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "User ID is required"})
	}

	var updateData models.User
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	updateData.ID = userID
	if err := h.db.UpdateUser(&updateData); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update user"})
	}

	return c.JSON(fiber.Map{"message": "User updated successfully"})
}

// UploadMod handles mod file upload and analysis
func (h *Handlers) UploadMod(c *fiber.Ctx) error {
	ctx := context.Background()

	// Get the uploaded file
	file, err := c.FormFile("mod_file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "No file uploaded"})
	}

	// Validate file type and size
	if !isValidModFile(file.Filename) {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid file type. Only .jar, .zip, .json files are allowed"})
	}

	if file.Size > 100*1024*1024 { // 100MB limit
		return c.Status(400).JSON(fiber.Map{"error": "File too large. Maximum size is 100MB"})
	}

	// Open the file
	src, err := file.Open()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to open file"})
	}
	defer src.Close()

	// Read file content
	content, err := io.ReadAll(src)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to read file"})
	}

	// Upload to storage
	fileURL, err := h.storage.UploadFile(ctx, content, file.Filename, file.Header.Get("Content-Type"))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("Failed to upload file: %v", err)})
	}

	// Detect mod type
	modType, err := mods.DetectModType(content, file.Filename)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": fmt.Sprintf("Failed to detect mod type: %v", err)})
	}

	// Create a new job
	job := &models.Job{
		ID:          uuid.New().String(),
		UserID:      c.Get("user_id", "anonymous"), // Get from auth middleware
		Status:      "pending",
		ModType:     modType,
		OriginalURL: fileURL,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	// Set required fields that can't be null
	filename := file.Filename
	fileSize := file.Size
	presetType := "default"

	job.OriginalFilename = &filename
	job.OriginalFileSize = &fileSize
	job.PresetType = &presetType

	if err := h.db.CreateJob(job); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create job"})
	}

	return c.JSON(fiber.Map{
		"job_id":   job.ID,
		"status":   job.Status,
		"mod_type": job.ModType,
		"message":  "File uploaded successfully",
	})
}

// GetJobStatus returns the status of a processing job
func (h *Handlers) GetJobStatus(c *fiber.Ctx) error {
	jobID := c.Params("id")
	if jobID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Job ID is required"})
	}

	job, err := h.db.GetJobByID(jobID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Job not found"})
	}

	return c.JSON(job)
}

// ProcessMod processes a mod using AI
func (h *Handlers) ProcessMod(c *fiber.Ctx) error {
	ctx := context.Background()
	jobID := c.Params("id")
	if jobID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Job ID is required"})
	}

	// Get processing parameters
	var params struct {
		PresetID    string `json:"preset_id"`
		Prompt      string `json:"prompt"`
		ModelConfig string `json:"model_config"`
	}
	if err := c.BodyParser(&params); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Get the job
	job, err := h.db.GetJobByID(jobID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Job not found"})
	}

	// Update job status to processing
	job.Status = "processing"
	job.UpdatedAt = time.Now()
	if err := h.db.UpdateJob(job); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update job status"})
	}

	// Process in background (for now, we'll do it synchronously)
	go func() {
		h.processModInBackground(ctx, job, params.PresetID, params.Prompt)
	}()

	return c.JSON(fiber.Map{
		"message": "Processing started",
		"job_id":  job.ID,
		"status":  job.Status,
	})
}

// processModInBackground handles the actual mod processing
func (h *Handlers) processModInBackground(ctx context.Context, job *models.Job, presetID, prompt string) {
	// Download original file
	content, err := h.storage.DownloadFile(ctx, job.OriginalURL)
	if err != nil {
		h.updateJobStatus(job.ID, "failed", fmt.Sprintf("Failed to download file: %v", err))
		return
	}

	// Check if we should use mock AI processing (for testing when OpenAI quota exceeded)
	useMockAI := true // Set to false when you have OpenAI quota

	var processedResponse *ai.ProcessModResponse

	if useMockAI {
		// Mock AI response for testing
		processedResponse = &ai.ProcessModResponse{
			ProcessedContent: h.generateMockEnhancedContent(string(content), prompt),
			Changelog:        "Mock AI Enhancement: " + prompt,
			TokensUsed:       100, // Mock token usage
		}
	} else {
		// Use real AI to process the mod
		realResponse, err := h.aiClient.ProcessMod(ctx, ai.ProcessModRequest{
			Content:        string(content),
			PromptTemplate: prompt,
			GameType:       job.ModType,
			Variables:      map[string]string{},
		})
		if err != nil {
			h.updateJobStatus(job.ID, "failed", fmt.Sprintf("AI processing failed: %v", err))
			return
		}
		processedResponse = realResponse
	}

	// Upload processed file
	filename := fmt.Sprintf("processed_%s_%s", job.ID, filepath.Base(job.OriginalURL))
	processedURL, err := h.storage.UploadFile(ctx, []byte(processedResponse.ProcessedContent), filename, "application/octet-stream")
	if err != nil {
		h.updateJobStatus(job.ID, "failed", fmt.Sprintf("Failed to upload processed file: %v", err))
		return
	}

	// Update job with results
	job.Status = "completed"
	job.ProcessedURL = &processedURL
	job.TokensUsed = &processedResponse.TokensUsed
	creditsUsed := 2 // Mock credits used
	job.CreditsUsed = &creditsUsed
	job.UpdatedAt = time.Now()
	if err := h.db.UpdateJob(job); err != nil {
		h.updateJobStatus(job.ID, "failed", fmt.Sprintf("Failed to update job: %v", err))
		return
	}
}

// generateMockEnhancedContent creates a mock enhanced version of the content
func (h *Handlers) generateMockEnhancedContent(originalContent, prompt string) string {
	// Create a simple mock enhancement based on the prompt
	enhancedContent := originalContent

	// Add some mock enhancements
	enhancedContent = strings.ReplaceAll(enhancedContent, "minecraft:diamond", "minecraft:netherite_ingot")
	enhancedContent = strings.ReplaceAll(enhancedContent, "enhanced_diamond_sword", "legendary_netherite_sword")
	enhancedContent = strings.ReplaceAll(enhancedContent, "\"count\": 1", "\"count\": 1,\n      \"components\": {\n        \"minecraft:enchantments\": {\n          \"minecraft:sharpness\": 3,\n          \"minecraft:unbreaking\": 2\n        }\n      }")

	// Add a comment about the enhancement
	enhancement := fmt.Sprintf("\n// Enhanced by ModForge.ai with prompt: %s\n// - Upgraded materials from diamond to netherite\n// - Added enchantments for better gameplay\n// - Balanced recipe cost\n", prompt)

	return enhancement + enhancedContent
}

// updateJobStatus is a helper to update job status
func (h *Handlers) updateJobStatus(jobID, status, errorMsg string) {
	job, err := h.db.GetJobByID(jobID)
	if err != nil {
		return
	}
	job.Status = status
	if errorMsg != "" {
		job.ErrorMessage = &errorMsg
	}
	job.UpdatedAt = time.Now()
	h.db.UpdateJob(job)
}

// DownloadMod handles mod download
func (h *Handlers) DownloadMod(c *fiber.Ctx) error {
	ctx := context.Background()
	jobID := c.Params("id")
	if jobID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Job ID is required"})
	}

	job, err := h.db.GetJobByID(jobID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Job not found"})
	}

	if job.Status != "completed" || job.ProcessedURL == nil {
		return c.Status(400).JSON(fiber.Map{"error": "Job not completed or no processed file available"})
	}

	// Get presigned URL for download
	downloadURL, err := h.storage.GetPresignedURL(ctx, *job.ProcessedURL, 1*time.Hour)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate download URL"})
	}

	return c.JSON(fiber.Map{
		"download_url": downloadURL,
		"expires_in":   3600, // 1 hour
	})
}

// GetUserJobs returns all jobs for a user
func (h *Handlers) GetUserJobs(c *fiber.Ctx) error {
	userID := c.Get("user_id", "anonymous")

	// Parse query parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	status := c.Query("status")

	jobs, err := h.db.GetUserJobs(userID, page, limit, status)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch jobs"})
	}

	return c.JSON(fiber.Map{
		"jobs":  jobs,
		"page":  page,
		"limit": limit,
	})
}

// GetPresets returns all available presets
func (h *Handlers) GetPresets(c *fiber.Ctx) error {
	presets, err := h.db.GetPresets()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch presets"})
	}

	return c.JSON(fiber.Map{"presets": presets})
}

// GetPresetsByType returns presets filtered by mod type
func (h *Handlers) GetPresetsByType(c *fiber.Ctx) error {
	modType := c.Params("type")
	if modType == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Mod type is required"})
	}

	presets, err := h.db.GetPresetsByType(modType)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch presets"})
	}

	return c.JSON(fiber.Map{"presets": presets})
}

// GetCredits returns user's credit balance
func (h *Handlers) GetCredits(c *fiber.Ctx) error {
	userID := c.Get("user_id", "anonymous")

	user, err := h.db.GetUserByID(userID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	return c.JSON(fiber.Map{
		"credits": user.Credits,
		"user_id": user.ID,
	})
}

// PurchaseCredits handles credit purchases
func (h *Handlers) PurchaseCredits(c *fiber.Ctx) error {
	userID := c.Get("user_id", "anonymous")

	var purchase struct {
		Amount int    `json:"amount"`
		Token  string `json:"payment_token"`
	}
	if err := c.BodyParser(&purchase); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if purchase.Amount <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid credit amount"})
	}

	// TODO: Implement actual payment processing
	// For now, just add credits directly (for testing)
	user, err := h.db.GetUserByID(userID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	user.Credits += purchase.Amount
	if err := h.db.UpdateUser(user); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update credits"})
	}

	return c.JSON(fiber.Map{
		"message": "Credits purchased successfully",
		"credits": user.Credits,
	})
}

// isValidModFile checks if the uploaded file is a valid mod file
func isValidModFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	validExtensions := []string{".jar", ".zip", ".json", ".mcmeta"}

	for _, validExt := range validExtensions {
		if ext == validExt {
			return true
		}
	}
	return false
}
