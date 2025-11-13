# Commands to Run Backend

## Quick Start (Copy & Paste)

### 1. Stop any existing backend:
```bash
lsof -ti:8080 | xargs kill -9
```

### 2. Start Backend:
```bash
cd /Users/slowermaster/DEALSHARE/hacathon/backend
export AI_PROVIDER=did
export AI_API_KEY='ZGluZXNoLnJhbUBkZWFsc2hhcmUuaW46:31Ln-cQ0VJ-_sj61DbZm3'
go run main.go
```

---

## Where to See Logs

### Option 1: In the Terminal (Where you ran the command)
- When you run `go run main.go`, logs appear **directly in that terminal window**
- You'll see all D-ID API requests, responses, and curl commands there
- **Keep that terminal window open!**

### Option 2: In Log File (If running in background)
```bash
# View logs in real-time
tail -f /Users/slowermaster/DEALSHARE/hacathon/backend.log

# Or from project root:
tail -f backend/backend.log
```

---

## Easy Restart Script

I created a script for you:

```bash
cd /Users/slowermaster/DEALSHARE/hacathon/backend
./restart.sh
```

This will:
- Stop old backend
- Start new backend
- Show you where logs are

---

## What You'll See in Logs

When you generate a video, you'll see:

```
=== D-ID API Request ===
URL: https://api.d-id.com/talks
Method: POST
Raw API Key format check:
  Format detected: base64(email):apikey
  Decoded email: dinesh.ram@dealshare.in
  Using format: email:apikey (then base64 encoded)

Curl equivalent:
curl -X POST 'https://api.d-id.com/talks' \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Basic ...' \
  -d '{...}'

=== D-ID API Response ===
Status: 201 Created
Status Code: 201
Success!
```

---

## Troubleshooting

### "Port already in use"
```bash
lsof -ti:8080 | xargs kill -9
```

### "Command not found"
Make sure you're in the right directory:
```bash
cd /Users/slowermaster/DEALSHARE/hacathon/backend
```

### "Can't see logs"
- If you ran `go run main.go` in terminal, logs are in **that terminal**
- If running in background, check: `tail -f backend/backend.log`


