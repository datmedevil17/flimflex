package service

import (
	"context"

	"time"

	"github.com/datmedevil17/micro-flex/pkg/jwt"
	pb "github.com/datmedevil17/micro-flex/proto/auth"
	"github.com/datmedevil17/micro-flex/services/auth/internal/models"
	"github.com/datmedevil17/micro-flex/services/auth/internal/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthService struct {
	pb.UnimplementedAuthServiceServer
	repo       repository.AuthRepository
	jwtManager *jwt.JWTManager
}

func NewAuthService(repo repository.AuthRepository, jwtManager *jwt.JWTManager) *AuthService {
	return &AuthService{
		repo:       repo,
		jwtManager: jwtManager,
	}
}

func (s *AuthService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	// Validate input
	if req.Email == "" || req.Password == "" || req.FullName == "" {
		return nil, status.Error(codes.InvalidArgument, "email, password, and full name are required")
	}

	// Check if user already exists
	existingUser, _ := s.repo.GetUserByEmail(req.Email)
	if existingUser != nil {
		return nil, status.Error(codes.AlreadyExists, "user with this email already exists")
	}

	// Create new user
	user := &models.User{
		Email:    req.Email,
		FullName: req.FullName,
	}

	if err := user.HashPassword(req.Password); err != nil {
		return nil, status.Error(codes.Internal, "failed to hash password")
	}

	if err := s.repo.CreateUser(user); err != nil {
		return nil, status.Error(codes.Internal, "failed to create user")
	}

	// Generate tokens
	accessToken, err := s.jwtManager.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate access token")
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID, user.Email)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate refresh token")
	}

	// Save refresh token
	rt := &models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	if err := s.repo.SaveRefreshToken(rt); err != nil {
		return nil, status.Error(codes.Internal, "failed to save refresh token")
	}

	return &pb.RegisterResponse{
		UserId:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Message:      "User registered successfully",
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	// Validate input
	if req.Email == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "email and password are required")
	}

	// Get user by email
	user, err := s.repo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, status.Error(codes.NotFound, "invalid email or password")
	}

	// Check password
	if !user.CheckPassword(req.Password) {
		return nil, status.Error(codes.Unauthenticated, "invalid email or password")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, status.Error(codes.PermissionDenied, "account is inactive")
	}

	// Generate tokens
	accessToken, err := s.jwtManager.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate access token")
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID, user.Email)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate refresh token")
	}

	// Save refresh token
	rt := &models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	if err := s.repo.SaveRefreshToken(rt); err != nil {
		return nil, status.Error(codes.Internal, "failed to save refresh token")
	}

	return &pb.LoginResponse{
		UserId:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Message:      "Login successful",
	}, nil
}

func (s *AuthService) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	if req.Token == "" {
		return nil, status.Error(codes.InvalidArgument, "token is required")
	}

	claims, err := s.jwtManager.ValidateToken(req.Token)
	if err != nil {
		return &pb.ValidateTokenResponse{
			Valid: false,
		}, nil
	}

	return &pb.ValidateTokenResponse{
		Valid:  true,
		UserId: claims.UserID,
		Email:  claims.Email,
	}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	if req.RefreshToken == "" {
		return nil, status.Error(codes.InvalidArgument, "refresh token is required")
	}

	// Validate refresh token
	rt, err := s.repo.GetRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid or expired refresh token")
	}

	// Get user
	user, err := s.repo.GetUserByID(rt.UserID)
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	// Generate new tokens
	accessToken, err := s.jwtManager.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate access token")
	}

	newRefreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID, user.Email)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate refresh token")
	}

	// Delete old refresh token
	s.repo.DeleteRefreshToken(req.RefreshToken)

	// Save new refresh token
	newRT := &models.RefreshToken{
		UserID:    user.ID,
		Token:     newRefreshToken,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	if err := s.repo.SaveRefreshToken(newRT); err != nil {
		return nil, status.Error(codes.Internal, "failed to save refresh token")
	}

	return &pb.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
