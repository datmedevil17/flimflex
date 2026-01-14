package clients

import (
	"fmt"

	"github.com/datmedevil17/micro-flex/api-gateway/config"
	authpb "github.com/datmedevil17/micro-flex/proto/auth"
	moviepb "github.com/datmedevil17/micro-flex/proto/movie"
	recommendationpb "github.com/datmedevil17/micro-flex/proto/recommendation"
	streamingpb "github.com/datmedevil17/micro-flex/proto/streaming"
	uploadpb "github.com/datmedevil17/micro-flex/proto/upload"
	userpb "github.com/datmedevil17/micro-flex/proto/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClients struct {
	AuthClient           authpb.AuthServiceClient
	UserClient           userpb.UserServiceClient
	MovieClient          moviepb.MovieServiceClient
	UploadClient         uploadpb.UploadServiceClient
	StreamingClient      streamingpb.StreamingServiceClient
	RecommendationClient recommendationpb.RecommendationServiceClient
}

func NewGRPCClients(cfg config.ServicesConfig) (*GRPCClients, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	authConn, err := grpc.Dial(cfg.AuthURL, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to auth service: %w", err)
	}

	userConn, err := grpc.Dial(cfg.UserURL, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user service: %w", err)
	}

	movieConn, err := grpc.Dial(cfg.MovieURL, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to movie service: %w", err)
	}

	uploadConn, err := grpc.Dial(cfg.UploadURL, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to upload service: %w", err)
	}

	streamingConn, err := grpc.Dial(cfg.StreamingURL, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to streaming service: %w", err)
	}

	recommendationConn, err := grpc.Dial(cfg.RecommendationURL, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to recommendation service: %w", err)
	}

	return &GRPCClients{
		AuthClient:           authpb.NewAuthServiceClient(authConn),
		UserClient:           userpb.NewUserServiceClient(userConn),
		MovieClient:          moviepb.NewMovieServiceClient(movieConn),
		UploadClient:         uploadpb.NewUploadServiceClient(uploadConn),
		StreamingClient:      streamingpb.NewStreamingServiceClient(streamingConn),
		RecommendationClient: recommendationpb.NewRecommendationServiceClient(recommendationConn),
	}, nil
}
