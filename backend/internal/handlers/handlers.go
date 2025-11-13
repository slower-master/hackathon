package handlers

import (
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

	// Create project record
	project := &models.Project{
		ProductImagePath: productPath,
		PersonMediaPath:  personPath,
		PersonMediaType:  personMediaType,
		Status:           "uploaded",
	}

	if err := h.db.Create(project).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to create project"})
		return
	}

	c.JSON(201, gin.H{
		"project_id": project.ID,
		"status":     project.Status,
		"message":    "Files uploaded successfully",
	})
}

// GenerateVideo generates a promotional video from uploaded media
func (h *Handlers) GenerateVideo(c *gin.Context) {
	projectID := c.Param("id")

	// Parse request body for custom script and video options
	var requestBody struct {
		Script            string `json:"script"`
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

	// Generate video using AI service with custom script and options
	videoPath, err := h.aiService.GenerateVideo(
		project.ProductImagePath,
		project.PersonMediaPath,
		project.PersonMediaType,
		requestBody.Script,
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
