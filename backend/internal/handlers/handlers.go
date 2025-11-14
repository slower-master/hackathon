package handlers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dealshare/hacathon/backend/internal/config"
	"github.com/dealshare/hacathon/backend/internal/models"
	"github.com/dealshare/hacathon/backend/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handlers struct {
	db        *gorm.DB
	config    *config.Config
	aiService *services.AIService
}

func New(db *gorm.DB, cfg *config.Config) *Handlers {
	return &Handlers{
		db:        db,
		config:    cfg,
		aiService: services.NewAIService(cfg),
	}
}

// UploadMedia handles product image and person media uploads
func (h *Handlers) UploadMedia(c *gin.Context) {
	// Parse multipart form
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(400, gin.H{"error": "Failed to parse form"})
		return
	}

	// Get product image
	productFiles := form.File["product_image"]
	if len(productFiles) == 0 {
		c.JSON(400, gin.H{"error": "product_image is required"})
		return
	}
	productFile := productFiles[0]

	// Get person media (image or video)
	personFiles := form.File["person_media"]
	if len(personFiles) == 0 {
		c.JSON(400, gin.H{"error": "person_media is required"})
		return
	}
	personFile := personFiles[0]

	// Get product details from form
	productName := c.PostForm("product_name")
	productDescription := c.PostForm("product_description")
	productCategory := c.PostForm("product_category")
	productPrice := c.PostForm("product_price")

	// Determine media type
	personMediaType := "image"
	if strings.HasSuffix(strings.ToLower(personFile.Filename), ".mp4") ||
		strings.HasSuffix(strings.ToLower(personFile.Filename), ".mov") ||
		strings.HasSuffix(strings.ToLower(personFile.Filename), ".avi") {
		personMediaType = "video"
	}

	// Create upload directories
	os.MkdirAll(h.config.UploadPath, 0755)

	// Save files
	productPath := filepath.Join(h.config.UploadPath, productFile.Filename)
	if err := c.SaveUploadedFile(productFile, productPath); err != nil {
		c.JSON(500, gin.H{"error": "Failed to save product image"})
		return
	}

	personPath := filepath.Join(h.config.UploadPath, personFile.Filename)
	if err := c.SaveUploadedFile(personFile, personPath); err != nil {
		c.JSON(500, gin.H{"error": "Failed to save person media"})
		return
	}

	// Generate AI script from product description using Google Gemini Pro - MANDATORY!
	fmt.Printf("\n" + strings.Repeat("=", 60) + "\n")
	fmt.Printf("🤖 GEMINI AI SCRIPT GENERATION (REQUIRED)\n")
	fmt.Printf(strings.Repeat("=", 60) + "\n")
	
	geminiAPIKey := os.Getenv("GOOGLE_GEMINI_API_KEY")
	if geminiAPIKey == "" {
		geminiAPIKey = os.Getenv("GEMINI_API_KEY") // Fallback to shorter name
	}
	
	if geminiAPIKey == "" {
		fmt.Printf("❌ ERROR: GOOGLE_GEMINI_API_KEY is REQUIRED!\n")
		fmt.Printf("   Please set GOOGLE_GEMINI_API_KEY in your .env file\n")
		fmt.Printf("   Get your API key from: https://ai.google.dev/\n")
		fmt.Printf(strings.Repeat("=", 60) + "\n\n")
		c.JSON(400, gin.H{
			"error": "Gemini API key is required. Please set GOOGLE_GEMINI_API_KEY in your environment variables.",
		})
		return
	}
	
	fmt.Printf("🔑 API Key: %s...%s (%d chars)\n", geminiAPIKey[:10], geminiAPIKey[len(geminiAPIKey)-4:], len(geminiAPIKey))
	fmt.Printf("📦 Product: %s\n", productName)
	fmt.Printf("📝 Description: %s\n", productDescription)
	fmt.Printf(strings.Repeat("-", 60) + "\n")
	
	geminiService := services.NewGeminiService(geminiAPIKey)
	generatedScript, err := geminiService.GenerateMarketingScript(productName, productDescription, productCategory, productPrice)
	if err != nil {
		fmt.Printf("❌ GEMINI FAILED: %v\n", err)
		fmt.Printf("❌ Cannot proceed without Gemini-generated script!\n")
		fmt.Printf(strings.Repeat("=", 60) + "\n\n")
		c.JSON(500, gin.H{
			"error": fmt.Sprintf("Failed to generate AI script with Gemini: %v. Please check your API key and try again.", err),
		})
		return
	}
	
	fmt.Printf("✅ GEMINI SUCCESS!\n")
	fmt.Printf("📝 Generated Script:\n")
	fmt.Printf("   \"%s\"\n", generatedScript)
	fmt.Printf(strings.Repeat("=", 60) + "\n\n")

	// Create project record
	project := &models.Project{
		ProductImagePath:   productPath,
		PersonMediaPath:    personPath,
		PersonMediaType:    personMediaType,
		ProductName:        productName,
		ProductDescription: productDescription,
		ProductCategory:    productCategory,
		ProductPrice:       productPrice,
		GeneratedScript:    generatedScript,
		Status:             "uploaded",
	}

	if err := h.db.Create(project).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to create project"})
		return
	}

	c.JSON(201, gin.H{
		"project_id":        project.ID,
		"status":            project.Status,
		"message":           "Files uploaded successfully",
		"generated_script":  generatedScript,
		"product_name":      productName,
		"product_description": productDescription,
	})
}

