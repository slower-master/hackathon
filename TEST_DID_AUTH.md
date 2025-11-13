# D-ID Authentication Test

## Current Status: 401 Unauthorized

From your logs, the system is:
- Decoding email correctly: `dinesh.ram@dealshare.in`
- Creating auth header: `ZGluZXNoLnJhbUBkZWFsc2hhcmUuaW46MzFMbi1jUTBWSi1fc2o2MURiWm0z`
- But still getting 401 Unauthorized

## Possible Issues

### 1. API Key Might Be Invalid
- The API key may have been regenerated
- The account might need verification
- Credits might be exhausted

### 2. Test the API Key Directly

Run this curl command to test:

```bash
curl -X POST 'https://api.d-id.com/talks' \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Basic ZGluZXNoLnJhbUBkZWFsc2hhcmUuaW46MzFMbi1jUTBWSi1fc2o2MURiWm0z' \
  -d '{
    "source_url": "https://create-images-results.d-id.com/api_docs/assets/noelle.jpeg",
    "script": {
      "type": "text",
      "input": "This is a test video",
      "provider": {
        "type": "microsoft",
        "voice_id": "en-US-JennyNeural"
      }
    }
  }'
```

### 3. Get Fresh API Key

1. **Login to D-ID**:
   - Go to: https://studio.d-id.com/
   - Login with: `dinesh.ram@dealshare.in`

2. **Check Account Status**:
   - Verify email is confirmed
   - Check remaining credits
   - Look for any account restrictions

3. **Get New API Key**:
   - Go to: Account Settings → API Keys
   - Click "Create API Key" or "Regenerate"
   - Copy the NEW key
   - Format will be: `your_api_key_string`

4. **Update Backend**:
```bash
# Stop backend
lsof -ti:8080 | xargs kill -9

# Start with NEW API key
cd /Users/slowermaster/DEALSHARE/hacathon/backend
export AI_PROVIDER=did
export AI_API_KEY='YOUR_NEW_API_KEY_HERE'  # Just the key, not base64 encoded
go run main.go
```

## Alternative: Try Without Email Encoding

D-ID might accept the API key in a different format. Let me create a test:

```bash
# If your D-ID API key is just: abc123xyz456
# Try sending it as:

# Option 1: Just the key with colon
echo -n "abc123xyz456:" | base64
# Then use: Authorization: Basic <result>

# Option 2: Email:key format  
echo -n "dinesh.ram@dealshare.in:abc123xyz456" | base64
# Then use: Authorization: Basic <result>
```

## Check D-ID Dashboard

Login and check:
1. **Credits**: Do you have credits remaining?
2. **API Key**: Is it active (not revoked)?
3. **Account**: Is email verified?
4. **Plan**: Are you on trial or paid plan?
5. **Usage**: Check API usage logs

## If Still Failing

The issue might be:
- ❌ API key format is wrong for D-ID
- ❌ Account needs verification
- ❌ Credits exhausted
- ❌ API key was regenerated
- ❌ Wrong account email

## Next Steps

1. **Run the curl test above** - see exact error
2. **Login to D-ID dashboard** - verify account
3. **Get fresh API key** - from dashboard
4. **Update the API key** in backend and restart

Copy this command to test:
```bash
curl -v -X POST 'https://api.d-id.com/talks' \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Basic ZGluZXNoLnJhbUBkZWFsc2hhcmUuaW46MzFMbi1jUTBWSi1fc2o2MURiWm0z' \
  -d '{"source_url":"https://create-images-results.d-id.com/api_docs/assets/noelle.jpeg","script":{"type":"text","input":"Test","provider":{"type":"microsoft","voice_id":"en-US-JennyNeural"}}}'
```

The `-v` flag will show you the exact error from D-ID.


