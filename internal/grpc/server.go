package grpc

import (
	"context"
	thumbnail1 "getthumbnails/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverAPI struct {
	thumbnail1.UnimplementedThumbnailServer // такая заглушка если еще не все ручки есть
	load                                    Load
}

type Load interface {
	LoadThumbnail(ctx context.Context,
		url string) (urlID string, err error)
}

func Register(gRPC *grpc.Server, load Load) { // регистрация сервера
	thumbnail1.RegisterThumbnailServer(gRPC, &serverAPI{load: load})
}

func (s *serverAPI) GetThumbnail(ctx context.Context,
	req *thumbnail1.ThumbnailRequest,
) (*thumbnail1.ThumbnailResponse, error) {
	if req.GetUrl() == "" {
		return nil, status.Error(codes.InvalidArgument, "Empty url")
	}
	id, err := s.load.LoadThumbnail(ctx, req.GetUrl())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Failed load thumbnail")
	}

	return &thumbnail1.ThumbnailResponse{
		Thumbnail: id,
	}, nil
}
