# Setup Guide

This guide will help you set up the AI Product Marketing Agent project.

## Prerequisites

- **Go** 1.21 or higher ([Install Go](https://golang.org/doc/install))
- **Node.js** 18+ and npm ([Install Node.js](https://nodejs.org/))
- **Docker** (optional, for containerized deployment)

## Step-by-Step Setup

### 1. Clone or Navigate to Project

```bash
cd /Users/slowermaster/DEALSHARE/hacathon
```

### 2. Backend Setup

```bash
cd backend

# Initialize Go modules (if not already done)
go mod tidy

# Download dependencies
go mod download

# Create necessary directories
mkdir -p data uploads generated/videos generated/websites

# Copy environment file (optional)
cp .env.example .env

# Run the server
go run main.go
```

The backend should now be running on `http://localhost:8080`

### 3. Frontend Setup

Open a new terminal:

```bash
cd frontend

# Install dependencies
npm install

# Copy environment file (optional)
cp .env.example .env.local

# Run development server
npm run dev
```

The frontend should now be running on `http://localhost:3000`

### 4. Verify Setup

1. Open `http://localhost:3000` in your browser
2. You should see the upload interface
3. Try uploading a product image and person photo/video
4. Check that the backend API responds at `http://localhost:8080/api/v1/projects`

## Using Make Commands

For convenience, use the Makefile:

```bash
# Setup everything
make setup-backend
make setup-frontend

# Run services
make run-backend    # Terminal 1
make run-frontend   # Terminal 2

# Or use Docker
make run-docker
```

## Troubleshooting

### Backend Issues

**Issue**: `go: cannot find module`
```bash
cd backend
go mod tidy
go mod download
```

**Issue**: Port already in use
```bash
# Change PORT in .env or set environment variable
export PORT=8081
go run main.go
```

**Issue**: Database errors
```bash
# Remove existing database and recreate
rm -rf backend/data
mkdir -p backend/data
go run main.go
```

### Frontend Issues

**Issue**: `npm install` fails
```bash
# Clear cache and reinstall
rm -rf node_modules package-lock.json
npm cache clean --force
npm install
```

**Issue**: Cannot connect to backend
- Check that backend is running on port 8080
- Verify `NEXT_PUBLIC_API_URL` in `.env.local` matches backend URL
- Check CORS settings in backend router

### Docker Issues

**Issue**: Port conflicts
```bash
# Edit docker-compose.yml to change ports
# Or stop conflicting services
```

**Issue**: Build fails
```bash
# Rebuild without cache
docker-compose build --no-cache
docker-compose up
```

## Testing the Application

1. **Upload Test**:
   - Upload a product image (JPG, PNG)
   - Upload a person photo or video
   - Click "Upload & Create Project"
   - Should see project ID and status

2. **Generate Video**:
   - After upload, click "Generate Video"
   - Status should change to "video_generating" then "video_complete"
   - Click "View Video" to see result

3. **Generate Website**:
   - After video generation, click "Generate Website"
   - Status should change to "website_generating" then "website_complete"
   - Click "View Website" to see generated site

## Next Steps

- Configure AI service integration (see README.md)
- Customize website templates
- Add authentication if needed
- Deploy to production

## Getting Help

- Check the main README.md for API documentation
- Review error logs in terminal output
- Check database at `backend/data/app.db` (use SQLite browser)


