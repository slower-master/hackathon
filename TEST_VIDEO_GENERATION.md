# Test Video Generation with D-ID

## 🎥 Your System is Ready!

You now have **D-ID API integrated** and ready to generate professional marketing videos!

## ✅ What's New

### 1. Text Input for Video Scripts
- Optional text area in the UI
- Enter your own marketing message
- Or leave empty for AI-generated script

### 2. Smart Prompt Enhancement
- Automatically improves short prompts
- Adds professional greeting if missing
- Adds call-to-action if missing
- Optimizes length for 30-second videos

### 3. Professional Voice
- Using Microsoft Neural Voice (Jenny)
- Clear, professional narration
- Natural-sounding speech

---

## 🎬 How to Generate Your First Video

### Step 1: Open the App
Go to: http://localhost:3000

### Step 2: Upload Files
1. **Product Image**: Upload a clear product photo (JPG/PNG)
2. **Person Photo**: Upload a headshot photo of a person

### Step 3: Write Your Script (Optional)

**Option A: Let AI Write It**
- Leave the text box empty
- AI will generate a professional marketing script

**Option B: Write Your Own**
Try one of these examples:

**Simple Product Introduction:**
```
This product is amazing! It helps you save time and get more done.
```
*AI will enhance it to:*
"Hello! This product is amazing! It helps you save time and get more done. Get started today and experience the difference!"

**Product Features:**
```
Our product has three key features: fast performance, easy to use, and affordable price
```

**Benefits Focused:**
```
Imagine saving 5 hours every week. That's what our product delivers. More time for what matters.
```

**Call to Action:**
```
Don't wait! Transform your workflow today with our innovative solution
```

### Step 4: Generate!
1. Click "Upload & Create Project"
2. Wait for upload to complete
3. Click "Generate Video"
4. **Wait 2-5 minutes** (D-ID generates the video)
5. Click "View Video" when complete!

---

## 💡 Script Writing Tips

### For Best Results:

**DO:**
- ✅ Mention product name or type
- ✅ List key benefits (not just features)
- ✅ Use emotional language
- ✅ Keep it 50-100 words (30 seconds)
- ✅ End with a call to action

**DON'T:**
- ❌ Write super long paragraphs (>150 words)
- ❌ Use technical jargon
- ❌ Forget to mention benefits
- ❌ Make it boring or robotic

### Great Examples:

**Example 1 - Software Product:**
```
Tired of wasting hours on repetitive tasks? Our automation tool does the work for you. 
Set it up once, and watch it handle everything automatically. Join 10,000+ happy users 
who've reclaimed their time. Try it free today!
```

**Example 2 - Physical Product:**
```
Meet the SmartBottle - your hydration companion. It tracks your water intake, reminds 
you to drink, and keeps your water cold for 24 hours. Perfect for busy professionals 
and fitness enthusiasts. Order yours now and stay healthy!
```

**Example 3 - Service:**
```
Growing your business shouldn't be complicated. We handle your marketing, so you can 
focus on what you do best. From social media to email campaigns, we've got you covered. 
Schedule your free consultation today!
```

---

## 🔍 What Happens Behind the Scenes

1. **Upload**: Files saved to server
2. **Script Enhancement**: Your text is professionally enhanced
3. **Length Optimization**: Script adjusted to ~30 seconds
4. **D-ID API Call**: Sends person photo + script
5. **AI Processing**: D-ID creates talking head video (2-5 min)
6. **Download**: Video saved to server
7. **Complete**: Ready to view!

---

## 🎯 Testing Checklist

### Test 1: AI-Generated Script
- [ ] Upload product image
- [ ] Upload person photo
- [ ] Leave script box empty
- [ ] Generate video
- [ ] Check result

### Test 2: Simple User Input
- [ ] Upload files
- [ ] Write: "This product is great!"
- [ ] Generate video
- [ ] Verify AI enhanced it

### Test 3: Full Custom Script
- [ ] Upload files
- [ ] Write detailed marketing message (50-100 words)
- [ ] Generate video
- [ ] Check quality

### Test 4: Long Input
- [ ] Upload files
- [ ] Write very long text (200+ words)
- [ ] Generate video
- [ ] Verify it was optimized to 30 seconds

---

## 🐛 Troubleshooting

### "Failed to generate video"

**Check 1: API Key**
```bash
cd backend
# Verify API key is set
echo $AI_API_KEY
# Should show: ZGluZXNoLnJhbUBkZWFsc2hhcmUuaW46:31Ln-cQ0VJ-_sj61DbZm3
```

**Check 2: D-ID Credits**
- Login to https://studio.d-id.com/
- Check remaining credits
- Need at least 1 credit per video

**Check 3: Image Format**
- Must be JPG or PNG
- Person photo should show a clear face
- File size < 5MB

### "Video generation timeout"

This is normal! D-ID takes 2-5 minutes:
- Don't close browser
- Don't refresh page
- Wait patiently
- Check terminal logs for progress

### "Invalid image format"

D-ID requires:
- Clear face photo
- Front-facing
- Good lighting
- JPG or PNG format

---

## 📊 Expected Results

### Video Quality:
- ✅ 720p HD quality
- ✅ 30 seconds duration
- ✅ Professional voice
- ✅ Smooth lip-sync
- ✅ Natural expressions

### Processing Time:
- Upload: < 1 second
- Video Generation: 2-5 minutes
- Download: 5-10 seconds
- **Total: ~3-6 minutes**

---

## 🎉 Next Steps After First Video

1. **Generate Website**
   - Click "Generate Website"
   - Professional marketing page created
   - Includes your video

2. **Try Different Scripts**
   - Test various tones
   - Compare results
   - Find what works best

3. **Experiment with Voices**
   - Edit `backend/internal/services/video_generator.go`
   - Change voice_id (line 122)
   - Options: en-US-JennyNeural, en-US-GuyNeural, etc.

4. **Customize Templates**
   - Modify website templates
   - Add your branding
   - Customize colors

---

## 📈 Cost Tracking

Each video generation uses:
- **1 D-ID credit** = ~$0.30
- **Your account**: Check at https://studio.d-id.com/

Free tier: 20 credits = 20 videos

---

## 🎬 Ready to Create!

Your system is configured and ready. Open http://localhost:3000 and create your first professional marketing video!

**Tips for Success:**
1. Use clear, well-lit person photos
2. Write benefit-focused scripts
3. Keep it concise (50-100 words)
4. Be patient during generation
5. Test different approaches

Happy creating! 🚀

