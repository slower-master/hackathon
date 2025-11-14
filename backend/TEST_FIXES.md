# 🧪 Test Your Fixes

## ✅ Fix #1: Person Video Visibility (Z-Index Fix)

### What Changed:
**REVERSED the track order** - Person video is now Track 0 (rendered on top)

### Before:
```
Track 0: Product video (background)
Track 1: Person video (overlay)
Result: Person hidden behind product ❌
```

### After:
```
Track 0: Person video (TOP LAYER)
Track 1: Product video (BACKGROUND)
Result: Person visible on top ✅
```

### Test It:
1. Restart backend: `./start.sh`
2. Generate a new video
3. Check logs for: `🎯 REVERSED TRACK ORDER: Track 0 = Person (TOP)`
4. Watch the final video - person should be VISIBLE in bottom-right

---

## ✅ Fix #2: Gemini Product-Specific Features

### Test Script:
```bash
# Set your Gemini API key
export GOOGLE_GEMINI_API_KEY="your-key-here"

# Run the test
cd backend
./test_gemini.sh
```

### Expected Output:
```
✅ Gemini API Key Found: AIzaSyD...

------------------------------------------------------
📦 Product: Namkeen papdi
📝 Description: Tasty Healthy Namkeen Packet
✅ Generated Features:

  🥘 Deliciously Wholesome
     Savor the irresistible crunch and authentic flavor...

  💪 Guilt-Free Goodness
     Enjoy your favorite snack without compromise...

  ⚡ Anytime, Anywhere Snack
     Perfectly portioned and easy to carry...

  🌟 Smart Snack Choice
     Get big flavor at an unbeatable price...
```

### If You See Errors:
1. **"No Gemini API key found"**
   - Set: `export GOOGLE_GEMINI_API_KEY="your-key"`
   - Get key: https://makersuite.google.com/app/apikey

2. **"API quota exceeded"**
   - Wait a minute (free tier: 60 requests/min)
   - Or get a new API key

3. **"Network error"**
   - Check internet connection
   - Check firewall settings

---

## 🎯 Complete Testing Workflow

### Step 1: Test Gemini (5 minutes)
```bash
export GOOGLE_GEMINI_API_KEY="your-key"
cd backend
./test_gemini.sh
```

✅ You should see **product-specific** features, not generic ones

### Step 2: Start Backend
```bash
./start.sh
```

✅ Look for in logs:
- `🔑 Gemini API Key: Found (39 chars)`
- `🎯 REVERSED TRACK ORDER: Track 0 = Person (TOP)`

### Step 3: Upload Product & Generate
1. Go to: http://localhost:8080
2. Upload: Product image + Person image/video
3. Generate video and website

### Step 4: Verify Results

**Website Features:**
✅ Should show product-specific features (not "Lightning Fast", "Premium Quality")

**Video:**
✅ Person should be VISIBLE in bottom-right corner
✅ Person audio should play with person video visible

---

## 📊 Quick Checklist

### Website Features (Gemini)
- [ ] Gemini API key set in environment
- [ ] Test script shows product-specific features
- [ ] Backend logs show "Gemini API Key: Found"
- [ ] Website shows product-relevant features

### Person Video (Z-Index)
- [ ] Backend rebuilt with reversed track order
- [ ] Logs show "REVERSED TRACK ORDER"
- [ ] Final video shows person in bottom-right
- [ ] Person is VISIBLE (not just audio)

---

## 🆘 Still Having Issues?

### Person Video Not Visible
Try these track orders in `video_generator_simple.go`:

**Option 1 (Current):** Person first, Product second
**Option 2:** Add explicit z-index values
**Option 3:** Use Shotstack "position" property

### Gemini Not Working
Check in this order:
1. API key is set: `echo $GOOGLE_GEMINI_API_KEY`
2. API key is valid (test with script)
3. Backend restarted after setting key
4. Check backend logs for errors

---

## 📝 Files Created
- `test_gemini.sh` - Test Gemini API
- `start.sh` - Start backend with environment
- `SETUP_GEMINI.md` - Complete Gemini setup guide
- `TEST_FIXES.md` - This file

---

**Ready to test! Run `./test_gemini.sh` first, then `./start.sh`** 🚀

