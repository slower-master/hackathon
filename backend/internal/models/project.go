package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Project struct {
	ID                 string    `json:"id" gorm:"primaryKey"`
	ProductImagePath   string    `json:"product_image_path"`
	PersonMediaPath    string    `json:"person_media_path"`
	PersonMediaType    string    `json:"person_media_type"` // "image" or "video"
	GeneratedVideoPath string    `json:"generated_video_path,omitempty"`
	WebsitePath        string    `json:"website_path,omitempty"`
	WebsiteURL         string    `json:"website_url,omitempty"`
	Status             string    `json:"status"` // "uploaded", "video_generating", "video_complete", "website_generating", "website_complete", "deployed"
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

func (p *Project) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	return nil
}

