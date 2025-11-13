# Quick Start: AI Video Generation

Get professional marketing videos in 3 simple steps!

## Step 1: Choose Your AI Service (5 minutes)

### Option A: D-ID (Recommended) 
**Best for**: MVP, Quick Start, Budget-friendly

1. Go to https://studio.d-id.com/
2. Sign up (free 20 credits)
3. Get API key: Account Settings → API Key → Copy

```bash
cd /Users/slowermaster/DEALSHARE/hacathon/backend
echo 'AI_PROVIDER=did' >> .env
echo 'AI_API_KEY=ZGluZXNoLnJhbUBkZWFsc2hhcmUuaW4:31Ln-cQ0VJ-_sj61DbZm3' >> .env
```

### Option B: Synthesia
**Best for**: Professional quality

1. Go to https://www.synthesia.io/
2. Sign up for trial
3. Request API access (takes 1-2 days)

### Option C: RunwayML
**Best for**: Creative videos

1. Go to https://app.runwayml.com/
2. Sign up
3. Get API key from Settings

## Step 2: Restart Backend

```bash
# Stop current server (Ctrl+C in the terminal running it)
# Or run this:
ps aux | grep "go run main.go" | grep -v grep | awk '{print $2}' | xargs kill

# Start with new config
cd /Users/slowermaster/DEALSHARE/hacathon/backend
go run main.go
```

## Step 3: Generate Your First Video

1. Open http://localhost:3000
2. Upload a product image
3. Upload a person photo/video
4. Click "Upload & Create Project"
5. Click "Generate Video" 
6. Wait 2-5 minutes
7. Click "View Video" when complete!

## What You Get

### Professional Marketing Video
- AI-generated presenter
- Natural voice narration
- Product showcase
- Professional quality
- Social media ready

### Marketing Website  
- Modern responsive design
- Hero section with product
- Feature highlights
- Video player
- Call-to-action buttons
- Mobile-friendly

## Sample Scripts

### For Product Demo (D-ID)
Edit in `backend/internal/services/video_generator.go` line ~80:

```go
"input": "Welcome! I'm excited to show you this amazing product. It's designed to solve your biggest challenges with innovative features and intuitive design. Let me walk you through what makes it special!",
```

### For Professional Presentation (Synthesia)
Edit line ~150:

```go
"scriptText": "Good morning! Today I'm presenting our flagship product - a revolutionary solution that transforms how businesses operate. With cutting-edge technology and user-centric design, we're setting new industry standards.",
```

## Costs

| Service | Free Tier | First Video | Monthly Plan |
|---------|-----------|-------------|--------------|
| D-ID | 20 credits free | $0 | $5.90 (60 videos) |
| Synthesia | Demo only | - | $22 (10 videos) |
| RunwayML | 125 credits | $0 | $12 (variable) |

## Troubleshooting

### "Video generation failed"
```bash
# Check your API key is set
cd backend
cat .env | grep AI_

# Should see:
# AI_PROVIDER=did
# AI_API_KEY=your_key_here
```

### "Connection timeout"
- Video generation takes 2-10 minutes
- Don't close the browser
- Check terminal logs for progress

### "Invalid API key"
- Make sure you copied the full key
- No spaces before/after
- Regenerate key if needed

## Testing Without API (Mock Mode)

```bash
# Use mock mode for testing
cd backend
export AI_PROVIDER=mock
go run main.go
```

This creates placeholder videos to test the workflow.

## Next Steps

Once your first video works:

1. **Customize the script** - Edit video_generator.go
2. **Try different voices** - Check AI service docs
3. **Adjust video length** - Configure duration settings
4. **Add branding** - Customize website templates
5. **Deploy** - See deployment guide

## Full Documentation

- **AI Services**: See `AI_SERVICES_GUIDE.md`
- **Backend API**: See `backend/README_BACKEND.md`
- **Main Setup**: See `README.md`

## Support

- Check terminal logs for errors
- Verify API keys are active
- Review AI service status pages
- Test with small files first

---

**🎉 You're Ready!**

Your AI-powered product marketing system is set up and ready to generate professional marketing videos and websites automatically!


