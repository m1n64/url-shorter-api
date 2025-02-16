package handlers

import (
	"context"
	links "link-service/internal/links/grpc"
	"link-service/internal/links/services"
)

type LinkServiceServer struct {
	links.UnimplementedLinkServiceServer
	slugService *services.SlugService
}

func NewLinkServiceServer(
	slugService *services.SlugService,
) *LinkServiceServer {
	return &LinkServiceServer{
		slugService: slugService,
	}
}

func (s *LinkServiceServer) GenerateSlug(ctx context.Context, request *links.Empty) (*links.GenerateSlugResponse, error) {
	return &links.GenerateSlugResponse{
		Slug: s.slugService.GenerateSlug(),
	}, nil
}
