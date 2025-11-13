# 🚀 Full AI Video Generation Pipeline - Setup Guide

## Overview

Your backend now supports **TWO modes** of video generation:

### Mode 1: Standard (D-ID Only) ✅ Currently Working
- **What it does:** Generates a talking avatar video
- **Output:** Avatar speaking about the product
- **Cost:** ~$0.30 per video
- **Time:** ~30 seconds

### Mode 2: Full AI Pipeline (NEW!) 🎬
- **What it does:** 
  1. D-ID → Talking avatar video
  2. RunwayML → Animated product video
  3. Shotstack → Composites both (picture-in-picture)
- **Output:** Product showcase video with talking avatar overlay
- **Cost:** ~$0.65 per video  
- **Time:** ~2-3 minutes

---

## 🔧 Configuration

Add these to your `.env` file:

```bash
# Enable Full AI Pipeline
USE_FULL_AI_PIPELINE=true

# D-ID API Key (you already have this)
AI_API_KEY=your_did_key_here

# RunwayML API Key (NEW - get from https://app.runwayml.com/account/secrets)
RUNWAYML_API_KEY=your_runwayml_key_here

# Shotstack API Key (NEW - get from https://dashboard.shotstack.io/)
SHOTSTACK_API_KEY=your_shotstack_key_here
```

---

## 📝 Getting API Keys

### 1. D-ID (Already Have ✅)
- You're already using this
- Current key works perfectly

### 2. RunwayML API Key
**Sign up:** https://runwayml.com/
**Get Key:** https://app.runwayml.com/account/secrets

**Pricing:**
- Free tier: $5 credits (~20 videos)
- Paid: $12/month for 625 credits

**Steps:**
1. Sign up for RunwayML
2. Go to Account > API Secrets
3. Create new secret
4. Copy the API key
5. Add to `.env`: `RUNWAYML_API_KEY=rw_...`

### 3. Shotstack API Key
**Sign up:** https://shotstack.io/
**Get Key:** https://dashboard.shotstack.io/

**Pricing:**
- Free tier: 20 renders/month
- Paid: $29/month unlimited

**Steps:**
1. Sign up for Shotstack
2. Go to Dashboard
3. Copy your API key
4. Add to `.env`: `SHOTSTACK_API_KEY=...`

---

## 🎯 How It Works

### Standard Mode (Current)
```
User uploads:
  ├─ Product image
  └─ Presenter image

Backend:
  └─ D-ID API
      └─ Generates talking avatar

Output:
  └─ Avatar talking about product
```

### Full AI Pipeline Mode (NEW)
```
User uploads:
  ├─ Product image
  └─ Presenter image

Backend:
  ├─ Step 1: D-ID API
  │   └─ Generates talking avatar video
  │
  ├─ Step 2: RunwayML Gen-3 API
  │   └─ Generates animated product showcase
  │       (product rotating, zooming, cinematic)
  │
  └─ Step 3: Shotstack API
      └─ Composites both videos
          (product main, avatar in corner)

Output:
  └─ Professional marketing video:
      Product showcase + Talking presenter
```

---

## 💰 Cost Breakdown

### Per Video:
| Service | Cost | Purpose |
|---------|------|---------|
| D-ID | $0.30 | Talking avatar |
| RunwayML | $0.25 | Product animation |
| Shotstack | $0.10 | Video compositing |
| **TOTAL** | **$0.65** | **Full pipeline** |

### Monthly (100 videos):
- **Standard Mode:** $30/month (D-ID only)
- **Full AI Pipeline:** $65/month (all 3 AI services)

---

## 🧪 Testing

### Test Standard Mode (Currently Working)
```bash
# Make sure USE_FULL_AI_PIPELINE=false in .env
curl 'http://localhost:8080/api/v1/projects/YOUR_PROJECT_ID/generate-video' \
  -H 'Content-Type: application/json' \
  --data-raw '{"script":"Test video"}'
```

