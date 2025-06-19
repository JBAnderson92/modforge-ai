package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {
	fmt.Println("🚀 ModForge.ai - AI-Powered Game Modding Platform")
	fmt.Println("==============================================")

	// Check if we should run API server or show help
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "api":
			fmt.Println("Starting API server...")
			runAPIServer()
		case "frontend":
			fmt.Println("Starting frontend development server...")
			runFrontend()
		case "dev":
			fmt.Println("Starting development environment (API + Frontend)...")
			runDevEnvironment()
		default:
			showHelp()
		}
	} else {
		// Default to API server in production, help in development
		if os.Getenv("RAILWAY_ENVIRONMENT") != "" || os.Getenv("PORT") != "" {
			fmt.Println("Production environment detected - starting API server...")
			runAPIServer()
		} else {
			showHelp()
		}
	}
}

func showHelp() {
	fmt.Print(`
Available commands:
  go run main.go api       - Start the API server
  go run main.go frontend  - Start the frontend development server
  go run main.go dev       - Start both API and frontend in development mode

For development setup:
  1. Copy .env.example to .env and configure your keys
  2. Run 'go run main.go dev' to start the development environment
  3. Access the frontend at http://localhost:5173
  4. API will be available at http://localhost:8080

Environment Setup:
  - OpenAI API key required for AI features
  - Cloudflare R2 credentials for file storage
  - Firebase config for authentication
  - VirusTotal API key for file scanning
`)
}

func runAPIServer() {
	// Check if we're in production (no Go available)
	if !isCommandAvailable("go") {
		// In production, we should have a pre-built binary
		// For now, let's just redirect to the api binary
		fmt.Println("Production mode detected - looking for API binary...")
		
		// Try to find and run the API binary
		if _, err := os.Stat("./api-server"); err == nil {
			cmd := exec.Command("./api-server")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				log.Fatalf("Failed to start API server: %v", err)
			}
		} else {
			log.Fatal("API server binary not found in production mode")
		}
	} else {
		// Development mode
		cmd := exec.Command("go", "run", "./api/main.go")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			log.Fatalf("Failed to start API server: %v", err)
		}
	}
}

func runFrontend() {
	// Check if Node.js is installed
	if !isCommandAvailable("npm") {
		log.Fatal("npm is required to run the frontend. Please install Node.js.")
	}

	// Check if frontend is set up
	if _, err := os.Stat("frontend/package.json"); os.IsNotExist(err) {
		log.Fatal("Frontend not set up. Run the setup script first.")
	}

	cmd := exec.Command("npm", "run", "dev")
	cmd.Dir = "frontend"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to start frontend: %v", err)
	}
}

func runDevEnvironment() {
	fmt.Println("Development environment will start API server on :8080 and frontend on :5173")
	fmt.Println("For now, run them separately:")
	fmt.Println("Terminal 1: go run main.go api")
	fmt.Println("Terminal 2: go run main.go frontend")
}

func isCommandAvailable(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}
