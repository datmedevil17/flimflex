package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
	pb "github.com/datmedevil17/micro-flex/proto/streaming"
	"github.com/datmedevil17/micro-flex/services/streaming/internal/models"
	"github.com/datmedevil17/micro-flex/services/streaming/internal/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type StreamingService struct {
	pb.UnimplementedStreamingServiceServer
	repo      repository.StreamingRepository
	secretKey string
}

func NewStreamingService(repo repository.StreamingRepository, secretKey string) *StreamingService {
	return &StreamingService{
		repo:      repo,
		secretKey: secretKey,
	}
}

func (s *StreamingService) GetStreamingURL(ctx context.Context, req *pb.GetStreamingURLRequest) (*pb.GetStreamingURLResponse, error) {
	if req.MovieId == "" || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "movie_id and user_id are required")
	}

	// Validate quality
	quality := req.Quality
	if quality == 0 {
		quality = 720 // Default to 720p
	}

	validQualities := map[int32]bool{360: true, 480: true, 720: true, 1080: true}
	if !validQualities[quality] {
		return nil, status.Error(codes.InvalidArgument, "invalid quality, must be 360, 480, 720, or 1080")
	}

	// Generate secure streaming token
	token := s.generateStreamToken(req.UserId, req.MovieId)
	expiresIn := 3600 // 1 hour

	// Create session
	session := &models.StreamSession{
		UserID:    req.UserId,
		MovieID:   req.MovieId,
		Quality:   int(quality),
		Token:     token,
		ExpiresAt: time.Now().Add(time.Duration(expiresIn) * time.Second),
	}

	if err := s.repo.CreateSession(session); err != nil {
		return nil, status.Error(codes.Internal, "failed to create streaming session")
	}

	// Generate streaming URL with token
	streamingURL := fmt.Sprintf("https://cdn.netflix.com/stream/%s?token=%s&quality=%d", req.MovieId, token, quality)

	return &pb.GetStreamingURLResponse{
		StreamingUrl: streamingURL,
		ExpiresIn:    int32(expiresIn),
		Message:      "Streaming URL generated successfully",
	}, nil
}

func (s *StreamingService) ValidateStreamAccess(ctx context.Context, req *pb.ValidateStreamAccessRequest) (*pb.ValidateStreamAccessResponse, error) {
	if req.UserId == "" || req.MovieId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id and movie_id are required")
	}

	// In a real implementation, you would check:
	// 1. User's subscription status
	// 2. Geographic restrictions
	// 3. Concurrent stream limits
	// 4. Parental controls

	// For now, we'll do basic validation
	hasAccess := true
	message := "Access granted"

	// Example: Check if user has active subscription (mock)
	// This would typically call the Auth or User service
	// For now, we always grant access
	
	return &pb.ValidateStreamAccessResponse{
		HasAccess: hasAccess,
		Message:   message,
	}, nil
}

// Helper function to generate secure stream token
func (s *StreamingService) generateStreamToken(userID, movieID string) string {
	data := fmt.Sprintf("%s:%s:%s:%d", userID, movieID, s.secretKey, time.Now().Unix())
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:]) + "-" + uuid.New().String()
}