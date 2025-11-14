package services

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
// productVideoStyle: "rotation", "zoom", "pan", "reveal", "auto" (default: "rotation")
// layout options (default: "presenter"):
//   - "presenter" (RECOMMENDED): Person 60% left, product 40% right - looks like real product explanation
//   - "split": Side-by-side 50/50 - balanced, professional
//   - "product_main": Product fullscreen + avatar overlay (traditional)
//   - "avatar_main": Avatar fullscreen + product overlay
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
	// Generate URLs for static assets (use actual uploaded files)
	productImageURL := fmt.Sprintf("/static/uploads/%s", filepath.Base(project.ProductImagePath))
	videoURL := ""
	if project.GeneratedVideoPath != "" {
		videoURL = fmt.Sprintf("/static/generated/videos/%s", filepath.Base(project.GeneratedVideoPath))
	}
	
	// Log URLs for debugging
	fmt.Printf("\n🖼️  Product Image URL: %s\n", productImageURL)
	fmt.Printf("🎥 Video URL: %s\n", videoURL)
	fmt.Printf("📁 Product Image Path: %s\n", project.ProductImagePath)

	// Use actual product details from the project
	productName := project.ProductName
	if productName == "" {
		productName = "Amazing Product"
	}
	
	productDescription := project.ProductDescription
	if productDescription == "" {
		productDescription = "Transform your experience with our innovative solution"
	}

	// Generate AI features using Gemini
	fmt.Printf("\n" + strings.Repeat("=", 60) + "\n")
	fmt.Printf("🤖 GEMINI: Generating Website Features\n")
	fmt.Printf(strings.Repeat("=", 60) + "\n")
	var features []map[string]string
	
	// HARDCODED Gemini API key for testing
	geminiAPIKey := "AIzaSyC_gI30tRdg-eYjVJn7ses22lrzrRB4vXc"
	
	fmt.Printf("🔑 Gemini API Key: HARDCODED for testing\n")
	fmt.Printf("📦 Product: %s\n", productName)
	fmt.Printf("📝 Description: %s\n", productDescription)
	fmt.Printf("🏷️  Category: %s\n", project.ProductCategory)
	fmt.Printf("💰 Price: %s\n", project.ProductPrice)
	
	geminiService := NewGeminiService(geminiAPIKey)
	aiFeatures, err := geminiService.GenerateWebsiteFeatures(productName, productDescription, project.ProductCategory, project.ProductPrice)
	if err != nil {
		fmt.Printf("❌ Gemini features generation failed: %v\n", err)
		fmt.Printf("⚠️  Using default features as fallback\n")
		features = getDefaultFeatures()
	} else {
		features = aiFeatures
		fmt.Printf("✅ Successfully generated %d AI features:\n", len(features))
		for i, f := range features {
			fmt.Printf("   %d. %s %s: %s\n", i+1, f["icon"], f["title"], f["description"])
		}
	}
	fmt.Printf(strings.Repeat("=", 60) + "\n\n")

	// Check if we should use v0.dev style generation (from config)
	if s.config.UseV0Style {
		fmt.Printf("🌐 Using v0.dev style website generation...\n")
		
		// Use v0.dev service for modern website generation
		v0Service := NewV0Service("")
		v0WebsiteDir, err := v0Service.GenerateWebsite(
			productName,
			productDescription,
			project.ProductPrice,
			productImageURL,
			videoURL,
			features,
		)
		if err == nil {
			// Copy generated files to target directory
			return s.copyWebsiteFiles(v0WebsiteDir, websiteDir)
		}
		fmt.Printf("⚠️  v0.dev generation failed, using default template\n")
	}

	// Generate professional marketing website HTML with dynamic product details and AI features
	htmlContent := MarketingWebsiteTemplate(
		productName,
		productDescription,
		videoURL,
		productImageURL,
		features,
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
	jsContent := ModernWebsiteJS()
	if err := os.WriteFile(filepath.Join(websiteDir, "script.js"), []byte(jsContent), 0644); err != nil {
		return err
	}

	return nil
}

// copyWebsiteFiles copies generated website files from source to target directory
func (s *AIService) copyWebsiteFiles(sourceDir, targetDir string) error {
	files := []string{"index.html", "styles.css", "script.js"}
	
	for _, filename := range files {
		sourcePath := filepath.Join(sourceDir, filename)
		targetPath := filepath.Join(targetDir, filename)
		
		// Read source file
		data, err := os.ReadFile(sourcePath)
		if err != nil {
			return fmt.Errorf("failed to read %s: %v", filename, err)
		}
		
		// Write to target
		if err := os.WriteFile(targetPath, data, 0644); err != nil {
			return fmt.Errorf("failed to write %s: %v", filename, err)
		}
	}

	return nil
}
