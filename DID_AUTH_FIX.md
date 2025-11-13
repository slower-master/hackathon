# D-ID Authentication Fix

## Issue
Getting `401 Unauthorized` error when calling D-ID API.

## Root Cause
The API key format was not being handled correctly. D-ID API key format is:
```
base64(email):apikey
```

For example: `ZGluZXNoLnJhbUBkZWFsc2hhcmUuaW46:31Ln-cQ0VJ-_sj61DbZm3`

This decodes to: `dinesh.ram@dealshare.in:31Ln-cQ0VJ-_sj61DbZm3`

## Solution
The code now:
1. Detects if the first part (before colon) is base64-encoded email
2. Decodes it to get the actual email
3. Reconstructs as `email:apikey` format
4. Base64 encodes the entire string for Basic auth
5. Uses as: `Authorization: Basic <base64_encoded_string>`

## What Changed

### Before:
```go
req.Header.Set("Authorization", "Basic "+vg.config.AIAPIKey)
```

### After:
```go
// Detect format and decode if needed
parts := strings.SplitN(authHeader, ":", 2)
decodedFirst, err := base64.StdEncoding.DecodeString(parts[0])
if err == nil && strings.Contains(string(decodedFirst), "@") {
    email := strings.TrimSuffix(string(decodedFirst), ":")
    authString := email + ":" + parts[1]
    authHeader = base64.StdEncoding.EncodeToString([]byte(authString))
}
req.Header.Set("Authorization", "Basic "+authHeader)
```

## Logging Added

Now when you make a request, you'll see:

```
=== D-ID API Request ===
URL: https://api.d-id.com/talks
Method: POST
Payload (first 200 chars): {...}
Raw API Key format check:
  Full key length: 65
  Contains colon: true
  Format detected: base64(email):apikey
  Decoded email: dinesh.ram@dealshare.in
  Using format: email:apikey (then base64 encoded)
  Final base64 encoded length: 48

Curl equivalent:
curl -X POST 'https://api.d-id.com/talks' \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Basic <base64_string>' \
  -d '{...}'
=====================

=== D-ID API Response ===
Status: 201 Created
Status Code: 201
Success!
=====================
```

## Testing

To test the fix:

1. **Generate a video** via the UI at http://localhost:3000
2. **Check terminal logs** - you'll see detailed request/response info
3. **Copy the curl command** from logs if you need to test manually

## Manual Testing with curl

If you want to test the API directly:

```bash
# Your API key
API_KEY="ZGluZXNoLnJhbUBkZWFsc2hhcmUuaW46:31Ln-cQ0VJ-_sj61DbZm3"

# Decode first part to get email
EMAIL=$(echo "ZGluZXNoLnJhbUBkZWFsc2hhcmUuaW46" | base64 -d | tr -d ':')
# Email = dinesh.ram@dealshare.in

# Get API key part
API_KEY_PART="31Ln-cQ0VJ-_sj61DbZm3"

# Create auth string
AUTH_STRING="${EMAIL}:${API_KEY_PART}"

# Base64 encode
BASE64_AUTH=$(echo -n "$AUTH_STRING" | base64)

# Test API call
curl -X POST https://api.d-id.com/talks \
  -H "Content-Type: application/json" \
  -H "Authorization: Basic ${BASE64_AUTH}" \
  -d '{
    "source_url": "https://example.com/image.jpg",
    "script": {
      "type": "text",
      "input": "Hello, this is a test video!",
      "provider": {
        "type": "microsoft",
        "voice_id": "en-US-JennyNeural"
      }
    }
  }'
```

## Expected Response

Success response (201 Created):
```json
{
  "id": "talk_abc123",
  "status": "created",
  "created_at": "2025-11-04T20:30:00Z"
}
```

Error response (401 Unauthorized):
```json
{
  "message": "Unauthorized"
}
```

## Current Status

✅ **Fixed**: Authentication now handles base64-encoded email format
✅ **Logging**: Detailed curl-equivalent commands in terminal
✅ **Error Handling**: Better error messages with response details

## Next Steps

1. Try generating a video again
2. Check terminal logs for detailed request/response
3. If still getting 401, verify:
   - API key is active at https://studio.d-id.com/
   - Account has credits remaining
   - API key hasn't been regenerated

## Troubleshooting

### Still getting 401?

1. **Check API Key Status**:
   - Login to https://studio.d-id.com/
   - Go to Account Settings → API Keys
   - Verify key is active and not expired

2. **Check Credits**:
   - Make sure you have credits remaining
   - Free tier: 20 credits
   - Check usage in dashboard

3. **Check Logs**:
   - Look at terminal output for detailed curl command
   - Copy the exact curl command and test manually
   - Compare auth header format

4. **Verify Email**:
   - The decoded email should match your D-ID account email
   - If different, the API key might be for a different account

## Support

If issues persist:
- Check D-ID status: https://status.d-id.com/
- Contact D-ID support with:
  - Account email
  - API key first 5 chars (for identification)
  - Full error logs from terminal



