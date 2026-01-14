package repository

import (
	"github.com/datmedevil17/micro-flex/services/upload/internal/models"
	"gorm.io/gorm"
)

type UploadRepository interface {
	SaveFile(file *models.File) error
	GetFileByPublicID(publicID string) (*models.File, error)
}

type uploadRepository struct {
	db *gorm.DB
}

func NewUploadRepository(db *gorm.DB) UploadRepository {
	return &uploadRepository{db: db}
}

func (r *uploadRepository) SaveFile(file *models.File) error {
	return r.db.Create(file).Error
}

func (r *uploadRepository) GetFileByPublicID(publicID string) (*models.File, error) {
	var file models.File
	err := r.db.Where("public_id = ?", publicID).First(&file).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}
