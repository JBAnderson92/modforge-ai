# ModForge.ai - AI-Powered Game Modding Platform

🤖 Transform your game mods with artificial intelligence. Upload Minecraft, Skyrim, or Lua mods and let AI rebalance, translate, enhance, and transform them into something extraordinary.

## 🎯 Features

### Core MVP Features

- **Multi-Game Support**: Minecraft JSON mods, Skyrim ESP files, Lua scripts
- **AI-Powered Transformations**: Rewrite content, rebalance stats, translate languages
- **Intelligent Presets**: Pre-built templates for common modding tasks
- **Secure File Processing**: Virus scanning and sandboxed execution
- **Modern Web Interface**: React + Tailwind CSS frontend
- **RESTful API**: Go + Fiber backend with proper error handling

### AI Presets (MVP)

- **Make Lore-Friendly** (1 credit): Rewrite item names and descriptions
- **Balance Items** (2 credits): Adjust stats, durability, and recipes
- **Translate Mod** (1 credit): Convert content to any language
- **Add Content** (3 credits): Generate new items, blocks, or recipes

## 🚀 Quick Start

### Prerequisites

- Go 1.21+
- Node.js 18+ and npm
- Redis (for development)

### 1. Setup Environment

```bash
# Copy environment template
cp .env.example .env

# Edit .env with your API keys:
# - OPENAI_API_KEY (required for AI features)
# - CLOUDFLARE_R2_* (required for file storage)
# - FIREBASE_CONFIG (required for auth)
```

### 2. Install Dependencies

```bash
# Backend dependencies
go mod tidy

# Frontend dependencies
cd frontend && npm install && cd ..
```

### 3. Run Development Servers

```bash
# Terminal 1: Start API server
go run main.go api

# Terminal 2: Start frontend dev server
go run main.go frontend

# Or see all available commands:
go run main.go
```

### 4. Access the Application

- **Frontend**: http://localhost:5173
- **API**: http://localhost:8080
- **Health Check**: http://localhost:8080/api/v1/health

1. Clone or download this project
2. Navigate to the project directory
3. The Go module is already initialized

### Running the Application

To run the main application:

```bash
go run main.go
```

To build the application:

```bash
go build -o project1
./project1
```

### Development

This project is configured for VS Code with:

- Go extension support
- Build and run tasks
- Proper Go formatting and linting

### Project Structure

```
.
├── .github/
│   └── copilot-instructions.md  # Copilot workspace instructions
├── .vscode/
│   └── tasks.json               # VS Code build/run tasks
├── go.mod                       # Go module definition
├── main.go                      # Main application entry point
└── README.md                    # This file
```

## Next Steps

- Add your business logic to `main.go`
- Create additional packages as needed
- Add tests in `*_test.go` files
- Configure any external dependencies in `go.mod`

## Commands

- `go run main.go` - Run the application
- `go build` - Build the application
- `go test ./...` - Run all tests
- `go mod tidy` - Clean up module dependencies
- `gofmt -w .` - Format all Go files
