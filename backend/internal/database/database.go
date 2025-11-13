package database

import (
	"os"
	"path/filepath"

	"github.com/dealshare/hacathon/backend/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Initialize(dbPath string) (*gorm.DB, error) {
	// Create directory if it doesn't exist
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&models.Project{})
}

