#!/bin/bash

# Test Gemini Feature Generation using Go code
# This tests the ACTUAL backend code

echo "======================================================"
echo "🧪 Testing Gemini Feature Generation (Go Code)"
echo "======================================================"
echo ""

# Try to find Gemini API key from various sources
if [ ! -z "$GOOGLE_GEMINI_API_KEY" ]; then
    echo "✅ Found GOOGLE_GEMINI_API_KEY in environment"
elif [ ! -z "$GEMINI_API_KEY" ]; then
    echo "✅ Found GEMINI_API_KEY in environment"
    export GOOGLE_GEMINI_API_KEY="$GEMINI_API_KEY"
elif [ -f .env ]; then
    echo "🔍 Loading from .env file..."
    export $(grep GOOGLE_GEMINI_API_KEY .env | xargs)
    if [ -z "$GOOGLE_GEMINI_API_KEY" ]; then
        export $(grep GEMINI_API_KEY .env | xargs)
        export GOOGLE_GEMINI_API_KEY="$GEMINI_API_KEY"
    fi
else
    echo "❌ No Gemini API key found!"
    echo ""
    echo "Please provide your Gemini API key:"
    echo ""
    echo "Option 1: As argument"
    echo "  ./test_gemini_go.sh YOUR_API_KEY"
    echo ""
    echo "Option 2: Set environment"
    echo "  export GOOGLE_GEMINI_API_KEY='your-key'"
    echo "  ./test_gemini_go.sh"
    echo ""
    echo "Option 3: In .env file"
    echo "  echo 'GOOGLE_GEMINI_API_KEY=your-key' > .env"
    echo "  ./test_gemini_go.sh"
    exit 1
fi

# If key provided as argument, use it
if [ ! -z "$1" ]; then
    export GOOGLE_GEMINI_API_KEY="$1"
    echo "✅ Using API key from command line argument"
fi

if [ -z "$GOOGLE_GEMINI_API_KEY" ]; then
    echo "❌ Still no API key found!"
    exit 1
fi

echo ""
echo "======================================================"
echo ""

# Run the Go test
go run test_gemini_features.go

# Check if it succeeded
if [ $? -eq 0 ]; then
    echo ""
    echo "======================================================"
    echo "🎉 SUCCESS! Gemini is generating custom features!"
    echo "======================================================"
    echo ""
    echo "Next steps:"
    echo "  1. Make sure backend has the same API key"
    echo "  2. Restart backend: ./start.sh"
    echo "  3. Upload product - features should be custom!"
else
    echo ""
    echo "======================================================"
    echo "❌ Test failed - check errors above"
    echo "======================================================"
fi

