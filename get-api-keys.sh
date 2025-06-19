#!/bin/bash

echo "🔑 Quick API Keys Setup for ModForge.ai"
echo "========================================"
echo ""
echo "I'll help you get the required API keys. Here are the direct links:"
echo ""

echo "1. 🤖 OPENAI API KEY (REQUIRED)"
echo "   → https://platform.openai.com/api-keys"
echo "   → Click 'Create new secret key'"
echo "   → Copy the key (starts with sk-)"
echo ""

echo "2. ☁️  CLOUDFLARE R2 STORAGE (REQUIRED)"
echo "   → https://dash.cloudflare.com/"
echo "   → Go to R2 Object Storage → Create bucket 'modforge-files'"
echo "   → Manage R2 API tokens → Create token"
echo "   → Copy: Account ID, Access Key ID, Secret Access Key"
echo ""

echo "3. 🔐 FIREBASE AUTH (REQUIRED)"
echo "   → https://console.firebase.google.com/"
echo "   → Create project → Authentication → Enable Email/Password"
echo "   → Project Settings → Service Accounts → Generate private key"
echo "   → Download the JSON file"
echo ""

echo "4. 🛡️  VIRUSTOTAL API (OPTIONAL)"
echo "   → https://www.virustotal.com/gui/my-apikey"
echo "   → Sign up → Get your API key"
echo ""

echo "5. 🗄️  REDIS (DEVELOPMENT)"
if command -v brew &> /dev/null; then
    echo "   → Run: brew install redis && brew services start redis"
elif command -v apt &> /dev/null; then
    echo "   → Run: sudo apt install redis-server && sudo systemctl start redis"
else
    echo "   → Run: docker run -d -p 6379:6379 redis:alpine"
fi
echo ""

echo "📝 NEXT STEPS:"
echo "1. Get your API keys from the links above"
echo "2. Edit the .env file and replace the placeholder values"
echo "3. Test with: go run main.go api"
echo ""

# Check if user wants to open the links
read -p "Open API key signup pages in browser? (y/n): " open_links
if [[ $open_links == "y" || $open_links == "Y" ]]; then
    echo "Opening signup pages..."
    
    if command -v open &> /dev/null; then
        # macOS
        open "https://platform.openai.com/api-keys"
        sleep 2
        open "https://dash.cloudflare.com/"
        sleep 2
        open "https://console.firebase.google.com/"
        sleep 2
        open "https://www.virustotal.com/gui/my-apikey"
    elif command -v xdg-open &> /dev/null; then
        # Linux
        xdg-open "https://platform.openai.com/api-keys"
        sleep 2
        xdg-open "https://dash.cloudflare.com/"
        sleep 2
        xdg-open "https://console.firebase.google.com/"
        sleep 2
        xdg-open "https://www.virustotal.com/gui/my-apikey"
    else
        echo "Please manually open the URLs above in your browser"
    fi
fi

echo ""
echo "💡 TIP: Start with just the OpenAI API key to test basic functionality!"
echo "    You can add the others later as needed."
