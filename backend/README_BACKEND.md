# Backend API Documentation

## Overview

Golang-based REST API server using Gin framework for AI-powered product marketing automation.

## Architecture

```
backend/
├── internal/
│   ├── config/           # Configuration management
│   ├── database/         # Database setup (SQLite/GORM)
│   ├── handlers/         # HTTP request handlers
│   ├── models/           # Data models (Project)
│   ├── router/           # API routes & CORS
│   └── services/         # Business logic
│       ├── ai_service.go          # Main AI orchestration
│       ├── video_generator.go    # AI video generation
│       └── website_templates.go  # Website generation
└── main.go              # Application entry point
```

## API Endpoints

### Upload Media
```http
POST /api/v1/upload
Content-Type: multipart/form-data

Form Data:
- product_image: File (JPG, PNG, GIF, WEBP)
- person_media: File (Image or Video: MP4, MOV, AVI)

Response:
{
  "project_id": "uuid",
  "status": "uploaded",
  "message": "Files uploaded successfully"
}
```

### List Projects
```http
GET /api/v1/projects

Response:
{
  "projects": [
    {
      "id": "uuid",
      "product_image_path": "/uploads/product.jpg",
      "person_media_path": "/uploads/person.jpg",
      "person_media_type": "image",
      "generated_video_path": "/generated/videos/video.mp4",
      "website_path": "/generated/websites/site-id",
      "status": "website_complete",
      "created_at": "2025-11-04T19:50:00Z",
      "updated_at": "2025-11-04T19:55:00Z"
    }
  ]
}
```

### Get Project
```http
GET /api/v1/projects/:id

Response:
{
  "id": "uuid",
  "product_image_path": "/uploads/product.jpg",
  ...
}
```

### Generate Video
```http
POST /api/v1/projects/:id/generate-video

Response:
{
  "project_id": "uuid",
  "video_path": "/generated/videos/video.mp4",
  "status": "video_complete"
}
```

### Generate Website
```http
POST /api/v1/projects/:id/generate-website

Response:
{
  "project_id": "uuid",
  "website_path": "/generated/websites/site-id",
  "status": "website_complete"
}
```

## Status Flow

1. `uploaded` - Files uploaded
2. `video_generating` - Video generation in progress
3. `video_complete` - Video ready
4. `website_generating` - Website generation in progress
5. `website_complete` - Website ready
6. `deployed` - (Future) Website deployed

## Database Schema

### Project Model
```go
type Project struct {
    ID                 string    // UUID
    ProductImagePath   string    // File path
    PersonMediaPath    string    // File path
    PersonMediaType    string    // "image" or "video"
    GeneratedVideoPath string    // Generated video path
    WebsitePath        string    // Generated website directory
    WebsiteURL         string    // Deployed URL (future)
    Status             string    // Current status
    CreatedAt          time.Time
    UpdatedAt          time.Time
}
```

## AI Service Integration

### Video Generation

**Supported Providers:**
- **D-ID**: Talking head videos with AI avatars
- **Synthesia**: Professional AI avatar presentations  
- **RunwayML**: Creative video generation
- **Mock**: Placeholder for testing

**Configuration:**
```bash
export AI_PROVIDER=did  # or synthesia, runwayml, mock
export AI_API_KEY=your_api_key_here
```

### Website Generation

Generates professional marketing websites with:
- Hero section with product showcase
- Feature highlights
- Video demo section
- Call-to-action
- Responsive design
- Modern CSS (gradients, animations)
- Interactive JavaScript

## Running the Server

### Development
```bash
cd backend
go mod tidy
go run main.go
```

### Production Build
```bash
cd backend
go build -o server main.go
./server
```

### With Docker
```bash
docker build -t marketing-backend .
docker run -p 8080:8080 marketing-backend
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| PORT | 8080 | Server port |
| DATABASE_PATH | ./data/app.db | SQLite database file |
| UPLOAD_PATH | ./uploads | Uploaded files directory |
| GENERATED_VIDEO_PATH | ./generated/videos | Generated videos directory |
| WEBSITE_PATH | ./generated/websites | Generated websites directory |
| AI_PROVIDER | mock | AI service provider |
| AI_API_KEY | - | API key for AI service |
| AI_API_URL | - | API URL (for some services) |

## Static File Serving

Files are served at:
- `/static/uploads/*` - Uploaded product images and person media
- `/static/generated/videos/*` - Generated promotional videos
- `/static/generated/websites/*` - Generated marketing websites

## CORS Configuration

CORS is enabled for all origins in development:
```go
Access-Control-Allow-Origin: *
Access-Control-Allow-Methods: POST, GET, PUT, DELETE, OPTIONS
Access-Control-Allow-Headers: Content-Type, Authorization
```

**⚠️ Important**: Configure specific origins for production!

## Error Handling

Standard HTTP status codes:
- `200` - Success
- `201` - Created
- `400` - Bad request (validation error)
- `404` - Not found
- `500` - Internal server error

Error response format:
```json
{
  "error": "Error message",
  "details": "Detailed error information (if available)"
}
```

## Performance Considerations

- Video generation: 2-10 minutes (depends on AI service)
- Website generation: < 1 second
- File uploads: Depends on file size
- Database: SQLite (suitable for MVP, use PostgreSQL for scale)

## Testing

```bash
# Run tests
go test ./...

# Test specific package
go test ./internal/handlers

# With coverage
go test -cover ./...
```

## Dependencies

Key dependencies (see `go.mod`):
- **gin-gonic/gin** - HTTP web framework
- **gorm.io/gorm** - ORM
- **gorm.io/driver/sqlite** - SQLite driver
- **google/uuid** - UUID generation

## Security Notes

- [ ] Add authentication middleware
- [ ] Validate file types and sizes
- [ ] Sanitize file names
- [ ] Rate limiting
- [ ] API key validation
- [ ] Secure CORS configuration
- [ ] HTTPS in production

## Monitoring & Logging

Current: Console logging with Gin's default logger

Production recommendations:
- Structured logging (zap, logrus)
- Error tracking (Sentry)
- Metrics (Prometheus)
- Distributed tracing (Jaeger)

## Scaling

For production scale:
1. Switch to PostgreSQL
2. Add Redis for caching
3. Use S3/Cloud Storage for files
4. Add message queue (RabbitMQ) for video processing
5. Implement worker pools
6. Add load balancer
7. Containerize with Kubernetes

## API Examples

### cURL Examples

**Upload files:**
```bash
curl -X POST http://localhost:8080/api/v1/upload \
  -F "product_image=@product.jpg" \
  -F "person_media=@person.jpg"
```

**Generate video:**
```bash
curl -X POST http://localhost:8080/api/v1/projects/PROJECT_ID/generate-video
```

**List projects:**
```bash
curl http://localhost:8080/api/v1/projects
```

## Support

See main README.md for setup instructions and troubleshooting.
See AI_SERVICES_GUIDE.md for AI integration details.


