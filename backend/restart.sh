#!/bin/bash

# Backend Restart Script
# Usage: ./restart.sh

echo "🛑 Stopping backend..."

# Kill processes on port 8080
lsof -ti:8080 | xargs kill -9 2>/dev/null

# Kill any Go backend processes
pkill -f "go run main.go" 2>/dev/null
pkill -f "backend/main.go" 2>/dev/null

sleep 2

echo "🚀 Starting backend..."

cd "$(dirname "$0")"

# Check if .env file exists
if [ ! -f .env ]; then
    echo "⚠️  Warning: .env file not found!"
    echo "Please create .env file with your configuration"
    echo "See .env.example for template"
fi

# Start backend with logging (it will load .env automatically)
nohup go run main.go > ../backend.log 2>&1 &

sleep 3

# Check if it's running
if curl -s http://localhost:8080/api/v1/projects > /dev/null; then
    echo "✅ Backend started successfully!"
    echo "📝 Logs: tail -f ../backend.log"
    echo "🌐 API: http://localhost:8080"
else
    echo "❌ Backend failed to start. Check logs: tail -f ../backend.log"
    exit 1
fi


