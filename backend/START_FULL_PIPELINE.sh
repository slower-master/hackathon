#!/bin/bash

# 🚀 Start Full AI Pipeline - Startup Script
# This script starts the backend with full AI pipeline enabled

echo "🚀 ========================================"
echo "🚀 Starting Full AI Video Pipeline"
echo "🚀 ========================================"
echo ""

# Navigate to backend directory
cd "$(dirname "$0")"

# Check if .env file exists
if [ ! -f .env ]; then
    echo "❌ ERROR: .env file not found!"
    echo ""
    echo "Please create .env file with the following content:"
    echo ""
    echo "# Backend Configuration"
    echo "PORT=8080"
    echo "DATABASE_PATH=./data/app.db"
    echo "UPLOAD_PATH=./uploads"
    echo "GENERATED_VIDEO_PATH=./generated/videos"
    echo "WEBSITE_PATH=./generated/websites"
    echo ""
    echo "# AI Provider"
    echo "AI_PROVIDER=did"
    echo ""
    echo "# Enable Full AI Pipeline"
    echo "USE_FULL_AI_PIPELINE=true"
    echo ""
    echo "# API Keys"
    echo "AI_API_KEY=your_did_key_here"
    echo "RUNWAYML_API_KEY=your_runwayml_key_here"
    echo "SHOTSTACK_API_KEY=your_shotstack_key_here"
    echo ""
    exit 1
fi

# Check if required API keys are set
echo "📋 Checking configuration..."
echo ""

# Source .env file
export $(grep -v '^#' .env | xargs)

# Verify D-ID API key
if [ -z "$AI_API_KEY" ]; then
    echo "❌ AI_API_KEY not set in .env!"
    exit 1
else
    echo "✅ D-ID API Key: ${AI_API_KEY:0:20}..."
fi

# Check if full pipeline is enabled
if [ "$USE_FULL_AI_PIPELINE" = "true" ]; then
    echo "✅ Full AI Pipeline: ENABLED"
    echo ""
    
    # Verify RunwayML API key
    if [ -z "$RUNWAYML_API_KEY" ]; then
        echo "⚠️  WARNING: RUNWAYML_API_KEY not set!"
        echo "   Full pipeline requires RunwayML API key"
        echo "   Get it from: https://app.runwayml.com/account/secrets"
        echo ""
    else
        echo "✅ RunwayML API Key: ${RUNWAYML_API_KEY:0:20}..."
    fi
    
    # Verify Shotstack API key
    if [ -z "$SHOTSTACK_API_KEY" ]; then
        echo "⚠️  WARNING: SHOTSTACK_API_KEY not set!"
        echo "   Full pipeline requires Shotstack API key"
        echo "   Get it from: https://dashboard.shotstack.io/"
        echo ""
    else
        echo "✅ Shotstack API Key: ${SHOTSTACK_API_KEY:0:20}..."
    fi
else
    echo "📝 Full AI Pipeline: DISABLED (using standard D-ID only)"
    echo ""
fi

echo ""
echo "🔧 ========================================"
echo "🔧 Starting Backend Server"
echo "🔧 ========================================"
echo ""

# Kill any existing backend processes
echo "Stopping existing backend processes..."
lsof -ti:8080 | xargs kill -9 2>/dev/null
pkill -f "go run main.go" 2>/dev/null
sleep 2

# Start backend
echo "Starting backend server..."
echo ""

# Run in foreground so we can see logs
go run main.go

# If go run fails, show error
if [ $? -ne 0 ]; then
    echo ""
    echo "❌ Failed to start backend!"
    echo "Check the error messages above"
    exit 1
fi