### Test Full AI Pipeline (Once API Keys Added)
```bash
# Set USE_FULL_AI_PIPELINE=true in .env
# Add RUNWAYML_API_KEY and SHOTSTACK_API_KEY to .env
# Restart backend

curl 'http://localhost:8080/api/v1/projects/YOUR_PROJECT_ID/generate-video' \
  -H 'Content-Type: application/json' \
  --data-raw '{"script":"Test full pipeline"}'
```

---

## 📊 Pipeline Progress Logs

When running full pipeline, you'll see:

```
🚀 ========================================
🚀 FULL AI PIPELINE STARTED
🚀 ========================================

📍 STEP 1/3: Generating Talking Avatar with D-ID
📸 Uploading presenter image to D-ID...
✅ D-ID talk created: tlk_xxx
⏳ Waiting for avatar generation...
✅ STEP 1 COMPLETE: Avatar video saved

📍 STEP 2/3: Generating Product Video with RunwayML Gen-3
🎬 Generating product video with RunwayML Gen-2...
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

## 🎬 Example Output

### Standard Mode Output:
```
┌─────────────────┐
│                 │
│      👤         │  ← Talking avatar
│     Person      │
│                 │
└─────────────────┘
```

### Full AI Pipeline Output:
```
┌─────────────────────────┐
│                         │
│   [Product Showcase]    │  ← Animated product video
│   (rotating, zooming)   │     (RunwayML generated)
│                         │
│              ┌────┐     │
│              │👤  │     │  ← Talking avatar overlay
│              └────┘     │     (D-ID generated)
└─────────────────────────┘
         ↑
  (Shotstack composited)
```

---

## 🔄 Switching Modes

### To Use Standard Mode (Avatar Only):
```bash
# In .env
USE_FULL_AI_PIPELINE=false
AI_PROVIDER=did
```

### To Use Full AI Pipeline:
```bash
# In .env
USE_FULL_AI_PIPELINE=true
AI_PROVIDER=did
RUNWAYML_API_KEY=your_key
SHOTSTACK_API_KEY=your_key
```

---

## ⚙️ Configuration File

Your `.env` should look like:

```bash
# Backend
PORT=8080
DATABASE_PATH=./data/app.db
UPLOAD_PATH=./uploads
GENERATED_VIDEO_PATH=./generated/videos
WEBSITE_PATH=./generated/websites

# AI Provider
AI_PROVIDER=did

# Full Pipeline Toggle
USE_FULL_AI_PIPELINE=true

# API Keys
AI_API_KEY=cmFrZXNoZGQ0NDU0QGdtYWlsLmNvbQ:DEGE6f5zBPjimAmsqg0oL
RUNWAYML_API_KEY=your_runwayml_key
SHOTSTACK_API_KEY=your_shotstack_key
```

---

## 🚨 Important Notes

1. **RunwayML Gen-3** takes ~60-90 seconds to generate product video
2. **Shotstack** takes ~30-60 seconds to composite
3. **Total pipeline time:** 2-3 minutes per video
4. **Standard mode** takes only ~30 seconds
5. **API keys must be valid** or the pipeline will fail at that step

---

## ✅ Implementation Status

- ✅ D-ID integration (working)
- ✅ RunwayML Gen-3 integration (implemented)
- ✅ Shotstack compositing (implemented)
- ✅ Full pipeline orchestration (implemented)
- ✅ Configuration toggle (implemented)
- ⏳ **Needs:** API keys for RunwayML and Shotstack

---

## 🎯 Next Steps

1. **Get RunwayML API key** → https://app.runwayml.com/account/secrets
2. **Get Shotstack API key** → https://dashboard.shotstack.io/
3. **Add keys to `.env` file**
4. **Set `USE_FULL_AI_PIPELINE=true`**
5. **Restart backend**
6. **Test video generation!**

---

**Questions? The full pipeline is ready to use once you add the API keys!** 🚀

