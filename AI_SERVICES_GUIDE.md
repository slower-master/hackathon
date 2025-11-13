# AI Services Integration Guide

This guide explains how to integrate various AI video generation services with your product marketing agent.

## Supported AI Services

### 1. **D-ID** (Recommended for MVP) ⭐
**Best for**: Talking head videos, avatar presentations

- **Website**: https://www.d-id.com/
- **Pricing**: Free tier available, $5.90/month for 20 credits
- **Setup Time**: 5 minutes
- **Quality**: Excellent for spokesperson-style videos

**Setup Steps**:
```bash
# 1. Sign up at https://studio.d-id.com/
# 2. Get your API key from https://studio.d-id.com/account-settings
# 3. Set environment variables

export AI_PROVIDER=did
export AI_API_KEY="your_d_id_api_key_here"
```

**Features**:
- Realistic talking avatars
- Multiple voice options
- 25+ languages
- Custom scripts
- Fast generation (2-5 minutes)

---

### 2. **Synthesia** 
**Best for**: Professional AI avatar presentations

- **Website**: https://www.synthesia.io/
- **Pricing**: Starts at $22/month
- **Setup Time**: 10 minutes
- **Quality**: Professional-grade for corporate videos

**Setup Steps**:
```bash
# 1. Sign up at https://www.synthesia.io/
# 2. Get API access
# 3. Set environment variables

export AI_PROVIDER=synthesia
export AI_API_KEY="your_synthesia_api_key"
```

**Features**:
- 140+ AI avatars
- 120+ languages
- Custom backgrounds
- Professional quality
- Longer videos (up to 30 min)

---

### 3. **RunwayML Gen-2**
**Best for**: Creative, cinematic video generation

- **Website**: https://runwayml.com/
- **Pricing**: $12/month for 625 credits
- **Setup Time**: 10 minutes
- **Quality**: Highly creative, artistic

**Setup Steps**:
```bash
# 1. Sign up at https://app.runwayml.com/
# 2. Get API key from account settings
# 3. Set environment variables

export AI_PROVIDER=runwayml
export AI_API_KEY="your_runwayml_api_key"
```

**Features**:
- Image-to-video generation
- Video-to-video transformation
- Motion brush
- Text prompts
- 4-18 second clips

---

## Quick Setup (Choose One)

### Option 1: D-ID (Easiest & Cheapest)

1. **Get API Key**:
   ```bash
   # Visit https://studio.d-id.com/
   # Create account → Account Settings → API Key
   ```

2. **Configure Backend**:
   ```bash
   cd backend
   
   # Create .env file
   cat > .env << EOF
   PORT=8080
   AI_PROVIDER=did
   AI_API_KEY=your_d_id_api_key_here
   DATABASE_PATH=./data/app.db
   UPLOAD_PATH=./uploads
   GENERATED_VIDEO_PATH=./generated/videos
   WEBSITE_PATH=./generated/websites
   EOF
   ```

3. **Restart Backend**:
   ```bash
   # Stop current server (Ctrl+C)
   go run main.go
   ```

### Option 2: Synthesia (More Professional)

1. **Get API Access**:
   - Sign up at https://www.synthesia.io/
   - Request API access through support
   - Wait for approval (usually 1-2 days)

2. **Configure**:
   ```bash
   export AI_PROVIDER=synthesia
   export AI_API_KEY="your_synthesia_key"
   ```

### Option 3: RunwayML (Most Creative)

1. **Get API Key**:
   - Sign up at https://app.runwayml.com/
   - Navigate to Settings → API Keys
   - Generate new key

2. **Configure**:
   ```bash
   export AI_PROVIDER=runwayml
   export AI_API_KEY="your_runway_key"
   ```

---

## Testing Your Integration

1. **Start the backend** (if not running):
   ```bash
   cd backend
   go run main.go
   ```

2. **Upload test files** via the UI at http://localhost:3000

3. **Click "Generate Video"** and wait (2-10 minutes depending on service)

4. **Check the result** in the generated videos folder

---

## Customizing Video Content

### Edit Video Scripts

Edit the script in `backend/internal/services/video_generator.go`:

```go
// For D-ID (line ~80)
"input": "YOUR CUSTOM SCRIPT HERE - talk about the product features, benefits, and call to action",

// For Synthesia (line ~150)
"scriptText": "YOUR CUSTOM MARKETING MESSAGE",

// For RunwayML (line ~40)
"prompt": "YOUR CUSTOM VIDEO DESCRIPTION",
```

### Customize Video Style

**D-ID Options**:
- Voice: `en-US-JennyNeural`, `en-GB-RyanNeural`, etc.
- Provider: `microsoft`, `elevenlabs`, `amazon`

**Synthesia Options**:
- Avatar: `anna_costume1_cameraA`, `jack_costume3_cameraA`
- Background: `green_screen`, `office`, `modern`

**RunwayML Options**:
- Duration: 4-18 seconds
- Resolution: 1280x720 or 1920x1080
- Motion: Low to high

---

## Cost Comparison

| Service | Free Tier | Paid Plan | Cost per Video | Best For |
|---------|-----------|-----------|----------------|----------|
| D-ID | 20 credits free | $5.90/month | ~$0.30 | MVP, Testing |
| Synthesia | Demo only | $22/month | ~$0.70 | Professional |
| RunwayML | 125 credits | $12/month | ~$0.50 | Creative |

---

## Troubleshooting

### "Video generation failed"
- ✅ Check API key is correct
- ✅ Verify you have credits/quota remaining
- ✅ Check image file formats (JPG, PNG only)
- ✅ Ensure file sizes < 10MB

### "Connection timeout"
- ✅ Video generation can take 5-10 minutes
- ✅ Check internet connection
- ✅ Verify API endpoint URLs

### "Invalid API key"
- ✅ Make sure you copied the full key
- ✅ Check for extra spaces
- ✅ Regenerate key if needed

---

## Production Recommendations

### For MVP / Testing
- **Use**: D-ID
- **Why**: Cheapest, fastest setup, good quality
- **Cost**: ~$6/month for 60 videos

### For Business / Scale
- **Use**: Synthesia
- **Why**: Professional quality, enterprise features
- **Cost**: $22-67/month for 10-40 videos

### For Creative / Marketing
- **Use**: RunwayML
- **Why**: Unique, artistic style
- **Cost**: $12-35/month for variable videos

---

## Alternative Services

If the above don't work, try:

- **HeyGen**: https://www.heygen.com/ (AI avatar videos)
- **Pictory**: https://pictory.ai/ (Text to video)
- **Invideo AI**: https://invideo.io/ (Marketing videos)
- **Lumen5**: https://lumen5.com/ (Social media videos)

---

## Mock Mode (No API Key)

For testing without API keys:

```bash
export AI_PROVIDER=mock
# No API key needed
```

This creates placeholder videos for development.

---

## Next Steps

1. Choose your AI service
2. Sign up and get API key
3. Configure environment variables
4. Test video generation
5. Customize scripts and prompts
6. Launch your product!

## Support

- Check service documentation links above
- Review error logs in terminal
- Test with small files first
- Monitor API quotas/limits


