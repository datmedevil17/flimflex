package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StreamSession struct {
	ID        string    `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    string    `gorm:"type:uuid;not null;index" json:"user_id"`
	MovieID   string    `gorm:"type:uuid;not null;index" json:"movie_id"`
	Quality   int       `gorm:"not null" json:"quality"`
	Token     string    `gorm:"uniqueIndex;not null" json:"token"`
	ExpiresAt time.Time `gorm:"not null;index" json:"expires_at"`
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	CreatedAt time.Time `json:"created_at"`
}

type StreamAnalytics struct {
	ID            string    `gorm:"type:uuid;primaryKey" json:"id"`
	UserID        string    `gorm:"type:uuid;not null;index" json:"user_id"`
	MovieID       string    `gorm:"type:uuid;not null;index" json:"movie_id"`
	SessionID     string    `gorm:"type:uuid;index" json:"session_id"`
	Duration      int       `json:"duration"` // seconds watched
	Quality       int       `json:"quality"`
	BufferCount   int       `json:"buffer_count"`
	AverageBitrate float64  `json:"average_bitrate"`
	WatchedAt     time.Time `gorm:"not null;index" json:"watched_at"`
	CreatedAt     time.Time `json:"created_at"`
}

func (s *StreamSession) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	return nil
}

func (a *StreamAnalytics) BeforeCreate(tx *gorm.DB) error {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	if a.WatchedAt.IsZero() {
		a.WatchedAt = time.Now()
	}
	return nil
}