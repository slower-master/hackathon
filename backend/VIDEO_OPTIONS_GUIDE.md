# 🎬 Video Generation Options Guide

## Overview

Your video generation API now supports **customizable product video styles** and **flexible layouts**!

---

## 📹 Product Video Styles

You can control **how the product is animated** in the video by choosing a style:

### Available Styles:

#### 1. **`"rotation"`** 🔄 (Recommended for most products)
- **What it does:** Product rotates 360 degrees smoothly
- **Best for:** Electronics, gadgets, 3D objects, items that look good from all angles
- **Example:** iPhone rotating to show all sides
- **Prompt used:** *"Professional product showcase with smooth 360-degree rotation, studio lighting, elegant spin, premium commercial feel"*

#### 2. **`"zoom"`** 🔍
- **What it does:** Camera zooms in from wide shot to close-up
- **Best for:** Products with fine details, textures, logos, intricate designs
- **Example:** Watch face zooming in to show details
- **Prompt used:** *"Professional product showcase with smooth zoom-in effect, starting wide and focusing on product details, studio lighting"*

#### 3. **`"pan"`** 📷
- **What it does:** Camera pans around the product, exploring different angles
- **Best for:** Large products, furniture, items with multiple interesting angles
- **Example:** Laptop being viewed from different sides
- **Prompt used:** *"Professional product showcase with smooth camera pan movement, exploring product from different angles, studio lighting"*

#### 4. **`"reveal"`** ✨
- **What it does:** Dramatic reveal with lighting effects, product emerging from shadows
- **Best for:** Premium products, luxury items, dramatic presentations
- **Example:** Jewelry emerging from darkness with dramatic lighting
- **Prompt used:** *"Professional product reveal with dramatic lighting, product emerging from shadows, cinematic reveal, premium commercial feel"*

#### 5. **`"auto"`** 🤖 (Default)
- **What it does:** Automatically chooses the best style (currently defaults to rotation)
- **Best for:** When you're not sure which style to use
- **Example:** Generic product showcase with smooth movement
- **Prompt used:** *"Professional product showcase with smooth camera movement, elegant rotation, studio lighting"*

---

## 📐 Layout Options

You can control **which element is the main focus** and which is the overlay:

### Available Layouts:

#### 1. **`"product_main"`** 📦 (Default)
- **Layout:**
  ```
  ┌─────────────────────────┐
  │                         │
  │   [Product Video]       │  ← Fullscreen (main)
  │   (rotating/zooming)    │
  │                         │
  │              ┌────┐     │
  │              │👤  │     │  ← Overlay (25% size)
  │              └────┘     │     (bottom-right)
  └─────────────────────────┘
  ```
- **Best for:** Product-focused marketing, showcasing product features
- **Use case:** When product is the star, presenter is supporting

#### 2. **`"avatar_main"`** 👤
- **Layout:**
  ```
  ┌─────────────────────────┐
  │                         │
  │      [Avatar Video]     │  ← Fullscreen (main)
  │      (talking)          │
  │                         │
  │              ┌────┐     │
  │              │📦  │     │  ← Overlay (30% size)
  │              └────┘     │     (bottom-right)
  └─────────────────────────┘
  ```
- **Best for:** Personal branding, influencer-style videos, presenter-focused content
- **Use case:** When presenter is the star, product is supporting

---

## 🚀 API Usage

### Request Format:

```json
POST /api/v1/projects/:id/generate-video
Content-Type: application/json

{
  "script": "Your marketing script here",
  "product_video_style": "rotation",  // Optional: "rotation", "zoom", "pan", "reveal", "auto"
  "layout": "product_main"            // Optional: "product_main", "avatar_main"
}
```

### Examples:

#### Example 1: Product-focused with rotation
```json
{
  "script": "Introducing our revolutionary smartphone!",
  "product_video_style": "rotation",
  "layout": "product_main"
}
```
**Result:** Product rotating fullscreen, avatar in corner

#### Example 2: Presenter-focused with zoom
```json
{
  "script": "Let me show you this amazing product!",
  "product_video_style": "zoom",
  "layout": "avatar_main"
}
```
**Result:** Presenter fullscreen, product zooming in corner

#### Example 3: Auto-detect everything
```json
{
  "script": "Check out our new product!"
}
```
**Result:** Uses defaults: `"auto"` style, `"product_main"` layout

---

## 🎯 Decision Guide

### How to Choose Product Video Style:

| Product Type | Recommended Style | Why |
|--------------|-------------------|-----|
| **Electronics** (phones, laptops) | `rotation` | Shows all sides, ports, features |
| **Jewelry/Watches** | `zoom` or `reveal` | Highlights details, premium feel |
| **Furniture** | `pan` | Shows size, different angles |
| **Clothing** | `pan` or `rotation` | Shows fit, style from all angles |
| **Food/Beverages** | `zoom` or `reveal` | Appetizing, highlights freshness |
| **Cosmetics** | `zoom` | Shows texture, color details |
| **Not sure?** | `auto` | Let AI decide |

### How to Choose Layout:

| Marketing Goal | Recommended Layout | Why |
|----------------|-------------------|-----|
| **Product showcase** | `product_main` | Product is the hero |
| **Personal branding** | `avatar_main` | Presenter is the hero |
| **Influencer marketing** | `avatar_main` | Trust-building, personal connection |
| **E-commerce** | `product_main` | Focus on product features |
| **Testimonials** | `avatar_main` | Person's story is important |

---

## 📊 Visual Comparison

### Layout: `product_main` (Default)
```
┌─────────────────────────┐
│                         │
│   📦 PRODUCT            │  ← 100% screen
│   (Animated showcase)   │
│                         │
│              ┌────┐     │
│              │👤  │     │  ← 25% screen
│              └────┘     │
└─────────────────────────┘
```

### Layout: `avatar_main`
```
┌─────────────────────────┐
│                         │
│      👤 PRESENTER        │  ← 100% screen
│      (Talking)           │
│                         │
│              ┌────┐     │
│              │📦  │     │  ← 30% screen
│              └────┘     │
└─────────────────────────┘
```

---

## 💡 Pro Tips

1. **Start with defaults:** Use `"auto"` and `"product_main"` first, then customize based on results

2. **Match style to product:** 
   - Round objects → `rotation`
   - Flat objects → `pan` or `zoom`
   - Premium items → `reveal`

3. **Consider your audience:**
   - B2B → `product_main` (focus on features)
   - B2C → `avatar_main` (personal connection)

4. **Test different combinations:**
   - Try `rotation` + `product_main` for electronics
   - Try `reveal` + `avatar_main` for luxury items

---

## 🔄 Defaults

If you don't specify options:
- **Product Video Style:** `"auto"` (uses rotation)
- **Layout:** `"product_main"` (product fullscreen, avatar overlay)

---

## 📝 Summary

✅ **Product Video Styles:** Control how product animates (rotation, zoom, pan, reveal, auto)  
✅ **Layout Options:** Control which element is main vs overlay  
✅ **Fully Configurable:** All options via API request  
✅ **Smart Defaults:** Works great out of the box  

**You now have full control over your video generation!** 🎬

