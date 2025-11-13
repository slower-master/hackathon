#!/bin/bash

# Test D-ID API Key
# Run: ./test_did_key.sh

API_KEY="ZGluZXNoLnJhbUBkZWFsc2hhcmUuaW46:31Ln-cQ0VJ-_sj61DbZm3"

echo "Testing D-ID API Key..."
echo "Key format: $API_KEY"
echo ""

# Method 1: Use as-is (base64 encode entire string)
echo "Method 1: Base64 encode entire key as-is"
AUTH1=$(echo -n "$API_KEY" | base64)
echo "Auth: $AUTH1"
RESPONSE1=$(curl -s -X POST 'https://api.d-id.com/talks' \
  -H 'Content-Type: application/json' \
  -H "Authorization: Basic $AUTH1" \
  -d '{"source_url":"https://create-images-results.d-id.com/api_docs/assets/noelle.jpeg","script":{"type":"text","input":"Test","provider":{"type":"microsoft","voice_id":"en-US-JennyNeural"}}}')
echo "Response: $RESPONSE1"
echo ""

# Method 2: Decode first part, then combine
echo "Method 2: Decode first part (if base64), then combine"
USERNAME=$(echo -n "ZGluZXNoLnJhbUBkZWFsc2hhcmUuaW46" | base64 -d 2>/dev/null)
if [ $? -eq 0 ]; then
    EMAIL=$(echo "$USERNAME" | tr -d ':')
    PASSWORD="31Ln-cQ0VJ-_sj61DbZm3"
    AUTH2=$(echo -n "$EMAIL:$PASSWORD" | base64)
    echo "Decoded email: $EMAIL"
    echo "Auth: $AUTH2"
    RESPONSE2=$(curl -s -X POST 'https://api.d-id.com/talks' \
      -H 'Content-Type: application/json' \
      -H "Authorization: Basic $AUTH2" \
      -d '{"source_url":"https://create-images-results.d-id.com/api_docs/assets/noelle.jpeg","script":{"type":"text","input":"Test","provider":{"type":"microsoft","voice_id":"en-US-JennyNeural"}}}')
    echo "Response: $RESPONSE2"
else
    echo "First part is not base64, skipping Method 2"
fi
echo ""

# Method 3: Try without base64 encoding (just the key)
echo "Method 3: Use key directly without base64"
RESPONSE3=$(curl -s -X POST 'https://api.d-id.com/talks' \
  -H 'Content-Type: application/json' \
  -H "Authorization: Basic $API_KEY" \
  -d '{"source_url":"https://create-images-results.d-id.com/api_docs/assets/noelle.jpeg","script":{"type":"text","input":"Test","provider":{"type":"microsoft","voice_id":"en-US-JennyNeural"}}}')
echo "Response: $RESPONSE3"
echo ""

echo "Check which method works (look for 'id' or 'created_at' in response)"
echo "If all show 'Unauthorized', the API key itself is invalid/expired."