// GenerateVideo generates a promotional video from uploaded media
func (h *Handlers) GenerateVideo(c *gin.Context) {
	projectID := c.Param("id")

	// Parse request body for video options (NO custom script - always use Gemini-generated)
	var requestBody struct {
		ProductVideoStyle string `json:"product_video_style"` // "rotation", "zoom", "pan", "reveal", "auto"
		Layout            string `json:"layout"`              // "product_main" or "avatar_main"
	}
	c.BindJSON(&requestBody)

	var project models.Project
	if err := h.db.First(&project, "id = ?", projectID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Project not found"})
		return
	}

	// Update status
	project.Status = "video_generating"
	h.db.Save(&project)

	// ALWAYS use Gemini-generated script - no custom script override
	if project.GeneratedScript == "" {
		project.Status = "uploaded" // Revert status
		h.db.Save(&project)
		c.JSON(400, gin.H{
			"error": "No Gemini-generated script found. Please ensure product description was processed correctly.",
		})
		return
	}

	fmt.Printf("\n" + strings.Repeat("=", 60) + "\n")
	fmt.Printf("🎬 USING GEMINI-GENERATED SCRIPT FOR VIDEO\n")
	fmt.Printf(strings.Repeat("=", 60) + "\n")
	fmt.Printf("📝 Script: \"%s\"\n", project.GeneratedScript)
	fmt.Printf(strings.Repeat("=", 60) + "\n\n")

	// Generate video using AI service with Gemini-generated script
	videoPath, err := h.aiService.GenerateVideo(
		project.ProductImagePath,
		project.PersonMediaPath,
		project.PersonMediaType,
		project.GeneratedScript, // ALWAYS use Gemini script
		requestBody.ProductVideoStyle,
		requestBody.Layout,
	)
	if err != nil {
		project.Status = "uploaded" // Revert status
		h.db.Save(&project)
		c.JSON(500, gin.H{"error": "Failed to generate video", "details": err.Error()})
		return
	}

	// Update project with generated video path
	project.GeneratedVideoPath = videoPath
	project.Status = "video_complete"
	h.db.Save(&project)

	c.JSON(200, gin.H{
		"project_id": project.ID,
		"video_path": videoPath,
		"status":     project.Status,
	})
}

