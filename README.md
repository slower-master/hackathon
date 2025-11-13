# AI-Powered Product Marketing Agent

An end-to-end product marketing automation system that generates promotional videos and websites from product images and person photos/videos.

## MVP Features

- ✅ Upload product image and person photo/video
- ✅ AI-powered promotional video generation (placeholder - ready for integration)
- ✅ Automatic website creation with editable interface
- ✅ Project management and status tracking
- 🚧 Social media integration (future)
- 🚧 Engagement analytics (future)

## Tech Stack

- **Backend**: Golang (Gin framework)
- **Frontend**: React/Next.js with TypeScript
- **Database**: SQLite (MVP), PostgreSQL (production)
- **AI Services**: [Ready for integration - RunwayML/Pika Labs for video generation]

## Project Structure

```
├── backend/
│   ├── internal/
│   │   ├── config/      # Configuration management
│   │   ├── database/    # Database setup and migrations
│   │   ├── handlers/    # HTTP request handlers
│   │   ├── models/      # Data models
│   │   ├── router/      # API routes
│   │   └── services/    # Business logic (AI integration)
│   ├── main.go          # Application entry point
│   └── go.mod           # Go dependencies
├── frontend/
│   ├── app/             # Next.js app directory
│   ├── lib/             # API client utilities
│   └── package.json     # Node dependencies
├── docker-compose.yml   # Docker orchestration
├── Makefile            # Build commands
└── README.md
```

## Quick Start

### Option 1: Using Make (Recommended)

```bash
# Setup both backend and frontend
make setup-backend
make setup-frontend

# Run backend (terminal 1)
make run-backend

# Run frontend (terminal 2)
make run-frontend
```

### Option 2: Docker Compose

```bash
# Build and run everything
make run-docker
# or
docker-compose up --build
```

### Option 3: Manual Setup

#### Backend Setup
```bash
cd backend
go mod tidy
go mod download
go run main.go
```

The backend will start on `http://localhost:8080`

#### Frontend Setup
```bash
cd frontend
npm install
npm run dev
```

The frontend will start on `http://localhost:3000`

## Environment Variables

### Backend (.env or environment)
```bash
PORT=8080
DATABASE_PATH=./data/app.db
UPLOAD_PATH=./uploads
GENERATED_VIDEO_PATH=./generated/videos
WEBSITE_PATH=./generated/websites

# AI Service Configuration (when ready)
AI_API_KEY=your_api_key_here
AI_API_URL=https://api.example.com/v1
```

### Frontend (.env.local)
```bash
NEXT_PUBLIC_API_URL=http://localhost:8080
```

## API Endpoints

### Upload & Project Management
- `POST /api/v1/upload` - Upload product image and person media
  - Form data: `product_image` (file), `person_media` (file)
  - Returns: `{ project_id, status, message }`

- `GET /api/v1/projects` - List all projects
  - Returns: `{ projects: [...] }`

- `GET /api/v1/projects/:id` - Get project details
  - Returns: `Project` object

### Generation
- `POST /api/v1/projects/:id/generate-video` - Generate promotional video
  - Returns: `{ project_id, video_path, status }`

- `POST /api/v1/projects/:id/generate-website` - Generate website
  - Returns: `{ project_id, website_path, status }`

## Project Status Flow

1. `uploaded` - Files uploaded successfully
2. `video_generating` - Video generation in progress
3. `video_complete` - Video generated successfully
4. `website_generating` - Website generation in progress
5. `website_complete` - Website generated successfully
6. `deployed` - Website deployed (future)

## AI Service Integration

The system is ready for AI service integration. To integrate:

1. Update `backend/internal/services/ai_service.go`
2. Implement the `GenerateVideo` method with your chosen AI service:
   - RunwayML
   - Pika Labs
   - Synthesia
   - D-ID
   - Or any other video generation API

3. Set `AI_API_KEY` and `AI_API_URL` environment variables

## Development

### Running Tests
```bash
# Backend tests
cd backend
go test ./...

# Frontend tests
cd frontend
npm test
```

### Building for Production
```bash
# Backend
cd backend
go build -o bin/server main.go

# Frontend
cd frontend
npm run build
```

## Next Steps

- [ ] Integrate actual AI video generation service
- [ ] Add website editing interface
- [ ] Implement website deployment/hosting
- [ ] Add user authentication
- [ ] Add project templates
- [ ] Implement social media posting (future)
- [ ] Add analytics and engagement tracking (future)

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Submit a pull request

## License

MIT License

