package main

import (
	"log"
	"os"

	"modforge.ai/ai"
	"modforge.ai/api/config"
	"modforge.ai/api/database"
	"modforge.ai/api/middleware"
	"modforge.ai/api/routes"
	"modforge.ai/api/storage"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.Initialize(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := database.RunMigrations(cfg.DatabaseURL); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize storage client
	storageClient, err := storage.NewClient(storage.Config{
		AccountID:  cfg.CloudflareR2.AccountID,
		APIToken:   cfg.CloudflareR2.APIToken,
		BucketName: cfg.CloudflareR2.BucketName,
		Region:     cfg.CloudflareR2.Region,
	})
	if err != nil {
		log.Fatalf("Failed to initialize storage client: %v", err)
	}

	// Initialize AI client
	aiClient := ai.NewClient(cfg.OpenAIAPIKey)

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})

	// Add middleware
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.AllowedOrigins,
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	// Initialize routes
	routes.Setup(app, db, cfg, storageClient, aiClient)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("ModForge.ai API server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
