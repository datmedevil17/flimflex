package repository

import (
	"errors"
	"time"

	"github.com/datmedevil17/micro-flex/services/streaming/internal/models"
	"gorm.io/gorm"
)

type StreamingRepository interface {
	CreateSession(session *models.StreamSession) error
	GetSessionByToken(token string) (*models.StreamSession, error)
	DeleteExpiredSessions() error
	SaveAnalytics(analytics *models.StreamAnalytics) error
	GetUserStreamingStats(userID string, from, to time.Time) ([]models.StreamAnalytics, error)
}

type streamingRepository struct {
	db *gorm.DB
}

func NewStreamingRepository(db *gorm.DB) StreamingRepository {
	return &streamingRepository{db: db}
}

func (r *streamingRepository) CreateSession(session *models.StreamSession) error {
	return r.db.Create(session).Error
}

func (r *streamingRepository) GetSessionByToken(token string) (*models.StreamSession, error) {
	var session models.StreamSession
	err := r.db.Where("token = ? AND expires_at > ?", token, time.Now()).First(&session).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid or expired session")
		}
		return nil, err
	}
	return &session, nil
}

func (r *streamingRepository) DeleteExpiredSessions() error {
	return r.db.Where("expires_at < ?", time.Now()).Delete(&models.StreamSession{}).Error
}

func (r *streamingRepository) SaveAnalytics(analytics *models.StreamAnalytics) error {
	return r.db.Create(analytics).Error
}

func (r *streamingRepository) GetUserStreamingStats(userID string, from, to time.Time) ([]models.StreamAnalytics, error) {
	var analytics []models.StreamAnalytics
	err := r.db.Where("user_id = ? AND watched_at BETWEEN ? AND ?", userID, from, to).
		Find(&analytics).Error
	if err != nil {
		return nil, err
	}
	return analytics, nil
}