#!/bin/bash

# Backend Startup Script with Environment Variables
# Make sure to set your API keys before running

# ==============================================
# ⚙️  CONFIGURATION - SET YOUR API KEYS HERE
# ==============================================

# Gemini API Key (for website features generation)
# Get your key from: https://makersuite.google.com/app/apikey
export GOOGLE_GEMINI_API_KEY="AIzaSyC_gI30tRdg-eYjVJn7ses22lrzrRB4vXc"

# D-ID API Key (for avatar video generation)
export DID_API_KEY="your-did-api-key-here"

# Shotstack API Key (for video merging)
export SHOTSTACK_API_KEY="your-shotstack-api-key-here"

# ==============================================
# 🎨 WEBSITE GENERATION SETTINGS
# ==============================================

# Use modern v0.dev style for websites (true/false)
export USE_V0_STYLE="true"

# ==============================================
# 🎬 VIDEO GENERATION SETTINGS
# ==============================================

# AI Provider: "did", "runwayml", "synthesia", or "mock"
export AI_PROVIDER="did"

# Use full AI pipeline (DID + Shotstack)
export USE_FULL_AI_PIPELINE="true"

# ==============================================
# 🚀 START BACKEND
# ==============================================

echo "======================================================"
echo "🚀 Starting Backend Server"
echo "======================================================"
echo ""
echo "📋 Configuration:"
echo "  - AI Provider: $AI_PROVIDER"
echo "  - Use V0 Style: $USE_V0_STYLE"
echo "  - Full AI Pipeline: $USE_FULL_AI_PIPELINE"
echo ""

if [ "$GOOGLE_GEMINI_API_KEY" = "your-gemini-api-key-here" ]; then
    echo "⚠️  WARNING: Gemini API key not set!"
    echo "💡 Edit start.sh and add your Gemini API key"
    echo "   Get key from: https://makersuite.google.com/app/apikey"
    echo ""
fi

if [ "$DID_API_KEY" = "your-did-api-key-here" ]; then
    echo "⚠️  WARNING: D-ID API key not set!"
    echo "💡 Edit start.sh and add your D-ID API key"
    echo ""
fi

if [ "$SHOTSTACK_API_KEY" = "your-shotstack-api-key-here" ]; then
    echo "⚠️  WARNING: Shotstack API key not set!"
    echo "💡 Edit start.sh and add your Shotstack API key"
    echo ""
fi

echo "======================================================"
echo ""

# Start the backend
./backend

