package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/datmedevil17/micro-flex/pkg/response"
	recommendationpb "github.com/datmedevil17/micro-flex/proto/recommendation"
	"github.com/gin-gonic/gin"
)

type RecommendationHandler struct {
	recommendationClient recommendationpb.RecommendationServiceClient
}

func NewRecommendationHandler(recommendationClient recommendationpb.RecommendationServiceClient) *RecommendationHandler {
	return &RecommendationHandler{
		recommendationClient: recommendationClient,
	}
}

func (h *RecommendationHandler) GetRecommendations(c *gin.Context) {
	userID := c.GetString("user_id")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	resp, err := h.recommendationClient.GetRecommendations(context.Background(), &recommendationpb.GetRecommendationsRequest{
		UserId: userID,
		Limit:  int32(limit),
	})

	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Recommendations retrieved successfully", resp.Recommendations)
}
