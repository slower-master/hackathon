#!/bin/bash

# 🧪 Test Video Generation Script
# Tests the full AI pipeline with different options

echo "🧪 ========================================"
echo "🧪 Video Generation Test Script"
echo "🧪 ========================================"
echo ""

# Check if backend is running
if ! lsof -ti:8080 > /dev/null 2>&1; then
    echo "❌ Backend is not running on port 8080!"
    echo "Please start the backend first using:"
    echo "  ./START_FULL_PIPELINE.sh"
    echo ""
    exit 1
fi

echo "✅ Backend is running"
echo ""

# Get project ID (you'll need to replace this with actual project ID)
echo "📋 You need a project ID to test video generation"
echo "   Get it by uploading images via the frontend first"
echo ""
echo "Enter your project ID (or press Enter to use test ID):"
read PROJECT_ID

if [ -z "$PROJECT_ID" ]; then
    PROJECT_ID="5f96fe9f-19be-45db-b859-efe30e356cec"  # Default test ID
    echo "Using test project ID: $PROJECT_ID"
fi

echo ""
echo "🎬 ========================================"
echo "🎬 Test 1: Standard Mode (D-ID Only)"
echo "🎬 ========================================"
echo ""

curl -X POST "http://localhost:8080/api/v1/projects/$PROJECT_ID/generate-video" \
  -H "Content-Type: application/json" \
  -d '{
    "script": "Test video with standard D-ID mode"
  }' 2>&1

echo ""
echo ""
echo "⏳ Waiting 5 seconds before next test..."
sleep 5
echo ""

# Check if USE_FULL_AI_PIPELINE is enabled
if grep -q "USE_FULL_AI_PIPELINE=true" .env 2>/dev/null; then
    echo "🎬 ========================================"
    echo "🎬 Test 2: Full AI Pipeline - Product Focus"
    echo "🎬 ========================================"
    echo ""
    
    curl -X POST "http://localhost:8080/api/v1/projects/$PROJECT_ID/generate-video" \
      -H "Content-Type: application/json" \
      -d '{
        "script": "Introducing our revolutionary product with advanced technology!",
        "product_video_style": "rotation",
        "layout": "product_main"
      }' 2>&1
    
    echo ""
    echo ""
    echo "⏳ Waiting 5 seconds before next test..."
    sleep 5
    echo ""
    
    echo "🎬 ========================================"
    echo "🎬 Test 3: Full AI Pipeline - Avatar Focus"
    echo "🎬 ========================================"
    echo ""
    
    curl -X POST "http://localhost:8080/api/v1/projects/$PROJECT_ID/generate-video" \
      -H "Content-Type: application/json" \
      -d '{
        "script": "Let me show you this amazing product that will change your life!",
        "product_video_style": "zoom",
        "layout": "avatar_main"
      }' 2>&1
    
    echo ""
    echo ""
else
    echo "ℹ️  Full AI Pipeline is disabled (USE_FULL_AI_PIPELINE=false)"
    echo "   Only testing standard D-ID mode"
    echo ""
fi

echo ""
echo "✅ ========================================"
echo "✅ Tests Complete!"
echo "✅ ========================================"
echo ""
echo "Check backend logs for detailed progress"
echo "Videos will be saved in: generated/videos/"
echo ""

