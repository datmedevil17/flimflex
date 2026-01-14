package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Movie struct {
	ID           string         `gorm:"type:uuid;primaryKey" json:"id"`
	Title        string         `gorm:"not null;index" json:"title"`
	Description  string         `gorm:"type:text" json:"description"`
	ThumbnailURL string         `json:"thumbnail_url"`
	VideoURL     string         `json:"video_url"`
	Duration     int            `gorm:"not null" json:"duration"` // in seconds
	Genre        string         `gorm:"not null;index" json:"genre"`
	Rating       float64        `gorm:"default:0" json:"rating"`
	ReleaseYear  int            `gorm:"not null;index" json:"release_year"`
	Director     string         `json:"director"`
	Cast         pq.StringArray `gorm:"type:text[]" json:"cast"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

func (m *Movie) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return nil
}