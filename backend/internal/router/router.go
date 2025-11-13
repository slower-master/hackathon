package router

import (
	"github.com/dealshare/hacathon/backend/internal/handlers"
	"github.com/gin-gonic/gin"
)

func Setup(h *handlers.Handlers) *gin.Engine {
	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// API routes
	api := r.Group("/api/v1")
	{
		api.POST("/upload", h.UploadMedia)
		api.GET("/projects", h.GetProjects)
		api.GET("/projects/:id", h.GetProject)
		api.POST("/projects/:id/generate-video", h.GenerateVideo)
		api.POST("/projects/:id/generate-website", h.GenerateWebsite)
	}

	// Serve static files (generated videos and websites)
	r.Static("/static/uploads", "./uploads")
	r.Static("/static/generated/videos", "./generated/videos")
	r.Static("/static/generated/websites", "./generated/websites")

	return r
}
