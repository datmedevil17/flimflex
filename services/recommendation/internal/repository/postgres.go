package repository

import (
	"errors"

	"github.com/datmedevil17/micro-flex/services/recommendation/internal/models"
	"gorm.io/gorm"
)

type RecommendationRepository interface {
	GetUserPreference(userID string) (*models.UserPreference, error)
	SaveUserPreference(pref *models.UserPreference) error
	GetTrendingMovies(timeRange string, limit int) ([]models.TrendingCache, error)
	SaveTrendingCache(cache *models.TrendingCache) error
	UpdateTrendingCache(cache *models.TrendingCache) error
}

type recommendationRepository struct {
	db *gorm.DB
}

func NewRecommendationRepository(db *gorm.DB) RecommendationRepository {
	return &recommendationRepository{db: db}
}

func (r *recommendationRepository) GetUserPreference(userID string) (*models.UserPreference, error) {
	var pref models.UserPreference
	err := r.db.Where("user_id = ?", userID).First(&pref).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user preference not found")
		}
		return nil, err
	}
	return &pref, nil
}

func (r *recommendationRepository) SaveUserPreference(pref *models.UserPreference) error {
	return r.db.Save(pref).Error
}

func (r *recommendationRepository) GetTrendingMovies(timeRange string, limit int) ([]models.TrendingCache, error) {
	var trending []models.TrendingCache
	err := r.db.Where("time_range = ?", timeRange).
		Order("trending_score DESC").
		Limit(limit).
		Find(&trending).Error
	if err != nil {
		return nil, err
	}
	return trending, nil
}

func (r *recommendationRepository) SaveTrendingCache(cache *models.TrendingCache) error {
	return r.db.Create(cache).Error
}

func (r *recommendationRepository) UpdateTrendingCache(cache *models.TrendingCache) error {
	return r.db.Save(cache).Error
}