# 🔑 ModForge.ai API Keys Setup Guide

This guide will help you obtain and configure all the required API keys for ModForge.ai.

## Required API Keys

### 1. 🤖 OpenAI API Key (REQUIRED)

**Purpose**: Powers the AI mod transformations using GPT-4-turbo

**How to get it**:

1. Go to [OpenAI Platform](https://platform.openai.com/)
2. Sign up or log in to your account
3. Navigate to API Keys section
4. Click "Create new secret key"
5. Copy the key (starts with `sk-`)

**Cost**: Pay-per-use, ~$0.01-0.03 per 1K tokens
**Free Tier**: $5 free credit for new accounts

### 2. ☁️ Cloudflare R2 Storage (REQUIRED)

**Purpose**: Stores uploaded mod files and processed results

**How to get it**:

1. Go to [Cloudflare Dashboard](https://dash.cloudflare.com/)
2. Sign up/login and navigate to R2 Object Storage
3. Create a new R2 bucket named `modforge-files`
4. Go to "Manage R2 API tokens"
5. Create new token with Object Read & Write permissions
6. Note down: Account ID, Access Key ID, Secret Access Key

**Cost**: Very cheap - $0.015/GB/month storage, $4.50/million requests
**Free Tier**: 10GB storage, 1M Class A operations per month

### 3. 🔐 Firebase Authentication (REQUIRED)

**Purpose**: User authentication and session management

**How to get it**:

1. Go to [Firebase Console](https://console.firebase.google.com/)
2. Create new project or use existing
3. Enable Authentication with Email/Password
4. Go to Project Settings > Service Accounts
5. Generate new private key (downloads JSON file)
6. Also get your Web App config from Project Settings > General

**Cost**: Free for up to 50K MAU (Monthly Active Users)

### 4. 🛡️ VirusTotal API (RECOMMENDED)

**Purpose**: Scans uploaded files for malware

**How to get it**:

1. Go to [VirusTotal](https://www.virustotal.com/)
2. Sign up for free account
3. Go to your profile and get API key
4. Free tier: 4 requests/minute, 500 requests/day

**Cost**: Free tier available, paid plans for higher limits

### 5. 🗄️ Redis (DEVELOPMENT)

**Purpose**: Rate limiting and job queue management

**For Development**:

```bash
# macOS
brew install redis
brew services start redis

# Ubuntu/Debian
sudo apt install redis-server
sudo systemctl start redis

# Docker
docker run -d -p 6379:6379 redis:alpine
```

**For Production**: Use managed Redis (AWS ElastiCache, Google Cloud Memorystore, etc.)

## Setup Instructions

### Step 1: Copy Environment File

```bash
cp .env.example .env
```

### Step 2: Edit .env File

Open `.env` and replace the placeholder values with your actual API keys.

### Step 3: Test Configuration

```bash
# Test that all services are accessible
go run main.go api
```

## Security Best Practices

1. **Never commit .env to version control** (already in .gitignore)
2. **Use different keys for development/production**
3. **Rotate keys regularly**
4. **Set up billing alerts for paid services**
5. **Use least-privilege access for all tokens**

## Troubleshooting

### OpenAI Issues

- Check billing setup and credit balance
- Verify key format (starts with `sk-`)
- Ensure you have access to GPT-4 models

### Cloudflare R2 Issues

- Verify bucket name matches exactly
- Check R2 token permissions include Object Read/Write
- Ensure account ID is correct

### Firebase Issues

- Verify service account JSON path is correct
- Check Firebase project has Authentication enabled
- Ensure web app configuration is correct

### Redis Issues

- Check if Redis server is running: `redis-cli ping`
- Verify connection URL format
- For Docker: ensure port 6379 is not in use

## Cost Estimation (Monthly)

**Development**: ~$0-5/month

- OpenAI: $1-3 (light testing)
- Cloudflare R2: Free tier
- Firebase: Free tier
- VirusTotal: Free tier

**Production (1000 users)**: ~$20-50/month

- OpenAI: $15-30 (depends on usage)
- Cloudflare R2: $5-10
- Firebase: Free (under 50K users)
- VirusTotal: $10 (if exceeding free tier)

## Next Steps

1. Set up all API keys using this guide
2. Test the application: `go run main.go api`
3. Verify health endpoint: `curl http://localhost:8080/api/v1/health`
4. Start frontend development server
5. Test end-to-end file upload flow
