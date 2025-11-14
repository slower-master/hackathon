#!/bin/bash

# Test Script: Verify Gemini is Generating Product-Specific Features
# This script tests if Gemini API is working correctly

echo "======================================================"
echo "🧪 Testing Gemini Feature Generation"
echo "======================================================"
echo ""

# Check if Gemini API key is set
if [ -z "$GOOGLE_GEMINI_API_KEY" ]; then
    if [ -z "$GEMINI_API_KEY" ]; then
        echo "❌ ERROR: No Gemini API key found!"
        echo ""
        echo "Please set one of these environment variables:"
        echo "  export GOOGLE_GEMINI_API_KEY=\"your-key-here\""
        echo "  or"
        echo "  export GEMINI_API_KEY=\"your-key-here\""
        echo ""
        echo "Get your FREE key from:"
        echo "  https://makersuite.google.com/app/apikey"
        exit 1
    fi
    GEMINI_API_KEY=$GEMINI_API_KEY
else
    GEMINI_API_KEY=$GOOGLE_GEMINI_API_KEY
fi

echo "✅ Gemini API Key Found: ${GEMINI_API_KEY:0:10}..."
echo ""

# Test products
declare -a products=(
    "Namkeen papdi:Tasty Healthy Namkeen Packet:Snacks:29"
    "Smart Watch:Fitness tracker with heart rate monitor:Electronics:2999"
    "Organic Honey:Pure natural honey from Himalayan bees:Food:499"
)

echo "Testing with 3 different products..."
echo ""

for product_data in "${products[@]}"; do
    IFS=':' read -r name description category price <<< "$product_data"
    
    echo "------------------------------------------------------"
    echo "📦 Product: $name"
    echo "📝 Description: $description"
    echo "🏷️  Category: $category"
    echo "💰 Price: ₹$price"
    echo ""
    
    # Create the request payload
    PROMPT="You are an expert marketing copywriter. Generate exactly 4 compelling product features/benefits for a website landing page.

Product Name: $name
Description: $description
Category: $category
Price: ₹$price

REQUIREMENTS:
1. Generate EXACTLY 4 features
2. Each feature should have:
   - An emoji icon (🚀, 💎, 🔒, ⚡, 🎯, ✨, 🌟, 💪, 🎨, 🔥, etc.)
   - A short title (2-4 words)
   - A description (15-25 words)
3. Features should be relevant to the product description
4. Make them compelling and benefit-focused
5. Use varied emojis

OUTPUT FORMAT (JSON only, no other text):
{
  \"features\": [
    {
      \"icon\": \"🚀\",
      \"title\": \"Feature Title\",
      \"description\": \"Feature description here.\"
    }
  ]
}

Now generate features for this product:"

    # Call Gemini API
    echo "🔄 Calling Gemini API..."
    
    RESPONSE=$(curl -s -X POST \
        "https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent?key=$GEMINI_API_KEY" \
        -H "Content-Type: application/json" \
        -d "{
            \"contents\":[{
                \"parts\":[{
                    \"text\": $(echo "$PROMPT" | jq -Rs .)
                }]
            }],
            \"generationConfig\": {
                \"temperature\": 0.9,
                \"topK\": 40,
                \"topP\": 0.95,
                \"maxOutputTokens\": 1024
            }
        }")
    
    # Check for errors
    if echo "$RESPONSE" | jq -e '.error' > /dev/null 2>&1; then
        echo "❌ API Error:"
        echo "$RESPONSE" | jq '.error'
        echo ""
        continue
    fi
    
    # Extract the generated text
    GENERATED_TEXT=$(echo "$RESPONSE" | jq -r '.candidates[0].content.parts[0].text')
    
    if [ -z "$GENERATED_TEXT" ] || [ "$GENERATED_TEXT" = "null" ]; then
        echo "❌ No response generated"
        echo "Raw response:"
        echo "$RESPONSE" | jq '.'
        echo ""
        continue
    fi
    
    echo "✅ Generated Features:"
    echo ""
    
    # Try to parse as JSON and display nicely
    if echo "$GENERATED_TEXT" | jq -e '.features' > /dev/null 2>&1; then
        echo "$GENERATED_TEXT" | jq -r '.features[] | "  \(.icon) \(.title)\n     \(.description)\n"'
    else
        # Try to extract JSON from markdown code blocks
        CLEAN_JSON=$(echo "$GENERATED_TEXT" | sed -n '/```json/,/```/p' | sed '1d;$d')
        if [ ! -z "$CLEAN_JSON" ]; then
            echo "$CLEAN_JSON" | jq -r '.features[] | "  \(.icon) \(.title)\n     \(.description)\n"'
        else
            echo "$GENERATED_TEXT"
        fi
    fi
    
    echo ""
done

echo "======================================================"
echo "✅ Gemini Test Complete!"
echo ""
echo "If you see product-specific features above, Gemini is working! 🎉"
echo ""
echo "If you see errors or generic features:"
echo "  1. Check your API key is valid"
echo "  2. Verify you have API quota remaining"
echo "  3. Check network connectivity"
echo "======================================================"

