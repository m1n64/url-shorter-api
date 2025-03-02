package repositories

import (
	"analytics-service/internal/analytics/entities"
	"time"
)

type AnalyticsEventRepository interface {
	Save(event *entities.AnalyticsEvent) error
	GetAll(shortLink string, startDate *time.Time, endDate *time.Time) ([]*entities.AnalyticsEvent, error)
	Delete(shortLink string) error
}
