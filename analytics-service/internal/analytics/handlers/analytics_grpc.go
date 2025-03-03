package handlers

import (
	analytics "analytics-service/internal/analytics/grpc"
	"analytics-service/internal/analytics/services"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AnalyticsServiceServer struct {
	analytics.UnimplementedAnalyticsServiceServer
	analyticsEventService *services.AnalyticsEventService
}

func NewAnalyticsServiceServer(analyticsEventService *services.AnalyticsEventService) *AnalyticsServiceServer {
	return &AnalyticsServiceServer{
		analyticsEventService: analyticsEventService,
	}
}

func (s *AnalyticsServiceServer) GetGeneralStats(ctx context.Context, request *analytics.AnalyticsRequest) (*analytics.AnalyticsResponse, error) {
	if request.ShortUrl == "" {
		return nil, status.Error(codes.InvalidArgument, "short_url is required")
	}

	stats, err := s.analyticsEventService.GetAll(request.ShortUrl, request.StartDate, request.EndDate, request.Device, request.Browser, request.Os, request.Country, request.Page, request.PerPage)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var analyticsList []*analytics.Analytics
	for _, stat := range stats.Events {
		analyticsList = append(analyticsList, &analytics.Analytics{
			ShortUrl:       stat.ShortLink,
			Ip:             stat.IP,
			Country:        stat.Country,
			Browser:        stat.Browser,
			BrowserVersion: stat.BrowserVer,
			Os:             stat.OS,
			OsVersion:      stat.OSVersion,
			UserAgent:      stat.UserAgent,
			Device:         stat.Device,
			Timestamp:      uint64(stat.Timestamp.Unix()),
		})
	}

	return &analytics.AnalyticsResponse{
		ShortUrl:     request.ShortUrl,
		Analytics:    analyticsList,
		TotalClicks:  stats.TotalClicks,
		UniqueClicks: stats.UniqueClicks,
		Page:         stats.Page,
		TotalPages:   stats.TotalPages,
		PerPage:      stats.PerPage,
	}, nil
}

func (s *AnalyticsServiceServer) GetClicksPerDay(ctx context.Context, request *analytics.AnalyticsCleanRequest) (*analytics.ClicksByDayResponse, error) {
	if request.ShortUrl == "" {
		return nil, status.Error(codes.InvalidArgument, "short_url is required")
	}

	clicksByDay, err := s.analyticsEventService.GetClicksPerDay(request.ShortUrl)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var totalClicksList []*analytics.Click
	for _, stat := range clicksByDay {
		totalClicksList = append(totalClicksList, &analytics.Click{
			TotalClicks: stat.TotalClicks,
			Timestamp:   uint64(stat.Timestamp.Unix()),
		})
	}

	return &analytics.ClicksByDayResponse{
		Clicks: totalClicksList,
	}, nil
}
