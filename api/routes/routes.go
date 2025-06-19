package routes

import (
	"modforge.ai/ai"
	"modforge.ai/api/config"
	"modforge.ai/api/database"
	"modforge.ai/api/handlers"
	"modforge.ai/api/storage"

	"github.com/gofiber/fiber/v2"
)

// Setup initializes all routes for the application
func Setup(app *fiber.App, db *database.DB, cfg *config.Config, storageClient *storage.Client, aiClient *ai.Client) {
	// Initialize handlers
	h := handlers.New(db, cfg, storageClient, aiClient)

	// Static file serving for local uploads (for MVP)
	app.Static("/uploads", "./uploads")

	// API v1 routes
	v1 := app.Group("/api/v1")

	// Health check
	v1.Get("/health", h.HealthCheck)

	// Authentication routes
	auth := v1.Group("/auth")
	auth.Post("/register", h.Register)
	auth.Post("/login", h.Login)
	auth.Post("/verify", h.VerifyToken)

	// Protected routes (require authentication)
	protected := v1.Group("")
	protected.Use(h.VerifyToken) // Auth middleware

	// User routes
	users := protected.Group("/users")
	users.Get("/profile", h.GetUserProfile)
	users.Put("/profile", h.UpdateUserProfile)

	// Mod processing routes
	mods := protected.Group("/mods")
	mods.Post("/upload", h.UploadMod)
	mods.Get("/jobs/:id", h.GetJobStatus)
	mods.Post("/jobs/:id/process", h.ProcessMod)
	mods.Get("/jobs/:id/download", h.DownloadMod)
	mods.Get("/jobs", h.GetUserJobs)

	// Mod presets
	presets := v1.Group("/presets") // Public endpoints
	presets.Get("/", h.GetPresets)
	presets.Get("/:type", h.GetPresetsByType)

	// Credits and billing (protected)
	billing := protected.Group("/billing")
	billing.Get("/credits", h.GetCredits)
	billing.Post("/credits/purchase", h.PurchaseCredits)
}
