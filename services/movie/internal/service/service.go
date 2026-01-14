package service

import (
	"context"

	pb "github.com/datmedevil17/micro-flex/proto/movie"
	"github.com/datmedevil17/micro-flex/services/movie/internal/models"
	"github.com/datmedevil17/micro-flex/services/movie/internal/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MovieService struct {
	pb.UnimplementedMovieServiceServer
	repo repository.MovieRepository
}

func NewMovieService(repo repository.MovieRepository) *MovieService {
	return &MovieService{
		repo: repo,
	}
}

func (s *MovieService) CreateMovie(ctx context.Context, req *pb.CreateMovieRequest) (*pb.CreateMovieResponse, error) {
	if req.Title == "" {
		return nil, status.Error(codes.InvalidArgument, "title is required")
	}

	movie := &models.Movie{
		Title:        req.Title,
		Description:  req.Description,
		ThumbnailURL: req.ThumbnailUrl,
		VideoURL:     req.VideoUrl,
		Duration:     int(req.Duration),
		Genre:        req.Genre,
		ReleaseYear:  int(req.ReleaseYear),
		Director:     req.Director,
		Cast:         req.Cast,
	}

	if err := s.repo.CreateMovie(movie); err != nil {
		return nil, status.Error(codes.Internal, "failed to create movie")
	}

	return &pb.CreateMovieResponse{
		Movie:   s.movieToProto(movie),
		Message: "Movie created successfully",
	}, nil
}

func (s *MovieService) GetMovie(ctx context.Context, req *pb.GetMovieRequest) (*pb.GetMovieResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	movie, err := s.repo.GetMovieByID(req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "movie not found")
	}

	return &pb.GetMovieResponse{
		Movie: s.movieToProto(movie),
	}, nil
}

func (s *MovieService) ListMovies(ctx context.Context, req *pb.ListMoviesRequest) (*pb.ListMoviesResponse, error) {
	limit := int(req.Limit)
	offset := int(req.Offset)

	if limit <= 0 {
		limit = 20
	}

	movies, total, err := s.repo.ListMovies(limit, offset)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to list movies")
	}

	var pbMovies []*pb.Movie
	for _, movie := range movies {
		pbMovies = append(pbMovies, s.movieToProto(&movie))
	}

	return &pb.ListMoviesResponse{
		Movies: pbMovies,
		Total:  int32(total),
	}, nil
}

func (s *MovieService) UpdateMovie(ctx context.Context, req *pb.UpdateMovieRequest) (*pb.UpdateMovieResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	movie, err := s.repo.GetMovieByID(req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "movie not found")
	}

	// Update fields if provided
	if req.Title != "" {
		movie.Title = req.Title
	}
	if req.Description != "" {
		movie.Description = req.Description
	}
	if req.ThumbnailUrl != "" {
		movie.ThumbnailURL = req.ThumbnailUrl
	}
	if req.VideoUrl != "" {
		movie.VideoURL = req.VideoUrl
	}
	if req.Duration > 0 {
		movie.Duration = int(req.Duration)
	}
	if req.Genre != "" {
		movie.Genre = req.Genre
	}
	if req.Rating > 0 {
		movie.Rating = req.Rating
	}
	if req.ReleaseYear > 0 {
		movie.ReleaseYear = int(req.ReleaseYear)
	}
	if req.Director != "" {
		movie.Director = req.Director
	}
	if len(req.Cast) > 0 {
		movie.Cast = req.Cast
	}

	if err := s.repo.UpdateMovie(movie); err != nil {
		return nil, status.Error(codes.Internal, "failed to update movie")
	}

	return &pb.UpdateMovieResponse{
		Movie:   s.movieToProto(movie),
		Message: "Movie updated successfully",
	}, nil
}

func (s *MovieService) DeleteMovie(ctx context.Context, req *pb.DeleteMovieRequest) (*pb.DeleteMovieResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	if err := s.repo.DeleteMovie(req.Id); err != nil {
		if err.Error() == "movie not found" {
			return nil, status.Error(codes.NotFound, "movie not found")
		}
		return nil, status.Error(codes.Internal, "failed to delete movie")
	}

	return &pb.DeleteMovieResponse{
		Message: "Movie deleted successfully",
	}, nil
}

func (s *MovieService) SearchMovies(ctx context.Context, req *pb.SearchMoviesRequest) (*pb.SearchMoviesResponse, error) {
	if req.Query == "" {
		return nil, status.Error(codes.InvalidArgument, "query is required")
	}

	limit := int(req.Limit)
	offset := int(req.Offset)

	if limit <= 0 {
		limit = 20
	}

	movies, total, err := s.repo.SearchMovies(req.Query, limit, offset)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to search movies")
	}

	var pbMovies []*pb.Movie
	for _, movie := range movies {
		pbMovies = append(pbMovies, s.movieToProto(&movie))
	}

	return &pb.SearchMoviesResponse{
		Movies: pbMovies,
		Total:  int32(total),
	}, nil
}

func (s *MovieService) GetMoviesByGenre(ctx context.Context, req *pb.GetMoviesByGenreRequest) (*pb.GetMoviesByGenreResponse, error) {
	if req.Genre == "" {
		return nil, status.Error(codes.InvalidArgument, "genre is required")
	}

	limit := int(req.Limit)
	offset := int(req.Offset)

	if limit <= 0 {
		limit = 20
	}

	movies, total, err := s.repo.GetMoviesByGenre(req.Genre, limit, offset)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get movies by genre")
	}

	var pbMovies []*pb.Movie
	for _, movie := range movies {
		pbMovies = append(pbMovies, s.movieToProto(&movie))
	}

	return &pb.GetMoviesByGenreResponse{
		Movies: pbMovies,
		Total:  int32(total),
	}, nil
}

// Helper function to convert model to proto
func (s *MovieService) movieToProto(movie *models.Movie) *pb.Movie {
	return &pb.Movie{
		Id:           movie.ID,
		Title:        movie.Title,
		Description:  movie.Description,
		ThumbnailUrl: movie.ThumbnailURL,
		VideoUrl:     movie.VideoURL,
		Duration:     int32(movie.Duration),
		Genre:        movie.Genre,
		Rating:       movie.Rating,
		ReleaseYear:  int32(movie.ReleaseYear),
		Director:     movie.Director,
		Cast:         movie.Cast,
		CreatedAt:    movie.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    movie.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}