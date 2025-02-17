package handlers

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	links "link-service/internal/links/grpc"
	"link-service/internal/links/models"
	"link-service/internal/links/services"
	"time"
)

type LinkServiceServer struct {
	links.UnimplementedLinkServiceServer
	linkService *services.LinkService
	slugService *services.SlugService
	validator   *validator.Validate
}

type CreateLinkRequest struct {
	UserId string  `validate:"required,uuid"`
	Url    string  `validate:"required,url"`
	Slug   *string `validate:"omitempty"`
}

func NewLinkServiceServer(
	linkService *services.LinkService,
	slugService *services.SlugService,
	validator *validator.Validate,
) *LinkServiceServer {
	return &LinkServiceServer{
		linkService: linkService,
		slugService: slugService,
		validator:   validator,
	}
}

func (s *LinkServiceServer) GetLinks(ctx context.Context, request *links.GetLinksRequest) (*links.GetLinksResponse, error) {
	if uuid.Validate(request.UserId) != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}

	linksList, err := s.linkService.GetAllByUserID(uuid.MustParse(request.UserId))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var linksResp []*links.LinkResponse
	for _, link := range linksList {
		linksResp = append(linksResp, s.getResponse(link))
	}

	return &links.GetLinksResponse{
		Links: linksResp,
	}, nil
}

func (s *LinkServiceServer) GetLink(ctx context.Context, request *links.GetLinkRequest) (*links.LinkResponse, error) {
	if uuid.Validate(request.Id) != nil || uuid.Validate(request.UserId) != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id or user id")
	}

	link, err := s.linkService.GetByID(uuid.MustParse(request.UserId), uuid.MustParse(request.Id))

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return s.getResponse(link), nil
}

func (s *LinkServiceServer) CreateLink(ctx context.Context, request *links.CreateLinkRequest) (*links.LinkResponse, error) {
	req := &CreateLinkRequest{
		UserId: request.UserId,
		Url:    request.Url,
		Slug:   request.Slug,
	}

	if err := s.validator.Struct(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	slug := s.slugService.GenerateSlug()

	if req.Slug != nil {
		slug = *req.Slug
	}

	link, err := s.linkService.Create(uuid.MustParse(req.UserId), req.Url, slug)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return s.getResponse(link), nil
}

func (s *LinkServiceServer) DeleteLink(ctx context.Context, request *links.DeleteLinkRequest) (*links.DeleteLinkResponse, error) {
	if uuid.Validate(request.Id) != nil || uuid.Validate(request.UserId) != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id or user id")
	}

	err := s.linkService.Delete(uuid.MustParse(request.UserId), uuid.MustParse(request.Id))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &links.DeleteLinkResponse{
		Deleted: true,
	}, nil
}

func (s *LinkServiceServer) GenerateSlug(ctx context.Context, request *links.Empty) (*links.GenerateSlugResponse, error) {
	return &links.GenerateSlugResponse{
		Slug: s.slugService.GenerateSlug(),
	}, nil
}

func (s *LinkServiceServer) getResponse(response *models.Link) *links.LinkResponse {
	return &links.LinkResponse{
		Id:        response.ID.String(),
		UserId:    response.UserID.String(),
		Url:       response.URL,
		Slug:      response.Slug,
		CreatedAt: response.CreatedAt.Format(time.RFC3339Nano),
	}
}
