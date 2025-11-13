# 🎨 Frontend UI Implementation Guide

## Current State
✅ Backend API supports product video styles and layouts  
❌ Frontend UI doesn't have controls for these options yet

---

## What You Need to Add to Frontend

### 1. Add UI Controls for Video Options

When generating a video, users should see these options:

```javascript
// Example UI structure needed in your video generation form:

<form onSubmit={handleGenerateVideo}>
  {/* Existing script input */}
  <textarea 
    name="script" 
    placeholder="Enter your marketing script..."
  />

  {/* NEW: Product Video Style Selector */}
  <div className="form-group">
    <label>Product Animation Style</label>
    <select name="product_video_style">
      <option value="auto">Auto (Let AI decide)</option>
      <option value="rotation">360° Rotation (Best for gadgets)</option>
      <option value="zoom">Zoom In (Best for details)</option>
      <option value="pan">Pan Around (Best for large items)</option>
      <option value="reveal">Dramatic Reveal (Best for luxury)</option>
    </select>
    <small>How should your product be animated?</small>
  </div>

  {/* NEW: Layout Selector */}
  <div className="form-group">
    <label>Video Layout</label>
    <select name="layout">
      <option value="product_main">Product Focus (Product fullscreen + Avatar in corner)</option>
      <option value="avatar_main">Presenter Focus (Avatar fullscreen + Product in corner)</option>
    </select>
    <small>Which element should be the main focus?</small>
  </div>

  <button type="submit">Generate Video</button>
</form>
```

---

### 2. Update API Call

Modify your video generation API call to include these options:

```javascript
// Example: Update your API call function

async function generateVideo(projectId, formData) {
  const response = await fetch(`/api/v1/projects/${projectId}/generate-video`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      script: formData.script,
      product_video_style: formData.product_video_style || 'auto',  // NEW
      layout: formData.layout || 'product_main',                     // NEW
    }),
  });

  return response.json();
}
```

---

### 3. Visual Previews (Optional but Recommended)

Add visual hints to help users understand the options:

```jsx
// Product Style Preview Icons
const styleIcons = {
  rotation: '🔄',
  zoom: '🔍',
  pan: '📷',
  reveal: '✨',
  auto: '🤖'
};

// Layout Preview
<div className="layout-preview">
  {layout === 'product_main' ? (
    <div className="preview-box">
      <div className="main-frame">Product Video</div>
      <div className="overlay-frame">Avatar</div>
    </div>
  ) : (
    <div className="preview-box">
      <div className="main-frame">Avatar</div>
      <div className="overlay-frame">Product</div>
    </div>
  )}
</div>
```

---

### 4. Smart Defaults Based on Use Case

You can provide presets for common scenarios:

```javascript
const videoPresets = {
  ecommerce: {
    product_video_style: 'rotation',
    layout: 'product_main',
    description: 'Perfect for product showcases'
  },
  influencer: {
    product_video_style: 'zoom',
    layout: 'avatar_main',
    description: 'Personal, engaging style'
  },
  luxury: {
    product_video_style: 'reveal',
    layout: 'product_main',
    description: 'Premium, dramatic feel'
  }
};

// UI for presets
<div className="presets">
  {Object.entries(videoPresets).map(([key, preset]) => (
    <button onClick={() => applyPreset(preset)}>
      {key} - {preset.description}
    </button>
  ))}
</div>
```

---

## Quick Implementation (5 Minutes)

### Minimal UI Addition:

Add these two dropdowns to your video generation form:

```html
<!-- Add before the "Generate Video" button -->

<label for="videoStyle">Product Animation:</label>
<select id="videoStyle" name="product_video_style">
  <option value="auto">Auto-detect (Recommended)</option>
  <option value="rotation">360° Rotation</option>
  <option value="zoom">Zoom In</option>
  <option value="pan">Pan Around</option>
  <option value="reveal">Dramatic Reveal</option>
</select>

<label for="layout">Layout:</label>
<select id="layout" name="layout">
  <option value="product_main">Product Main (Default)</option>
  <option value="avatar_main">Presenter Main</option>
</select>
```

Then update your JavaScript:

```javascript
// In your form submit handler:
const formData = {
  script: scriptInput.value,
  product_video_style: document.getElementById('videoStyle').value,
  layout: document.getElementById('layout').value
};

// Send to API
fetch(`/api/v1/projects/${projectId}/generate-video`, {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify(formData)
});
```

---

## Alternative: Use Defaults for Now

If you don't want to add UI immediately, the API will use sensible defaults:

```javascript
// This still works (uses default options):
fetch(`/api/v1/projects/${projectId}/generate-video`, {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    script: "Your script"
    // product_video_style: defaults to "auto"
    // layout: defaults to "product_main"
  })
});
```

**Defaults:**
- `product_video_style`: `"auto"` (AI chooses best style)
- `layout`: `"product_main"` (product fullscreen, avatar overlay)

---

## Recommended UI Flow

```
1. User uploads images ✅ (already implemented)
   ↓
2. User enters script ✅ (already implemented)
   ↓
3. User chooses options: ⚠️ (need to add)
   - Product animation style (dropdown)
   - Video layout (dropdown)
   ↓
4. Click "Generate Video" ✅ (already implemented)
   ↓
5. Backend generates with chosen options ✅ (already implemented)
```

---

## Summary

**Current:** Options via API only  
**Needed:** Add 2 dropdown selects to frontend form  
**Time:** ~5 minutes to implement basic version  
**Alternative:** Use defaults (no UI changes needed)

---

## Next Steps

1. **Quick fix:** Just use defaults (no frontend changes)
2. **Better UX:** Add 2 dropdowns (5 min implementation)
3. **Best UX:** Add presets + visual previews (30 min implementation)

Choose based on your timeline! 🚀

