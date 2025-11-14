#!/bin/bash

# Simple Gemini Test - Quick verification
# Usage: ./test_gemini_simple.sh YOUR_GEMINI_API_KEY

API_KEY="$1"

if [ -z "$API_KEY" ]; then
    echo "❌ ERROR: Please provide your Gemini API key"
    echo ""
    echo "Usage:"
    echo "  ./test_gemini_simple.sh YOUR_GEMINI_API_KEY"
    echo ""
    echo "Get your FREE key from:"
    echo "  https://makersuite.google.com/app/apikey"
    exit 1
fi

echo "======================================================"
echo "🧪 Testing Gemini API - Simple Version"
echo "======================================================"
echo ""
echo "🔑 API Key: ${API_KEY:0:10}..."
echo ""

# Test with Namkeen product
PRODUCT="Namkeen papdi"
DESCRIPTION="Tasty Healthy Namkeen Packet"

echo "📦 Testing with product: $PRODUCT"
echo "📝 Description: $DESCRIPTION"
echo ""
echo "🔄 Calling Gemini API..."
echo ""

# Create simple prompt
PROMPT="Generate 4 product features for: $PRODUCT - $DESCRIPTION

Return ONLY JSON format:
{
  \"features\": [
    {\"icon\": \"🥘\", \"title\": \"Feature Name\", \"description\": \"Feature description here.\"}
  ]
}"

# Call Gemini API
RESPONSE=$(curl -s -X POST \
    "https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent?key=$API_KEY" \
    -H "Content-Type: application/json" \
    -d "{
        \"contents\":[{
            \"parts\":[{
                \"text\": \"$PROMPT\"
            }]
        }]
    }" 2>&1)

# Check for curl errors
if [ $? -ne 0 ]; then
    echo "❌ Network Error:"
    echo "$RESPONSE"
    exit 1
fi

# Check for API errors
if echo "$RESPONSE" | grep -q '"error"'; then
    echo "❌ API Error:"
    echo "$RESPONSE" | grep -o '"message":"[^"]*"' | sed 's/"message":"//' | sed 's/"$//'
    echo ""
    echo "Common issues:"
    echo "  1. Invalid API key"
    echo "  2. API key not activated"
    echo "  3. Quota exceeded"
    echo ""
    echo "Full response:"
    echo "$RESPONSE"
    exit 1
fi

# Extract generated text
GENERATED_TEXT=$(echo "$RESPONSE" | grep -o '"text":"[^}]*' | sed 's/"text":"//' | sed 's/\\n/\n/g')

if [ -z "$GENERATED_TEXT" ]; then
    echo "❌ No response generated"
    echo ""
    echo "Raw response:"
    echo "$RESPONSE"
    exit 1
fi

echo "✅ Gemini API is working!"
echo ""
echo "Generated Response:"
echo "======================================================"
echo "$GENERATED_TEXT"
echo "======================================================"
echo ""

# Check if response contains product-specific keywords
if echo "$GENERATED_TEXT" | grep -iq "namkeen\|snack\|healthy\|tasty"; then
    echo "✅ SUCCESS: Response is product-specific!"
    echo ""
    echo "🎉 Gemini is generating custom features for your products!"
else
    echo "⚠️  Warning: Response might not be product-specific"
    echo "   But API is working correctly"
fi

echo ""
echo "======================================================"
echo "Next Steps:"
echo "  1. Add this key to backend/start.sh"
echo "  2. Restart backend: ./start.sh"
echo "  3. Generate website - features will be custom!"
echo "======================================================"

