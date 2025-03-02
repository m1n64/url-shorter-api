package repositories

import (
	"analytics-service/internal/analytics/entities"
	"gorm.io/gorm"
	"time"
)

type analyticsEventClickHouseRepository struct {
	clickHouse *gorm.DB
}

func NewAnalyticsEventClickHouseRepository(clickHouse *gorm.DB) AnalyticsEventRepository {
	return &analyticsEventClickHouseRepository{
		clickHouse: clickHouse,
	}
}

func (r *analyticsEventClickHouseRepository) Save(event *entities.AnalyticsEvent) error {
	return r.clickHouse.Create(event).Error
}

func (r *analyticsEventClickHouseRepository) GetAll(shortLink string, startDate *time.Time, endDate *time.Time) ([]*entities.AnalyticsEvent, error) {
	var events []*entities.AnalyticsEvent

	query := r.clickHouse.Where("short_link = ?", shortLink)

	if startDate != nil {
		query = query.Where("timestamp >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("timestamp <= ?", *endDate)
	}

	err := query.Find(&events).Error
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (r *analyticsEventClickHouseRepository) Delete(shortLink string) error {
	return r.clickHouse.Where("short_link = ?", shortLink).Delete(&entities.AnalyticsEvent{}).Error
}
