package services

import (
	"analytics-service/internal/analytics/entities"
	"analytics-service/internal/analytics/repositories"
	"analytics-service/pkg/infrastructure/services"
	"go.uber.org/zap"
)

type AnalyticsEventService struct {
	analyticsEventRepo     repositories.AnalyticsEventRepository
	countryService         *services.CountryService
	userAgentParserService *services.UserAgentParserService
	logger                 *zap.Logger
}

func NewAnalyticsEventService(analyticsEventRepo repositories.AnalyticsEventRepository, countryService *services.CountryService, userAgentParserService *services.UserAgentParserService, logger *zap.Logger) *AnalyticsEventService {
	return &AnalyticsEventService{
		analyticsEventRepo:     analyticsEventRepo,
		countryService:         countryService,
		userAgentParserService: userAgentParserService,
		logger:                 logger,
	}
}

func (s *AnalyticsEventService) Save(event *entities.AnalyticsEvent) (*entities.AnalyticsEvent, error) {
	if event.Country == "" {
		event.Country = s.countryService.GetCountryByIP(event.IP)
	}

	parserUA := s.userAgentParserService.ParseUserAgent(event.UserAgent)

	event.OS = parserUA.OS
	event.OSVersion = parserUA.OSVersion
	event.Device = parserUA.Device
	event.Browser = parserUA.Browser
	event.BrowserVer = parserUA.BrowserVer

	if err := s.analyticsEventRepo.Save(event); err != nil {
		s.logger.Error("Error saving analytics event", zap.Error(err))
		return nil, err
	}

	if err := s.analyticsEventRepo.SaveClick(event.ShortLink, event.Timestamp); err != nil {
		s.logger.Error("Error saving click", zap.Error(err))
		return nil, err
	}

	return event, nil
}

func (s *AnalyticsEventService) GetAll(
	shortLink string,
	startDate *uint64,
	endDate *uint64,
	device *string,
	browser *string,
	os *string,
	country *string,
	page *uint32,
	perPage *uint32,
) (*entities.AnalyticsData, error) {
	filters := repositories.Filter{
		StartDate: startDate,
		EndDate:   endDate,
		Device:    device,
		Browser:   browser,
		OS:        os,
		Country:   country,
		Page:      page,
		PerPage:   perPage,
	}

	events, err := s.analyticsEventRepo.GetAll(shortLink, filters)
	if err != nil {
		s.logger.Error("Error getting analytics events", zap.Error(err))
		return nil, err
	}

	return events, nil
}

func (s *AnalyticsEventService) GetClicksPerDay(shortLink string) ([]*entities.Click, error) {
	return s.analyticsEventRepo.GetClicksPerDay(shortLink)
}
