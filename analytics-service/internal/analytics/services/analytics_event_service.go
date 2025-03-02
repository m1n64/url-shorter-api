package services

import (
	"analytics-service/internal/analytics/entities"
	"analytics-service/internal/analytics/repositories"
	"analytics-service/pkg/infrastructure/services"
	"go.uber.org/zap"
)

type AnalyticsEventService struct {
	analyticsEventRepo repositories.AnalyticsEventRepository
	countryService     *services.CountryService
	logger             *zap.Logger
}

func NewAnalyticsEventService(analyticsEventRepo repositories.AnalyticsEventRepository, countryService *services.CountryService, logger *zap.Logger) *AnalyticsEventService {
	return &AnalyticsEventService{
		analyticsEventRepo: analyticsEventRepo,
		countryService:     countryService,
		logger:             logger,
	}
}

func (s *AnalyticsEventService) Save(event *entities.AnalyticsEvent) (*entities.AnalyticsEvent, error) {
	if event.Country == "" {
		event.Country = s.countryService.GetCountryByIP(event.IP)
	}

	if err := s.analyticsEventRepo.Save(event); err != nil {
		s.logger.Error("Error saving analytics event", zap.Error(err))
		return nil, err
	}

	return event, nil
}
