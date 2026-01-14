package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type File struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	PublicID  string         `gorm:"not null" json:"public_id"`
	SecureURL string         `gorm:"not null" json:"secure_url"`
	FileType  string         `gorm:"not null" json:"file_type"` // image, video
	Format    string         `json:"format"`
	Bytes     int            `json:"bytes"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
