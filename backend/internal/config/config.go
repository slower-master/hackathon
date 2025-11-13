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
	UseFullAIPipeline  bool   // If true, uses RunwayML + Shotstack for complete AI pipeline
}

func Load() *Config {
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
		UseFullAIPipeline:  getEnv("USE_FULL_AI_PIPELINE", "false") == "true",
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
