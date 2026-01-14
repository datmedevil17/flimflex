package service

import (
	"context"
	"fmt"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	pb "github.com/datmedevil17/micro-flex/proto/upload"
	"github.com/datmedevil17/micro-flex/services/upload/config"
	"github.com/datmedevil17/micro-flex/services/upload/internal/models"
	"github.com/datmedevil17/micro-flex/services/upload/internal/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UploadService struct {
	pb.UnimplementedUploadServiceServer
	cloudinary *cloudinary.Cloudinary
	repo       repository.UploadRepository
}

func NewUploadService(cfg config.CloudinaryConfig, repo repository.UploadRepository) (*UploadService, error) {
	cld, err := cloudinary.NewFromParams(cfg.CloudName, cfg.APIKey, cfg.APISecret)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize cloudinary: %w", err)
	}

	return &UploadService{
		cloudinary: cld,
		repo:       repo,
	}, nil
}

func (s *UploadService) UploadVideo(ctx context.Context, req *pb.UploadVideoRequest) (*pb.UploadVideoResponse, error) {
	if len(req.FileData) == 0 {
		return nil, status.Error(codes.InvalidArgument, "file data is required")
	}

	if req.FileName == "" {
		return nil, status.Error(codes.InvalidArgument, "file name is required")
	}

	// Upload video to Cloudinary
	uploadParams := uploader.UploadParams{
		ResourceType: "video",
		Folder:       "netflix/videos",
		PublicID:     fmt.Sprintf("movie_%s_%s", req.MovieId, req.FileName),
	}

	result, err := s.cloudinary.Upload.Upload(ctx, req.FileData, uploadParams)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to upload video: %v", err))
	}

	file := &models.File{
		PublicID:  result.PublicID,
		SecureURL: result.SecureURL,
		FileType:  "video",
		Format:    result.Format,
		Bytes:     result.Bytes,
	}

	if err := s.repo.SaveFile(file); err != nil {
		// Log error but don't fail the request since upload succeeded
		fmt.Printf("failed to save file metadata: %v\n", err)
	}

	return &pb.UploadVideoResponse{
		Url:      result.SecureURL,
		PublicId: result.PublicID,
		Message:  "Video uploaded successfully",
	}, nil
}

func (s *UploadService) UploadThumbnail(ctx context.Context, req *pb.UploadThumbnailRequest) (*pb.UploadThumbnailResponse, error) {
	if len(req.FileData) == 0 {
		return nil, status.Error(codes.InvalidArgument, "file data is required")
	}

	if req.FileName == "" {
		return nil, status.Error(codes.InvalidArgument, "file name is required")
	}

	// Upload thumbnail to Cloudinary
	uploadParams := uploader.UploadParams{
		ResourceType:   "image",
		Folder:         "netflix/thumbnails",
		PublicID:       fmt.Sprintf("thumbnail_%s_%s", req.MovieId, req.FileName),
		Transformation: "c_fill,w_1280,h_720",
	}

	result, err := s.cloudinary.Upload.Upload(ctx, req.FileData, uploadParams)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to upload thumbnail: %v", err))
	}

	file := &models.File{
		PublicID:  result.PublicID,
		SecureURL: result.SecureURL,
		FileType:  "image",
		Format:    result.Format,
		Bytes:     result.Bytes,
	}

	if err := s.repo.SaveFile(file); err != nil {
		fmt.Printf("failed to save file metadata: %v\n", err)
	}

	return &pb.UploadThumbnailResponse{
		Url:      result.SecureURL,
		PublicId: result.PublicID,
		Message:  "Thumbnail uploaded successfully",
	}, nil
}

func (s *UploadService) UploadAvatar(ctx context.Context, req *pb.UploadAvatarRequest) (*pb.UploadAvatarResponse, error) {
	if len(req.FileData) == 0 {
		return nil, status.Error(codes.InvalidArgument, "file data is required")
	}

	if req.FileName == "" {
		return nil, status.Error(codes.InvalidArgument, "file name is required")
	}

	// Upload avatar to Cloudinary
	uploadParams := uploader.UploadParams{
		ResourceType:   "image",
		Folder:         "netflix/avatars",
		PublicID:       fmt.Sprintf("avatar_%s_%s", req.UserId, req.FileName),
		Transformation: "c_fill,w_200,h_200,g_face",
	}

	result, err := s.cloudinary.Upload.Upload(ctx, req.FileData, uploadParams)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to upload avatar: %v", err))
	}

	file := &models.File{
		PublicID:  result.PublicID,
		SecureURL: result.SecureURL,
		FileType:  "image",
		Format:    result.Format,
		Bytes:     result.Bytes,
	}

	if err := s.repo.SaveFile(file); err != nil {
		fmt.Printf("failed to save file metadata: %v\n", err)
	}

	return &pb.UploadAvatarResponse{
		Url:      result.SecureURL,
		PublicId: result.PublicID,
		Message:  "Avatar uploaded successfully",
	}, nil
}

func (s *UploadService) DeleteFile(ctx context.Context, req *pb.DeleteFileRequest) (*pb.DeleteFileResponse, error) {
	if req.PublicId == "" {
		return nil, status.Error(codes.InvalidArgument, "public_id is required")
	}

	resourceType := req.ResourceType
	if resourceType == "" {
		resourceType = "image"
	}

	// Delete file from Cloudinary
	deleteParams := uploader.DestroyParams{
		ResourceType: resourceType,
		PublicID:     req.PublicId,
	}

	_, err := s.cloudinary.Upload.Destroy(ctx, deleteParams)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to delete file: %v", err))
	}

	return &pb.DeleteFileResponse{
		Message: "File deleted successfully",
	}, nil
}
