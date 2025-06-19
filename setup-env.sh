#!/bin/bash

# ModForge.ai Environment Setup Script
echo "🚀 ModForge.ai Environment Setup"
echo "================================="
echo ""

# Create .env file from template
if [ ! -f .env ]; then
    cp .env.example .env
    echo "✅ Created .env file from template"
else
    echo "ℹ️  .env file already exists"
fi

echo ""
echo "🔑 API Keys Setup"
echo "=================="
echo ""

# Function to update env file
update_env() {
    local key=$1
    local value=$2
    if [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS
        sed -i '' "s|^${key}=.*|${key}=${value}|" .env
    else
        # Linux
        sed -i "s|^${key}=.*|${key}=${value}|" .env
    fi
}

# OpenAI API Key
echo "1. 🤖 OpenAI API Key"
echo "   Get yours at: https://platform.openai.com/api-keys"
echo "   Format: sk-..."
read -p "   Enter your OpenAI API key (or press Enter to skip): " openai_key
if [ ! -z "$openai_key" ]; then
    update_env "OPENAI_API_KEY" "$openai_key"
    echo "   ✅ OpenAI API key saved"
else
    echo "   ⚠️  Skipped OpenAI API key"
fi
echo ""

# Cloudflare R2
echo "2. ☁️  Cloudflare R2 Storage"
echo "   Get yours at: https://dash.cloudflare.com/profile/api-tokens"
read -p "   Enter your Cloudflare Account ID (or press Enter to skip): " cf_account_id
if [ ! -z "$cf_account_id" ]; then
    update_env "CLOUDFLARE_R2_ACCOUNT_ID" "$cf_account_id"
    
    read -p "   Enter your R2 Access Key ID: " cf_access_key
    update_env "CLOUDFLARE_R2_ACCESS_KEY_ID" "$cf_access_key"
    
    read -p "   Enter your R2 Secret Access Key: " cf_secret_key
    update_env "CLOUDFLARE_R2_SECRET_ACCESS_KEY" "$cf_secret_key"
    
    read -p "   Enter your R2 Bucket Name [modforge-files]: " cf_bucket
    cf_bucket=${cf_bucket:-modforge-files}
    update_env "CLOUDFLARE_R2_BUCKET_NAME" "$cf_bucket"
    
    echo "   ✅ Cloudflare R2 configuration saved"
else
    echo "   ⚠️  Skipped Cloudflare R2 configuration"
fi
echo ""

# Firebase
echo "3. 🔐 Firebase Authentication"
echo "   Get yours at: https://console.firebase.google.com/"
read -p "   Enter path to Firebase service account JSON file (or press Enter to skip): " firebase_config
if [ ! -z "$firebase_config" ]; then
    update_env "FIREBASE_CONFIG" "$firebase_config"
    echo "   ✅ Firebase configuration saved"
else
    echo "   ⚠️  Skipped Firebase configuration"
fi
echo ""

# VirusTotal
echo "4. 🛡️  VirusTotal API (Optional)"
echo "   Get yours at: https://www.virustotal.com/gui/my-apikey"
read -p "   Enter your VirusTotal API key (or press Enter to skip): " vt_key
if [ ! -z "$vt_key" ]; then
    update_env "VIRUSTOTAL_API_KEY" "$vt_key"
    echo "   ✅ VirusTotal API key saved"
else
    echo "   ⚠️  Skipped VirusTotal API key"
fi
echo ""

# Check Redis
echo "5. 🗄️  Redis Setup"
if command -v redis-cli &> /dev/null; then
    if redis-cli ping &> /dev/null; then
        echo "   ✅ Redis is running locally"
    else
        echo "   ⚠️  Redis is installed but not running"
        echo "   Start it with: brew services start redis (macOS) or sudo systemctl start redis (Linux)"
    fi
else
    echo "   ❌ Redis not found"
    echo "   Install it with:"
    echo "   - macOS: brew install redis"
    echo "   - Ubuntu: sudo apt install redis-server"
    echo "   - Docker: docker run -d -p 6379:6379 redis:alpine"
fi
echo ""

echo "🎉 Setup Complete!"
echo "=================="
echo ""
echo "Next steps:"
echo "1. Start Redis if not running"
echo "2. Test the API: go run main.go api"
echo "3. Install frontend deps: cd frontend && npm install"
echo "4. Start frontend: go run main.go frontend"
echo ""
echo "For detailed setup instructions, see: API_KEYS_SETUP.md"
