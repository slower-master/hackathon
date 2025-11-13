# 🚀 Quick Start Guide - Full AI Video Pipeline

## ✅ What's Been Implemented

Your backend now has a **complete AI-powered video generation pipeline**:

1. ✅ **D-ID Integration** - Talking avatar generation
2. ✅ **RunwayML Gen-3 Integration** - Animated product videos  
3. ✅ **Shotstack Integration** - Professional video compositing
4. ✅ **Configurable Options** - Product styles & layouts
5. ✅ **Complete Pipeline Orchestration** - Automated workflow

---

## 🎯 Quick Test (5 Minutes)

### Step 1: Update .env File

```bash
cd /Users/slowermaster/DEALSHARE/hacathon/backend

# Copy template
cp ENV_TEMPLATE.txt .env

# Edit .env and add your API keys:
nano .env
```

**Required in .env:**
```bash
# Basic (working already)
AI_API_KEY=your_did_key

# For full pipeline (add these)
USE_FULL_AI_PIPELINE=true
RUNWAYML_API_KEY=your_runwayml_key
SHOTSTACK_API_KEY=your_shotstack_key
```

### Step 2: Start Backend

```bash
# Use the startup script
./START_FULL_PIPELINE.sh

# Or manually:
go run main.go
```

### Step 3: Test Video Generation

```bash
# In a new terminal, run:
./TEST_VIDEO_GENERATION.sh

# Or manually test:
curl -X POST 'http://localhost:8080/api/v1/projects/YOUR_PROJECT_ID/generate-video' \
  -H 'Content-Type: application/json' \
  -d '{
    "script": "Test my amazing product!",
    "product_video_style": "rotation",
    "layout": "product_main"
  }'
```

---

## 📋 Configuration Options

### Product Video Styles:
- `"rotation"` - 360° rotation (best for electronics)
- `"zoom"` - Zoom into details (best for intricate items)
- `"pan"` - Camera pans around (best for large products)
- `"reveal"` - Dramatic reveal (best for luxury)
- `"auto"` - AI chooses automatically (default)

### Layouts:
- `"product_main"` - Product fullscreen, avatar overlay (default)
- `"avatar_main"` - Avatar fullscreen, product overlay

---

## 🧪 Test Commands

### Test 1: Standard Mode (D-ID Only)
```bash
curl -X POST 'http://localhost:8080/api/v1/projects/PROJECT_ID/generate-video' \
  -H 'Content-Type: application/json' \
  -d '{"script": "Test standard mode"}'
```
**Expected:** ~30 seconds, avatar video

### Test 2: Full Pipeline - Product Focus
```bash
curl -X POST 'http://localhost:8080/api/v1/projects/PROJECT_ID/generate-video' \
  -H 'Content-Type: application/json' \
  -d '{
    "script": "Amazing product showcase!",
    "product_video_style": "rotation",
    "layout": "product_main"
  }'
```
**Expected:** ~2-3 minutes, product video with avatar overlay

### Test 3: Full Pipeline - Avatar Focus
```bash
curl -X POST 'http://localhost:8080/api/v1/projects/PROJECT_ID/generate-video' \
  -H 'Content-Type: application/json' \
  -d '{
    "script": "Let me show you this product!",
    "product_video_style": "zoom",
    "layout": "avatar_main"
  }'
```
**Expected:** ~2-3 minutes, avatar video with product overlay

---

## 📊 Expected Output Logs

When running full pipeline, you'll see:

```
🚀 ========================================
🚀 FULL AI PIPELINE STARTED
🚀 ========================================

📋 Configuration:
   Product Video Style: rotation
   Layout: product_main

📍 STEP 1/3: Generating Talking Avatar with D-ID
📸 Uploading presenter image to D-ID...
✅ D-ID talk created: tlk_xxx
⏳ Waiting for avatar generation...
✅ STEP 1 COMPLETE: Avatar video saved

📍 STEP 2/3: Generating Product Video with RunwayML Gen-3
🎬 Generating product video with RunwayML Gen-3...
📤 Calling RunwayML API...
✅ RunwayML task created: task_xxx
⏳ Waiting for product video generation...
✅ STEP 2 COMPLETE: Product video saved

📍 STEP 3/3: Compositing Videos with Shotstack
🎨 Compositing videos with Shotstack API...
📤 Calling Shotstack API...
✅ Shotstack render started: render_xxx
⏳ Waiting for video compositing...
✅ STEP 3 COMPLETE: Final video saved

🎉 ========================================
🎉 FULL AI PIPELINE COMPLETED SUCCESSFULLY!
🎉 ========================================
```

