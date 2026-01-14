package repository

import (
	"errors"

	"github.com/datmedevil17/micro-flex/services/user/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	// Profile operations
	CreateProfile(profile *models.Profile) error
	GetProfileByID(userID string) (*models.Profile, error)
	GetProfileByEmail(email string) (*models.Profile, error)
	UpdateProfile(profile *models.Profile) error
	
	// Watch history operations
	AddWatchHistory(history *models.WatchHistory) error
	GetWatchHistory(userID string, limit, offset int) ([]models.WatchHistory, int64, error)
	UpdateWatchHistory(history *models.WatchHistory) error
	
	// Watchlist operations
	AddToWatchlist(watchlist *models.Watchlist) error
	GetWatchlist(userID string) ([]models.Watchlist, error)
	RemoveFromWatchlist(userID, movieID string) error
	IsInWatchlist(userID, movieID string) (bool, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// Profile operations
func (r *userRepository) CreateProfile(profile *models.Profile) error {
	return r.db.Create(profile).Error
}

func (r *userRepository) GetProfileByID(userID string) (*models.Profile, error) {
	var profile models.Profile
	err := r.db.Where("id = ?", userID).First(&profile).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("profile not found")
		}
		return nil, err
	}
	return &profile, nil
}

func (r *userRepository) GetProfileByEmail(email string) (*models.Profile, error) {
	var profile models.Profile
	err := r.db.Where("email = ?", email).First(&profile).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("profile not found")
		}
		return nil, err
	}
	return &profile, nil
}

func (r *userRepository) UpdateProfile(profile *models.Profile) error {
	return r.db.Save(profile).Error
}

// Watch history operations
func (r *userRepository) AddWatchHistory(history *models.WatchHistory) error {
	// Check if record already exists for this user and movie
	var existing models.WatchHistory
	err := r.db.Where("user_id = ? AND movie_id = ?", history.UserID, history.MovieID).First(&existing).Error
	
	if err == nil {
		// Update existing record
		existing.WatchDuration = history.WatchDuration
		existing.WatchedAt = history.WatchedAt
		return r.db.Save(&existing).Error
	}
	
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Create new record
		return r.db.Create(history).Error
	}
	
	return err
}

func (r *userRepository) GetWatchHistory(userID string, limit, offset int) ([]models.WatchHistory, int64, error) {
	var histories []models.WatchHistory
	var total int64
	
	query := r.db.Where("user_id = ?", userID)
	
	if err := query.Model(&models.WatchHistory{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	err := query.Order("watched_at DESC").Limit(limit).Offset(offset).Find(&histories).Error
	if err != nil {
		return nil, 0, err
	}
	
	return histories, total, nil
}

func (r *userRepository) UpdateWatchHistory(history *models.WatchHistory) error {
	return r.db.Save(history).Error
}

// Watchlist operations
func (r *userRepository) AddToWatchlist(watchlist *models.Watchlist) error {
	// Check if already in watchlist
	exists, err := r.IsInWatchlist(watchlist.UserID, watchlist.MovieID)
	if err != nil {
		return err
	}
	
	if exists {
		return errors.New("movie already in watchlist")
	}
	
	return r.db.Create(watchlist).Error
}

func (r *userRepository) GetWatchlist(userID string) ([]models.Watchlist, error) {
	var watchlist []models.Watchlist
	err := r.db.Where("user_id = ?", userID).Order("added_at DESC").Find(&watchlist).Error
	if err != nil {
		return nil, err
	}
	return watchlist, nil
}

func (r *userRepository) RemoveFromWatchlist(userID, movieID string) error {
	result := r.db.Where("user_id = ? AND movie_id = ?", userID, movieID).Delete(&models.Watchlist{})
	if result.Error != nil {
		return result.Error
	}
	
	if result.RowsAffected == 0 {
		return errors.New("item not found in watchlist")
	}
	
	return nil
}

func (r *userRepository) IsInWatchlist(userID, movieID string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Watchlist{}).Where("user_id = ? AND movie_id = ?", userID, movieID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}