// GenerateWebsite generates a website for the product
func (h *Handlers) GenerateWebsite(c *gin.Context) {
	projectID := c.Param("id")

	var project models.Project
	if err := h.db.First(&project, "id = ?", projectID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Project not found"})
		return
	}

	// Update status
	project.Status = "website_generating"
	h.db.Save(&project)

	// Generate website
	websitePath, err := h.aiService.GenerateWebsite(project)
	if err != nil {
		project.Status = "video_complete" // Revert status
		h.db.Save(&project)
		c.JSON(500, gin.H{"error": "Failed to generate website", "details": err.Error()})
		return
	}

	// Update project
	project.WebsitePath = websitePath
	project.Status = "website_complete"
	h.db.Save(&project)

	fmt.Printf("\n" + strings.Repeat("=", 60) + "\n")
	fmt.Printf("✅ WEBSITE GENERATION COMPLETE\n")
	fmt.Printf(strings.Repeat("=", 60) + "\n")
	fmt.Printf("📁 Website Path: %s\n", websitePath)
	fmt.Printf("🆔 Project ID: %s\n", project.ID)
	fmt.Printf("📊 Status: %s\n", project.Status)
	fmt.Printf(strings.Repeat("=", 60) + "\n\n")

	c.JSON(200, gin.H{
		"project_id":   project.ID,
		"website_path": websitePath,
		"status":       project.Status,
	})
}

// GetProjects lists all projects
func (h *Handlers) GetProjects(c *gin.Context) {
	var projects []models.Project
	if err := h.db.Order("created_at DESC").Find(&projects).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch projects"})
		return
	}

	c.JSON(200, gin.H{"projects": projects})
}

// GetProject gets a single project by ID
func (h *Handlers) GetProject(c *gin.Context) {
	projectID := c.Param("id")

	var project models.Project
	if err := h.db.First(&project, "id = ?", projectID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Project not found"})
		return
	}

	c.JSON(200, project)
}

// UploadToInstagram uploads the generated video to Instagram
func (h *Handlers) UploadToInstagram(c *gin.Context) {
	projectID := c.Param("id")

	// Parse request body for Instagram credentials and options
	var requestBody struct {
		InstagramAccessToken string `json:"instagram_access_token"`
		InstagramUserID      string `json:"instagram_user_id"`
		CustomCaption        string `json:"custom_caption"`
	}
	c.BindJSON(&requestBody)

	var project models.Project
	if err := h.db.First(&project, "id = ?", projectID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Project not found"})
		return
	}

	// Check if video exists
	if project.GeneratedVideoPath == "" {
		c.JSON(400, gin.H{"error": "No generated video found. Please generate video first."})
		return
	}

	// Validate Instagram credentials
	accessToken := requestBody.InstagramAccessToken
	if accessToken == "" {
		// Try to get from environment variable
		accessToken = os.Getenv("INSTAGRAM_ACCESS_TOKEN")
	}
	if accessToken == "" {
		c.JSON(400, gin.H{"error": "Instagram access token is required"})
		return
	}

	instagramUserID := requestBody.InstagramUserID
	if instagramUserID == "" {
		instagramUserID = os.Getenv("INSTAGRAM_USER_ID")
	}
	if instagramUserID == "" {
		c.JSON(400, gin.H{"error": "Instagram user ID is required"})
		return
	}

	// Update status
	project.Status = "instagram_uploading"
	h.db.Save(&project)

	// Generate caption
	caption := requestBody.CustomCaption
	if caption == "" {
		caption = services.GenerateInstagramCaption(
			project.ProductName,
			project.ProductDescription,
			project.ProductPrice,
		)
	}

	// Create Instagram service and upload
	instagramService := services.NewInstagramService(accessToken)
	postID, postURL, err := instagramService.UploadVideoToInstagram(
		project.GeneratedVideoPath,
		caption,
		instagramUserID,
	)
	if err != nil {
		project.Status = "video_complete" // Revert status
		h.db.Save(&project)
		c.JSON(500, gin.H{"error": "Failed to upload to Instagram", "details": err.Error()})
		return
	}

	// Update project with Instagram post details
	project.InstagramPostID = postID
	project.InstagramPostURL = postURL
	project.Status = "instagram_posted"
	h.db.Save(&project)

	c.JSON(200, gin.H{
		"project_id":          project.ID,
		"instagram_post_id":   postID,
		"instagram_post_url":  postURL,
		"caption":             caption,
		"status":              project.Status,
		"message":             "Video successfully posted to Instagram!",
	})
}