---

## 🔧 Troubleshooting

### Backend won't start
```bash
# Kill existing processes
lsof -ti:8080 | xargs kill -9
pkill -f "go run main.go"

# Restart
./START_FULL_PIPELINE.sh
```

### 401 Unauthorized (D-ID)
- Check AI_API_KEY in .env is correct
- Verify key format (should be Base64 encoded)

### 401 Unauthorized (RunwayML)
- Check RUNWAYML_API_KEY in .env
- Get fresh key from https://app.runwayml.com/account/secrets

### 401 Unauthorized (Shotstack)
- Check SHOTSTACK_API_KEY in .env
- Get key from https://dashboard.shotstack.io/

### Pipeline takes too long
- Normal: 2-3 minutes for full pipeline
- RunwayML: ~60-90 seconds (product video generation)
- Shotstack: ~30-60 seconds (compositing)

---

## 📁 File Structure

```
backend/
├── START_FULL_PIPELINE.sh     # Start backend with full pipeline
├── TEST_VIDEO_GENERATION.sh   # Test video generation
├── ENV_TEMPLATE.txt            # Environment template
├── .env                        # Your configuration (create this)
├── main.go                     # Backend entry point
├── internal/
│   ├── config/config.go       # Configuration management
│   ├── handlers/handlers.go   # API endpoints
│   └── services/
│       ├── ai_service.go      # AI service router
│       └── video_generator.go # Full pipeline implementation
└── generated/videos/          # Generated videos saved here
```

---

## 🎬 Video Output

Generated videos will be saved in:
```
backend/generated/videos/[uuid].mp4
```

Access via API response:
```json
{
  "project_id": "xxx",
  "status": "video_complete",
  "video_path": "generated/videos/xxx.mp4"
}
```

View in browser:
```
http://localhost:8080/static/generated/videos/xxx.mp4
```

---

## 💰 Cost Tracking

### Per Video Costs:
- **Standard Mode:** $0.30 (D-ID only)
- **Full Pipeline:** $0.65 (D-ID + RunwayML + Shotstack)

### Monthly Estimates:
- 50 videos (standard): $15/month
- 50 videos (full): $32.50/month
- 100 videos (standard): $30/month
- 100 videos (full): $65/month

---

## 📚 Documentation

All documentation is in the `backend/` directory:

- `AI_PIPELINE_SETUP.md` - Complete setup guide
- `IMPLEMENTATION_SUMMARY.md` - Technical details
- `VIDEO_OPTIONS_GUIDE.md` - Video style options
- `COMPOSITOR_EXPLAINED.md` - Shotstack explanation
- `FRONTEND_UI_GUIDE.md` - UI implementation guide

---

## ✅ Checklist

Before testing:
- [ ] .env file created with API keys
- [ ] USE_FULL_AI_PIPELINE=true (for full pipeline)
- [ ] Backend compiles (run: `go build`)
- [ ] Backend starts (run: `./START_FULL_PIPELINE.sh`)
- [ ] Port 8080 is available
- [ ] Have a valid project ID (upload images via frontend first)

---

## 🚀 Next Steps

1. **Test standard mode first** (just D-ID)
2. **Add RunwayML & Shotstack keys**
3. **Test full pipeline**
4. **Add frontend UI controls** (see FRONTEND_UI_GUIDE.md)
5. **Go to production!**

---

**Questions? Check the documentation files or review backend logs for detailed error messages.**

🎉 **You're ready to generate amazing AI-powered marketing videos!**
