package services

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dealshare/hacathon/backend/internal/config"
	"github.com/dealshare/hacathon/backend/internal/models"
	"github.com/google/uuid"
)

type AIService struct {
	config         *config.Config
	videoGenerator *VideoGenerator
}

func NewAIService(cfg *config.Config) *AIService {
	return &AIService{
		config:         cfg,
		videoGenerator: NewVideoGenerator(cfg),
	}
}

// GenerateVideo generates a promotional video combining product image and person media
// productVideoStyle: "rotation", "zoom", "pan", "reveal", "auto" (default: "auto")
// layout: "product_main" (product fullscreen + avatar overlay) or "avatar_main" (avatar fullscreen + product overlay) (default: "product_main")
func (s *AIService) GenerateVideo(productImagePath, personMediaPath, personMediaType, customScript, productVideoStyle, layout string) (string, error) {
	// Create output directory
	os.MkdirAll(s.config.GeneratedVideoPath, 0755)

	// Route to appropriate AI provider
	switch s.config.AIProvider {
	case "runwayml":
		return s.videoGenerator.GenerateWithRunwayML(productImagePath, personMediaPath, customScript)
	
	case "did":
		return s.videoGenerator.GenerateWithDID(productImagePath, personMediaPath, customScript, productVideoStyle, layout)
	
	case "synthesia":
		return s.videoGenerator.GenerateWithSynthesia(productImagePath, customScript)
	
	case "mock":
		fallthrough
	default:
		// Mock implementation for development/testing
		videoID := uuid.New().String()
		outputPath := filepath.Join(s.config.GeneratedVideoPath, fmt.Sprintf("%s.mp4", videoID))
		return s.generateMockVideo(productImagePath, personMediaPath, outputPath)
	}
}

// generateMockVideo creates a placeholder video file for MVP testing
func (s *AIService) generateMockVideo(productImagePath, personMediaPath, outputPath string) (string, error) {
	// Create a simple placeholder file
	// In production, this would be replaced with actual video generation
	file, err := os.Create(outputPath)
	if err != nil {
		return "", fmt.Errorf("failed to create video file: %w", err)
	}
	defer file.Close()

	// Write placeholder content (in real implementation, this would be the actual video bytes)
	file.WriteString("Placeholder video - integrate with AI service")
	return outputPath, nil
}

// GenerateWebsite generates a website for the product
func (s *AIService) GenerateWebsite(project models.Project) (string, error) {
	// Create output directory
	os.MkdirAll(s.config.WebsitePath, 0755)

	// Generate unique website directory
	websiteID := uuid.New().String()
	websiteDir := filepath.Join(s.config.WebsitePath, websiteID)
	os.MkdirAll(websiteDir, 0755)

	// Generate HTML, CSS, and JS files
	if err := s.generateWebsiteFiles(project, websiteDir); err != nil {
		return "", fmt.Errorf("failed to generate website files: %w", err)
	}

	return websiteDir, nil
}

// generateWebsiteFiles creates HTML, CSS, and JS files for the website
func (s *AIService) generateWebsiteFiles(project models.Project, websiteDir string) error {
	// Generate URLs for static assets
	productImageURL := fmt.Sprintf("/static/uploads/%s", filepath.Base(project.ProductImagePath))
	videoURL := ""
	if project.GeneratedVideoPath != "" {
		videoURL = fmt.Sprintf("/static/generated/videos/%s", filepath.Base(project.GeneratedVideoPath))
	}

	// Generate professional marketing website HTML
	htmlContent := MarketingWebsiteTemplate(
		"Amazing Product",
		"Transform your experience with our innovative solution",
		videoURL,
		productImageURL,
	)

	if err := os.WriteFile(filepath.Join(websiteDir, "index.html"), []byte(htmlContent), 0644); err != nil {
		return err
	}

	// Generate modern CSS
	cssContent := ModernWebsiteCSS()
	if err := os.WriteFile(filepath.Join(websiteDir, "styles.css"), []byte(cssContent), 0644); err != nil {
		return err
	}

	// Generate interactive JavaScript
	jsContent := InteractiveJS()
	if err := os.WriteFile(filepath.Join(websiteDir, "script.js"), []byte(jsContent), 0644); err != nil {
		return err
	}

	return nil
}
