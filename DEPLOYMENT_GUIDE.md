# 🚀 ModForge.ai Deployment Guide

## Make Your App Live in Production!

This guide will help you deploy ModForge.ai to make it accessible to real users worldwide.

## 🎯 Quick Start (Recommended)

### Option 1: Railway + Vercel (Easiest)

**Backend on Railway:**

1. Install Railway CLI:

   ```bash
   curl -fsSL https://railway.app/install.sh | sh
   ```

2. Login and deploy:

   ```bash
   railway login
   railway init
   railway up
   ```

3. Set environment variables in Railway dashboard:
   - `OPENAI_API_KEY`
   - `CLOUDFLARE_R2_ACCOUNT_ID`
   - `CLOUDFLARE_R2_API_TOKEN`
   - `CLOUDFLARE_R2_BUCKET_NAME`
   - `ALLOWED_ORIGINS`

**Frontend on Vercel:**

1. Install Vercel CLI:

   ```bash
   npm install -g vercel
   ```

2. Deploy frontend:

   ```bash
   cd frontend
   vercel --prod
   ```

3. Set `VITE_API_URL` to your Railway backend URL

**✅ Your app is now live!**

## 🛠️ Deployment Options

### Option A: Railway (Recommended for MVP)

- **Cost**: $5/month
- **Pros**: Zero config, automatic SSL, database included
- **Best for**: Getting started quickly

### Option B: Render

- **Cost**: Free tier available
- **Pros**: Great free tier, easy setup
- **Best for**: Testing and small projects

### Option C: Fly.io

- **Cost**: Pay-as-you-go
- **Pros**: Global edge deployment, fast
- **Best for**: Production apps with global users

### Option D: DigitalOcean App Platform

- **Cost**: $5/month
- **Pros**: Simple scaling, managed database
- **Best for**: Growing applications

## 🔧 Pre-Deployment Checklist

### 1. Environment Setup

```bash
# Copy production environment
cp .env.example .env.production

# Edit with production values
nano .env.production
```

### 2. Database Migration

```bash
# Ensure migrations are ready
ls migrations/
```

### 3. Build Test

```bash
# Test production build
./deploy.sh
# Choose option 6: Build for production
```

### 4. Security Check

- ✅ Remove any test API keys
- ✅ Set strong database passwords
- ✅ Configure CORS for your domain
- ✅ Enable HTTPS only

## 🌍 Domain Setup

### 1. Get a Domain

- [Namecheap](https://namecheap.com) - $10-15/year
- [Cloudflare](https://cloudflare.com) - $10/year
- [Google Domains](https://domains.google.com) - $12/year

### 2. Configure DNS

Point your domain to your deployment:

```
A record: @ -> your-server-ip
CNAME: www -> your-app.railway.app
```

### 3. SSL Certificate

Most platforms provide automatic SSL:

- Railway: Automatic
- Vercel: Automatic
- Render: Automatic

## 📊 Production Monitoring

### 1. Error Tracking

Add Sentry for error monitoring:

```bash
npm install @sentry/react @sentry/node
```

### 2. Analytics

Add Google Analytics or Plausible:

```html
<!-- Add to frontend/index.html -->
<script
  async
  src="https://www.googletagmanager.com/gtag/js?id=GA_MEASUREMENT_ID"
></script>
```

### 3. Health Checks

Your app includes health check endpoint:

```
GET /api/v1/health
```

## 💰 Cost Estimation

### Minimal Setup (MVP)

- **Domain**: $10/year
- **Backend (Railway)**: $5/month
- **Frontend (Vercel)**: Free
- **Storage (Cloudflare R2)**: $0.015/GB
- **OpenAI API**: Pay per use
- **Total**: ~$70/year + usage

### Growing App

- **Backend**: $20/month (more resources)
- **Database**: $15/month (managed)
- **CDN**: $10/month
- **Monitoring**: $10/month
- **Total**: ~$660/year

## 🚀 Scaling Strategy

### Phase 1: MVP (0-100 users)

- Single server deployment
- SQLite database
- Basic monitoring

### Phase 2: Growth (100-1k users)

- Load balancer
- PostgreSQL database
- Redis caching
- File CDN

### Phase 3: Scale (1k+ users)

- Multiple regions
- Microservices
- Advanced monitoring
- Auto-scaling

## 🔐 Security Best Practices

### 1. Environment Variables

```bash
# Never commit secrets to git
echo ".env*" >> .gitignore
```

### 2. API Rate Limiting

```go
// Already implemented in middleware
app.Use(limiter.New())
```

### 3. File Upload Security

- ✅ File type validation
- ✅ Size limits
- ✅ Virus scanning (VirusTotal)

### 4. HTTPS Only

```go
// Redirect HTTP to HTTPS in production
if os.Getenv("HTTPS_ONLY") == "true" {
    app.Use(func(c *fiber.Ctx) error {
        if !c.Secure() {
            return c.Redirect("https://" + c.Hostname() + c.OriginalURL())
        }
        return c.Next()
    })
}
```

## 📱 Mobile Responsiveness

Your app is already mobile-ready with Tailwind CSS:

- ✅ Responsive design
- ✅ Touch-friendly interface
- ✅ Fast loading

## 🧪 Testing in Production

### 1. Smoke Tests

```bash
# Test health endpoint
curl https://your-app.com/api/v1/health

# Test file upload
curl -X POST https://your-app.com/api/v1/mods/upload \
  -F "mod_file=@test.json"
```

### 2. Load Testing

```bash
# Install artillery
npm install -g artillery

# Run load test
artillery quick --count 10 --num 100 https://your-app.com
```

## 🎉 Go Live Checklist

- [ ] Deploy backend to Railway/Render/Fly
- [ ] Deploy frontend to Vercel/Netlify
- [ ] Configure custom domain
- [ ] Set up SSL certificate
- [ ] Test all functionality
- [ ] Set up monitoring
- [ ] Configure backups
- [ ] Update DNS records
- [ ] Test from different devices
- [ ] Announce launch! 🚀

## 📞 Support & Maintenance

### 1. Monitoring

- Set up uptime monitoring (UptimeRobot)
- Configure error alerts
- Monitor API usage and costs

### 2. Backups

```bash
# Automated database backups
crontab -e
0 2 * * * /usr/bin/sqlite3 /app/modforge.db ".backup /backups/modforge-$(date +%Y%m%d).db"
```

### 3. Updates

```bash
# Deploy updates
git push origin main
# Platform auto-deploys from main branch
```

## 🌟 Launch Strategy

1. **Soft Launch**: Share with friends and beta testers
2. **Product Hunt**: Submit to Product Hunt
3. **Social Media**: Share on Twitter, Reddit, Discord
4. **Blog Posts**: Write about your building journey
5. **Communities**: Share in developer and gaming communities

Your ModForge.ai app is now ready for the world! 🌍✨
