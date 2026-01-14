package handlers

import (
	"context"
	"net/http"

	"github.com/datmedevil17/micro-flex/pkg/response"
	userpb "github.com/datmedevil17/micro-flex/proto/user"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userClient userpb.UserServiceClient
}

func NewUserHandler(userClient userpb.UserServiceClient) *UserHandler {
	return &UserHandler{
		userClient: userClient,
	}
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetString("user_id")

	resp, err := h.userClient.GetProfile(context.Background(), &userpb.GetProfileRequest{
		UserId: userID,
	})

	if err != nil {
		response.NotFound(c, "Profile not found")
		return
	}

	response.Success(c, http.StatusOK, "Profile retrieved successfully", resp.Profile)
}

type UpdateProfileRequest struct {
	FullName  string `json:"full_name"`
	AvatarURL string `json:"avatar_url"`
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetString("user_id")

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	resp, err := h.userClient.UpdateProfile(context.Background(), &userpb.UpdateProfileRequest{
		UserId:    userID,
		FullName:  req.FullName,
		AvatarUrl: req.AvatarURL,
	})

	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, http.StatusOK, resp.Message, resp.Profile)
}

func (h *UserHandler) GetWatchlist(c *gin.Context) {
	userID := c.GetString("user_id")

	resp, err := h.userClient.GetWatchlist(context.Background(), &userpb.GetWatchlistRequest{
		UserId: userID,
	})

	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Watchlist retrieved successfully", resp.Items)
}

type AddToWatchlistRequest struct {
	MovieID string `json:"movie_id" binding:"required"`
}

func (h *UserHandler) AddToWatchlist(c *gin.Context) {
	userID := c.GetString("user_id")

	var req AddToWatchlistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	resp, err := h.userClient.AddToWatchlist(context.Background(), &userpb.AddToWatchlistRequest{
		UserId:  userID,
		MovieId: req.MovieID,
	})

	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, resp.Message, nil)
}

func (h *UserHandler) RemoveFromWatchlist(c *gin.Context) {
	userID := c.GetString("user_id")
	movieID := c.Param("id")

	resp, err := h.userClient.RemoveFromWatchlist(context.Background(), &userpb.RemoveFromWatchlistRequest{
		UserId:  userID,
		MovieId: movieID,
	})

	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, resp.Message, nil)
}
