# 🚀 One-Click Deploy

Deploy ModForge.ai to Railway with one click:

[![Deploy on Railway](https://railway.app/button.svg)](https://railway.app/template/8bDfGT?referralCode=alphasec)

## What gets deployed:

- ✅ Go backend API server
- ✅ SQLite database with migrations
- ✅ File upload and AI processing
- ✅ Automatic HTTPS and domain
- ✅ Environment variables setup

## After deployment:

1. **Set API Keys** in Railway dashboard:

   - `OPENAI_API_KEY`
   - `CLOUDFLARE_R2_ACCOUNT_ID`
   - `CLOUDFLARE_R2_API_TOKEN`
   - `CLOUDFLARE_R2_BUCKET_NAME`

2. **Deploy Frontend** to Vercel:

   ```bash
   cd frontend
   vercel --prod
   ```

3. **Update CORS** in Railway:

   - Set `ALLOWED_ORIGINS` to your Vercel URL

4. **Test Your App**:
   - Visit your Railway URL
   - Upload a test mod file
   - Verify AI processing works

Your app will be live in under 5 minutes! 🎉
