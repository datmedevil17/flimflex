package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserPreference struct {
	ID           string    `gorm:"type:uuid;primaryKey" json:"id"`
	UserID       string    `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`
	FavoriteGenres string  `gorm:"type:text" json:"favorite_genres"` // JSON array
	WatchedGenres  string  `gorm:"type:text" json:"watched_genres"`  // JSON map with counts
	AverageRating  float64 `json:"average_rating"`
	UpdatedAt      time.Time `json:"updated_at"`
	CreatedAt      time.Time `json:"created_at"`
}

type TrendingCache struct {
	ID           string    `gorm:"type:uuid;primaryKey" json:"id"`
	MovieID      string    `gorm:"type:uuid;not null;index" json:"movie_id"`
	TimeRange    string    `gorm:"not null;index" json:"time_range"` // day, week, month
	ViewCount    int       `gorm:"not null" json:"view_count"`
	TrendingScore float64  `gorm:"not null" json:"trending_score"`
	UpdatedAt    time.Time `json:"updated_at"`
	CreatedAt    time.Time `json:"created_at"`
}

func (up *UserPreference) BeforeCreate(tx *gorm.DB) error {
	if up.ID == "" {
		up.ID = uuid.New().String()
	}
	return nil
}

func (tc *TrendingCache) BeforeCreate(tx *gorm.DB) error {
	if tc.ID == "" {
		tc.ID = uuid.New().String()
	}
	return nil
}