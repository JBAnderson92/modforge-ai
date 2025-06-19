#!/bin/bash

# 🚀 ModForge.ai Deployment Script
# This script helps you deploy ModForge.ai to various platforms

set -e

echo "🚀 ModForge.ai Deployment Helper"
echo "================================"

# Check if user has required tools
check_requirements() {
    echo "📋 Checking requirements..."
    
    # Check for Git
    if ! command -v git &> /dev/null; then
        echo "❌ Git is required but not installed"
        exit 1
    fi
    
    # Check for Go
    if ! command -v go &> /dev/null; then
        echo "❌ Go is required but not installed"
        exit 1
    fi
    
    # Check for Node.js
    if ! command -v node &> /dev/null; then
        echo "❌ Node.js is required but not installed"
        exit 1
    fi
    
    echo "✅ All requirements met!"
}

# Deploy to Railway
deploy_railway() {
    echo "🚂 Deploying to Railway..."
    
    # Check if Railway CLI is installed
    if ! command -v railway &> /dev/null; then
        echo "Installing Railway CLI..."
        curl -fsSL https://railway.app/install.sh | sh
    fi
    
    # Login to Railway
    echo "Please login to Railway:"
    railway login
    
    # Create new project
    railway init
    
    # Set environment variables
    echo "Setting up environment variables..."
    echo "Please set these variables in your Railway dashboard:"
    echo "- OPENAI_API_KEY"
    echo "- CLOUDFLARE_R2_ACCOUNT_ID" 
    echo "- CLOUDFLARE_R2_API_TOKEN"
    echo "- CLOUDFLARE_R2_BUCKET_NAME"
    echo "- CLOUDFLARE_R2_REGION"
    echo "- ALLOWED_ORIGINS"
    
    # Deploy
    railway up
    
    echo "✅ Railway deployment started!"
    echo "Visit https://railway.app/dashboard to monitor your deployment"
}

# Deploy to Render
deploy_render() {
    echo "🎨 Deploying to Render..."
    echo "1. Go to https://render.com and create an account"
    echo "2. Connect your GitHub repository"
    echo "3. Create a new Web Service"
    echo "4. Use these settings:"
    echo "   - Build Command: go build -o main ./api/main.go"
    echo "   - Start Command: ./main"
    echo "   - Environment: Go"
    echo "5. Add environment variables in Render dashboard"
    echo "✅ Follow the instructions above to deploy to Render"
}

# Deploy to Fly.io
deploy_fly() {
    echo "🪰 Deploying to Fly.io..."
    
    # Check if Fly CLI is installed
    if ! command -v flyctl &> /dev/null; then
        echo "Installing Fly CLI..."
        curl -L https://fly.io/install.sh | sh
    fi
    
    # Login to Fly
    flyctl auth login
    
    # Launch app
    flyctl launch --no-deploy
    
    # Set secrets
    echo "Setting up secrets..."
    flyctl secrets set OPENAI_API_KEY="$OPENAI_API_KEY"
    flyctl secrets set CLOUDFLARE_R2_ACCOUNT_ID="$CLOUDFLARE_R2_ACCOUNT_ID"
    flyctl secrets set CLOUDFLARE_R2_API_TOKEN="$CLOUDFLARE_R2_API_TOKEN"
    flyctl secrets set CLOUDFLARE_R2_BUCKET_NAME="$CLOUDFLARE_R2_BUCKET_NAME"
    flyctl secrets set CLOUDFLARE_R2_REGION="$CLOUDFLARE_R2_REGION"
    
    # Deploy
    flyctl deploy
    
    echo "✅ Fly.io deployment complete!"
}

# Deploy frontend to Vercel
deploy_frontend() {
    echo "🌐 Deploying frontend to Vercel..."
    
    cd frontend
    
    # Check if Vercel CLI is installed
    if ! command -v vercel &> /dev/null; then
        echo "Installing Vercel CLI..."
        npm install -g vercel
    fi
    
    # Build frontend
    npm run build
    
    # Deploy to Vercel
    vercel --prod
    
    cd ..
    echo "✅ Frontend deployed to Vercel!"
}

# Setup production environment
setup_production() {
    echo "🔧 Setting up production environment..."
    
    # Create production .env
    cp .env.example .env.production
    
    echo "Please edit .env.production with your production values:"
    echo "- Set DATABASE_URL to production database"
    echo "- Set ALLOWED_ORIGINS to your frontend URL" 
    echo "- Verify all API keys are production keys"
    
    read -p "Press enter when ready to continue..."
    
    echo "✅ Production environment configured!"
}

# Build for production
build_production() {
    echo "🏗️  Building for production..."
    
    # Build backend
    echo "Building Go backend..."
    CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o modforge-api ./api/main.go
    
    # Build frontend
    echo "Building React frontend..."
    cd frontend
    npm run build
    cd ..
    
    echo "✅ Production build complete!"
}

# Main menu
show_menu() {
    echo ""
    echo "Choose deployment option:"
    echo "1) Railway (Recommended for beginners)"
    echo "2) Render (Great free tier)"
    echo "3) Fly.io (Fast global deployment)"
    echo "4) Deploy frontend to Vercel"
    echo "5) Setup production environment"
    echo "6) Build for production"
    echo "7) Check requirements"
    echo "0) Exit"
    echo ""
    read -p "Enter your choice: " choice
    
    case $choice in
        1) deploy_railway ;;
        2) deploy_render ;;
        3) deploy_fly ;;
        4) deploy_frontend ;;
        5) setup_production ;;
        6) build_production ;;
        7) check_requirements ;;
        0) echo "Goodbye! 👋"; exit 0 ;;
        *) echo "Invalid option"; show_menu ;;
    esac
}

# Start script
check_requirements
show_menu
