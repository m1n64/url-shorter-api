package repositories

import (
	"analytics-service/internal/analytics/entities"
	"time"
)

type Filter struct {
	StartDate *uint64
	EndDate   *uint64
	Device    *string
	Browser   *string
	OS        *string
	Country   *string
	Page      *uint32
	PerPage   *uint32
}

type AnalyticsEventRepository interface {
	Save(event *entities.AnalyticsEvent) error
	SaveClick(shortLink string, timestamp time.Time) error
	GetAll(shortLink string, filters Filter) (*entities.AnalyticsData, error)
	GetClicksPerDay(shortLink string) ([]*entities.Click, error)
	Delete(shortLink string) error
}
