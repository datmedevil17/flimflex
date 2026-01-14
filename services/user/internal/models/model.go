package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Profile struct {
	ID               string    `gorm:"type:uuid;primaryKey" json:"id"`
	Email            string    `gorm:"uniqueIndex;not null" json:"email"`
	FullName         string    `gorm:"not null" json:"full_name"`
	AvatarURL        string    `json:"avatar_url"`
	SubscriptionPlan string    `gorm:"default:'free'" json:"subscription_plan"`
	IsActive         bool      `gorm:"default:true" json:"is_active"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func (Profile) TableName() string {
	return "users"
}

type WatchHistory struct {
	ID            string    `gorm:"type:uuid;primaryKey" json:"id"`
	UserID        string    `gorm:"type:uuid;not null;index" json:"user_id"`
	MovieID       string    `gorm:"type:uuid;not null;index" json:"movie_id"`
	WatchDuration int       `gorm:"not null" json:"watch_duration"` // in seconds
	WatchedAt     time.Time `gorm:"not null;index" json:"watched_at"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Watchlist struct {
	ID        string    `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    string    `gorm:"type:uuid;not null;index" json:"user_id"`
	MovieID   string    `gorm:"type:uuid;not null;index" json:"movie_id"`
	AddedAt   time.Time `gorm:"not null" json:"added_at"`
	CreatedAt time.Time `json:"created_at"`
}

func (p *Profile) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	return nil
}

func (wh *WatchHistory) BeforeCreate(tx *gorm.DB) error {
	if wh.ID == "" {
		wh.ID = uuid.New().String()
	}
	if wh.WatchedAt.IsZero() {
		wh.WatchedAt = time.Now()
	}
	return nil
}

func (w *Watchlist) BeforeCreate(tx *gorm.DB) error {
	if w.ID == "" {
		w.ID = uuid.New().String()
	}
	if w.AddedAt.IsZero() {
		w.AddedAt = time.Now()
	}
	return nil
}
