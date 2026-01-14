package repository

import (
	"errors"
	"strings"

	"github.com/datmedevil17/micro-flex/services/movie/internal/models"
	"gorm.io/gorm"
)

type MovieRepository interface {
	CreateMovie(movie *models.Movie) error
	GetMovieByID(id string) (*models.Movie, error)
	ListMovies(limit, offset int) ([]models.Movie, int64, error)
	UpdateMovie(movie *models.Movie) error
	DeleteMovie(id string) error
	SearchMovies(query string, limit, offset int) ([]models.Movie, int64, error)
	GetMoviesByGenre(genre string, limit, offset int) ([]models.Movie, int64, error)
}

type movieRepository struct {
	db *gorm.DB
}

func NewMovieRepository(db *gorm.DB) MovieRepository {
	return &movieRepository{db: db}
}

func (r *movieRepository) CreateMovie(movie *models.Movie) error {
	return r.db.Create(movie).Error
}

func (r *movieRepository) GetMovieByID(id string) (*models.Movie, error) {
	var movie models.Movie
	err := r.db.Where("id = ?", id).First(&movie).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("movie not found")
		}
		return nil, err
	}
	return &movie, nil
}

func (r *movieRepository) ListMovies(limit, offset int) ([]models.Movie, int64, error) {
	var movies []models.Movie
	var total int64
	
	if err := r.db.Model(&models.Movie{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	err := r.db.Order("created_at DESC").Limit(limit).Offset(offset).Find(&movies).Error
	if err != nil {
		return nil, 0, err
	}
	
	return movies, total, nil
}

func (r *movieRepository) UpdateMovie(movie *models.Movie) error {
	return r.db.Save(movie).Error
}

func (r *movieRepository) DeleteMovie(id string) error {
	result := r.db.Where("id = ?", id).Delete(&models.Movie{})
	if result.Error != nil {
		return result.Error
	}
	
	if result.RowsAffected == 0 {
		return errors.New("movie not found")
	}
	
	return nil
}

func (r *movieRepository) SearchMovies(query string, limit, offset int) ([]models.Movie, int64, error) {
	var movies []models.Movie
	var total int64
	
	searchQuery := "%" + strings.ToLower(query) + "%"
	
	dbQuery := r.db.Where("LOWER(title) LIKE ? OR LOWER(description) LIKE ? OR LOWER(director) LIKE ?", 
		searchQuery, searchQuery, searchQuery)
	
	if err := dbQuery.Model(&models.Movie{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	err := dbQuery.Order("created_at DESC").Limit(limit).Offset(offset).Find(&movies).Error
	if err != nil {
		return nil, 0, err
	}
	
	return movies, total, nil
}

func (r *movieRepository) GetMoviesByGenre(genre string, limit, offset int) ([]models.Movie, int64, error) {
	var movies []models.Movie
	var total int64
	
	query := r.db.Where("LOWER(genre) = ?", strings.ToLower(genre))
	
	if err := query.Model(&models.Movie{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&movies).Error
	if err != nil {
		return nil, 0, err
	}
	
	return movies, total, nil
}