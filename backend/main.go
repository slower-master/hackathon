package main

import (
	"log"
	"os"

	"github.com/dealshare/hacathon/backend/internal/config"
	"github.com/dealshare/hacathon/backend/internal/database"
	"github.com/dealshare/hacathon/backend/internal/handlers"
	"github.com/dealshare/hacathon/backend/internal/router"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file (ignore error if file doesn't exist)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Load configuration
	cfg := config.Load()

	// Check for Gemini API key
	geminiKey := os.Getenv("GOOGLE_GEMINI_API_KEY")
	if geminiKey == "" {
		geminiKey = os.Getenv("GEMINI_API_KEY")
	}

	// Log configuration
	log.Printf("=== Configuration ===")
	log.Printf("AI Provider: %s", cfg.AIProvider)
	log.Printf("AI API Key: %s", func() string {
		if cfg.AIAPIKey != "" {
			return "***" + cfg.AIAPIKey[len(cfg.AIAPIKey)-10:]
		}
		return "(not set)"
	}())
	
	// Log Gemini status
	if geminiKey != "" {
		log.Printf("🤖 Google Gemini: ENABLED (%s...%s)", geminiKey[:8], geminiKey[len(geminiKey)-4:])
	} else {
		log.Printf("⚠️  Google Gemini: DISABLED (no API key found)")
		log.Printf("   Set GOOGLE_GEMINI_API_KEY environment variable to enable AI script generation")
	}
	log.Printf("====================")

	// Initialize database
	db, err := database.Initialize(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Run migrations
	if err := database.Migrate(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize handlers
	h := handlers.New(db, cfg)

	// Setup router
	r := router.Setup(h)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

