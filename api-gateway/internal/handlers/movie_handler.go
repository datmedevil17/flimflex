package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/datmedevil17/micro-flex/pkg/response"
	moviepb "github.com/datmedevil17/micro-flex/proto/movie"
	"github.com/gin-gonic/gin"
)

type MovieHandler struct {
	movieClient moviepb.MovieServiceClient
}

func NewMovieHandler(movieClient moviepb.MovieServiceClient) *MovieHandler {
	return &MovieHandler{
		movieClient: movieClient,
	}
}

type CreateMovieRequest struct {
	Title        string   `json:"title" binding:"required"`
	Description  string   `json:"description" binding:"required"`
	ThumbnailURL string   `json:"thumbnail_url"`
	VideoURL     string   `json:"video_url"`
	Duration     int32    `json:"duration"`
	Genre        string   `json:"genre" binding:"required"`
	ReleaseYear  int32    `json:"release_year"`
	Director     string   `json:"director"`
	Cast         []string `json:"cast"`
}

func (h *MovieHandler) CreateMovie(c *gin.Context) {
	var req CreateMovieRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	resp, err := h.movieClient.CreateMovie(context.Background(), &moviepb.CreateMovieRequest{
		Title:        req.Title,
		Description:  req.Description,
		ThumbnailUrl: req.ThumbnailURL,
		VideoUrl:     req.VideoURL,
		Duration:     req.Duration,
		Genre:        req.Genre,
		ReleaseYear:  req.ReleaseYear,
		Director:     req.Director,
		Cast:         req.Cast,
	})

	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, http.StatusOK, resp.Message, resp.Movie)
}

func (h *MovieHandler) UpdateMovie(c *gin.Context) {
	id := c.Param("id")
	var req CreateMovieRequest // Reusing struct as fields are same
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	resp, err := h.movieClient.UpdateMovie(context.Background(), &moviepb.UpdateMovieRequest{
		Id:           id,
		Title:        req.Title,
		Description:  req.Description,
		ThumbnailUrl: req.ThumbnailURL,
		VideoUrl:     req.VideoURL,
		Duration:     req.Duration,
		Genre:        req.Genre,
		ReleaseYear:  req.ReleaseYear,
		Director:     req.Director,
		Cast:         req.Cast,
	})

	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, http.StatusOK, resp.Message, resp.Movie)
}

func (h *MovieHandler) DeleteMovie(c *gin.Context) {
	id := c.Param("id")

	resp, err := h.movieClient.DeleteMovie(context.Background(), &moviepb.DeleteMovieRequest{
		Id: id,
	})

	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, http.StatusOK, resp.Message, nil)
}

func (h *MovieHandler) ListMovies(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	resp, err := h.movieClient.ListMovies(context.Background(), &moviepb.ListMoviesRequest{
		Limit:  int32(limit),
		Offset: int32(offset),
	})

	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Movies retrieved successfully", gin.H{
		"movies": resp.Movies,
		"total":  resp.Total,
	})
}

func (h *MovieHandler) GetMovie(c *gin.Context) {
	id := c.Param("id")

	resp, err := h.movieClient.GetMovie(context.Background(), &moviepb.GetMovieRequest{
		Id: id,
	})

	if err != nil {
		response.NotFound(c, "Movie not found")
		return
	}

	response.Success(c, http.StatusOK, "Movie retrieved successfully", resp.Movie)
}

func (h *MovieHandler) SearchMovies(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		response.Error(c, http.StatusBadRequest, "Query parameter 'q' is required")
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	resp, err := h.movieClient.SearchMovies(context.Background(), &moviepb.SearchMoviesRequest{
		Query:  query,
		Limit:  int32(limit),
		Offset: int32(offset),
	})

	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Search results", gin.H{
		"movies": resp.Movies,
		"total":  resp.Total,
	})
}

func (h *MovieHandler) GetMoviesByGenre(c *gin.Context) {
	genre := c.Param("genre")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	resp, err := h.movieClient.GetMoviesByGenre(context.Background(), &moviepb.GetMoviesByGenreRequest{
		Genre:  genre,
		Limit:  int32(limit),
		Offset: int32(offset),
	})

	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Movies retrieved successfully", gin.H{
		"movies": resp.Movies,
		"total":  resp.Total,
	})
}
