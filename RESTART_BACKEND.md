# Restart Backend - Simple Commands

## Stop & Start Backend

```bash
# Stop backend
lsof -ti:8080 | xargs kill -9

# Start backend
cd /Users/slowermaster/DEALSHARE/hacathon/backend
export AI_PROVIDER=did
export AI_API_KEY='ZGluZXNoLnJhbUBkZWFsc2hhcmUuaW46:31Ln-cQ0VJ-_sj61DbZm3'
go run main.go
```

## View Logs

The logs appear in the **terminal where you ran `go run main.go`**

Keep that terminal window open to see:
- D-ID API requests
- Response status codes
- Any errors

## Test API Key

If you get 401 errors, test the key:

```bash
./test_did_key.sh
```

This will try all authentication methods and show you which (if any) work.

## If Still 401

The code is correct now. The issue is likely:
- API key expired/invalid
- No credits in D-ID account
- Account not verified

Check: https://studio.d-id.com/


