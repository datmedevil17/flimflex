package service

import (
	"context"

	pb "github.com/datmedevil17/micro-flex/proto/user"
	"github.com/datmedevil17/micro-flex/services/user/internal/models"
	"github.com/datmedevil17/micro-flex/services/user/internal/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) GetProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.GetProfileResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	profile, err := s.repo.GetProfileByID(req.UserId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "profile not found")
	}

	return &pb.GetProfileResponse{
		Profile: &pb.Profile{
			UserId:           profile.ID,
			Email:            profile.Email,
			FullName:         profile.FullName,
			AvatarUrl:        profile.AvatarURL,
			SubscriptionPlan: profile.SubscriptionPlan,
			CreatedAt:        profile.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
	}, nil
}

func (s *UserService) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.UpdateProfileResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	profile, err := s.repo.GetProfileByID(req.UserId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "profile not found")
	}

	// Update fields if provided
	if req.FullName != "" {
		profile.FullName = req.FullName
	}
	if req.AvatarUrl != "" {
		profile.AvatarURL = req.AvatarUrl
	}

	if err := s.repo.UpdateProfile(profile); err != nil {
		return nil, status.Error(codes.Internal, "failed to update profile")
	}

	return &pb.UpdateProfileResponse{
		Profile: &pb.Profile{
			UserId:           profile.ID,
			Email:            profile.Email,
			FullName:         profile.FullName,
			AvatarUrl:        profile.AvatarURL,
			SubscriptionPlan: profile.SubscriptionPlan,
			CreatedAt:        profile.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
		Message: "Profile updated successfully",
	}, nil
}

func (s *UserService) GetWatchHistory(ctx context.Context, req *pb.GetWatchHistoryRequest) (*pb.GetWatchHistoryResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	limit := int(req.Limit)
	offset := int(req.Offset)

	if limit <= 0 {
		limit = 20
	}

	histories, total, err := s.repo.GetWatchHistory(req.UserId, limit, offset)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get watch history")
	}

	var items []*pb.WatchHistory
	for _, h := range histories {
		items = append(items, &pb.WatchHistory{
			Id:            h.ID,
			UserId:        h.UserID,
			MovieId:       h.MovieID,
			WatchDuration: int32(h.WatchDuration),
			WatchedAt:     h.WatchedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return &pb.GetWatchHistoryResponse{
		Items: items,
		Total: int32(total),
	}, nil
}

func (s *UserService) AddToWatchHistory(ctx context.Context, req *pb.AddToWatchHistoryRequest) (*pb.AddToWatchHistoryResponse, error) {
	if req.UserId == "" || req.MovieId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id and movie_id are required")
	}

	history := &models.WatchHistory{
		UserID:        req.UserId,
		MovieID:       req.MovieId,
		WatchDuration: int(req.WatchDuration),
	}

	if err := s.repo.AddWatchHistory(history); err != nil {
		return nil, status.Error(codes.Internal, "failed to add watch history")
	}

	return &pb.AddToWatchHistoryResponse{
		Message: "Watch history added successfully",
	}, nil
}

func (s *UserService) GetWatchlist(ctx context.Context, req *pb.GetWatchlistRequest) (*pb.GetWatchlistResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	watchlist, err := s.repo.GetWatchlist(req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get watchlist")
	}

	var items []*pb.Watchlist
	for _, w := range watchlist {
		items = append(items, &pb.Watchlist{
			Id:      w.ID,
			UserId:  w.UserID,
			MovieId: w.MovieID,
			AddedAt: w.AddedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return &pb.GetWatchlistResponse{
		Items: items,
	}, nil
}

func (s *UserService) AddToWatchlist(ctx context.Context, req *pb.AddToWatchlistRequest) (*pb.AddToWatchlistResponse, error) {
	if req.UserId == "" || req.MovieId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id and movie_id are required")
	}

	watchlist := &models.Watchlist{
		UserID:  req.UserId,
		MovieID: req.MovieId,
	}

	if err := s.repo.AddToWatchlist(watchlist); err != nil {
		if err.Error() == "movie already in watchlist" {
			return nil, status.Error(codes.AlreadyExists, "movie already in watchlist")
		}
		return nil, status.Error(codes.Internal, "failed to add to watchlist")
	}

	return &pb.AddToWatchlistResponse{
		Message: "Added to watchlist successfully",
	}, nil
}

func (s *UserService) RemoveFromWatchlist(ctx context.Context, req *pb.RemoveFromWatchlistRequest) (*pb.RemoveFromWatchlistResponse, error) {
	if req.UserId == "" || req.MovieId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id and movie_id are required")
	}

	if err := s.repo.RemoveFromWatchlist(req.UserId, req.MovieId); err != nil {
		if err.Error() == "item not found in watchlist" {
			return nil, status.Error(codes.NotFound, "item not found in watchlist")
		}
		return nil, status.Error(codes.Internal, "failed to remove from watchlist")
	}

	return &pb.RemoveFromWatchlistResponse{
		Message: "Removed from watchlist successfully",
	}, nil
}
