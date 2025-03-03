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

func (r *analyticsEventClickHouseRepository) SaveClick(shortLink string, timestamp time.Time) error {
	return r.clickHouse.Exec(`
		INSERT INTO clicks_summary (short_link, total_clicks, timestamp) 
		VALUES (?, 1, ?)
	`, shortLink, timestamp).Error
}

func (r *analyticsEventClickHouseRepository) GetAll(shortLink string, filters Filter) (*entities.AnalyticsData, error) {
	var events []*entities.AnalyticsEvent
	var totalClicks uint32
	var uniqueClicks uint32

	r.clickHouse.Debug()

	event := entities.AnalyticsEvent{}
	query := r.clickHouse.Table(event.TableName()).Where("short_link = ?", shortLink)

	if filters.StartDate != nil {
		query = query.Where("timestamp >= ?", *filters.StartDate)
	}
	if filters.EndDate != nil {
		query = query.Where("timestamp <= ?", *filters.EndDate)
	}
	if filters.Device != nil {
		query = query.Where("device = ?", *filters.Device)
	}
	if filters.Browser != nil {
		query = query.Where("browser = ?", *filters.Browser)
	}
	if filters.OS != nil {
		query = query.Where("os = ?", *filters.OS)
	}
	if filters.Country != nil {
		query = query.Where("country = ?", *filters.Country)
	}

	err := r.clickHouse.Table(event.TableName()).
		Where("short_link = ?", shortLink).
		Select("COUNT(*) AS total_clicks").
		Scan(&totalClicks).Error
	if err != nil {
		return nil, err
	}

	err = r.clickHouse.Table(event.TableName()).
		Where("short_link = ?", shortLink).
		Select("COUNT(DISTINCT concat(ip, user_agent)) AS unique_clicks").
		Scan(&uniqueClicks).Error
	if err != nil {
		return nil, err
	}

	page := uint32(1)
	perPage := uint32(10)

	if filters.Page != nil {
		page = *filters.Page
	}
	if filters.PerPage != nil {
		perPage = *filters.PerPage
	}

	totalPages := uint32((totalClicks + perPage - 1) / perPage)
	offset := (page - 1) * perPage

	err = query.
		Order("timestamp DESC").
		Limit(int(perPage)).
		Offset(int(offset)).
		Find(&events).Error
	if err != nil {
		return nil, err
	}

	return &entities.AnalyticsData{
		Events:       events,
		TotalClicks:  totalClicks,
		UniqueClicks: uniqueClicks,
		Page:         page,
		PerPage:      perPage,
		TotalPages:   totalPages,
	}, nil
}

func (r *analyticsEventClickHouseRepository) GetClicksPerDay(shortLink string) ([]*entities.Click, error) {
	var clicks []*entities.Click

	err := r.clickHouse.Raw(`
		SELECT 
			short_link, 
			sum(total_clicks) AS total_clicks, 
			toStartOfDay(timestamp) AS timestamp
		FROM clicks_summary
		WHERE short_link = ?
		GROUP BY short_link, timestamp
		ORDER BY timestamp DESC
	`, shortLink).Find(&clicks).Error

	if err != nil {
		return nil, err
	}

	return clicks, nil
}

func (r *analyticsEventClickHouseRepository) Delete(shortLink string) error {
	return r.clickHouse.Where("short_link = ?", shortLink).Delete(&entities.AnalyticsEvent{}).Error
}
