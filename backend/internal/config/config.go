package config

import (
	"os"
)

type Config struct {
	DatabasePath       string
	UploadPath         string
	GeneratedVideoPath string
	WebsitePath        string
	AIAPIKey           string
	AIAPIURL           string
	AIProvider         string // "runwayml", "did", "synthesia", or "mock"
	Port               string
	// New AI service API keys
	RunwayMLAPIKey     string
	ShotstackAPIKey    string
	GeminiAPIKey       string // For generating website content and features
	UseFullAIPipeline  bool   // If true, uses RunwayML + Shotstack for complete AI pipeline
	UseV0Style         bool   // If true, uses modern v0.dev style for websites
}

func Load() *Config {
	// Try multiple environment variable names for Gemini API key
	geminiAPIKey := getEnv("GOOGLE_GEMINI_API_KEY", "")
	if geminiAPIKey == "" {
		geminiAPIKey = getEnv("GEMINI_API_KEY", "")
	}
	
	return &Config{
		DatabasePath:       getEnv("DATABASE_PATH", "./data/app.db"),
		UploadPath:         getEnv("UPLOAD_PATH", "./uploads"),
		GeneratedVideoPath: getEnv("GENERATED_VIDEO_PATH", "./generated/videos"),
		WebsitePath:        getEnv("WEBSITE_PATH", "./generated/websites"),
		AIAPIKey:           getEnv("AI_API_KEY", ""),
		AIAPIURL:           getEnv("AI_API_URL", ""),
		AIProvider:         getEnv("AI_PROVIDER", "mock"), // Options: runwayml, did, synthesia, mock
		Port:               getEnv("PORT", "8080"),
		// New AI service API keys
		RunwayMLAPIKey:     getEnv("RUNWAYML_API_KEY", ""),
		ShotstackAPIKey:    getEnv("SHOTSTACK_API_KEY", ""),
		GeminiAPIKey:       geminiAPIKey,
		UseFullAIPipeline:  getEnv("USE_FULL_AI_PIPELINE", "false") == "true",
		UseV0Style:         getEnv("USE_V0_STYLE", "true") == "true", // Default to true for modern websites
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
