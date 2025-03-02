package entities

import (
	"github.com/go-json-experiment/json/v1"
	"time"
)

type AnalyticsEvent struct {
	ShortLink   string    `gorm:"column:short_link;type:String;primaryKey" json:"short_link"`
	Destination string    `gorm:"column:destination;type:String" json:"destination"`
	IP          string    `gorm:"column:ip;type:String" json:"ip"`
	Country     string    `gorm:"column:country;type:String" json:"country"`
	Referer     string    `gorm:"column:referer;type:String" json:"referer"`
	UserAgent   string    `gorm:"column:user_agent;type:String" json:"user_agent"`
	Timestamp   time.Time `gorm:"column:timestamp;type:DateTime64(3)" json:"timestamp"`
}

func (*AnalyticsEvent) TableName() string {
	return "analytics_events"
}

func (e *AnalyticsEvent) UnmarshalJSON(data []byte) error {
	type Alias AnalyticsEvent
	aux := &struct {
		Timestamp int64 `json:"timestamp"`
		*Alias
	}{
		Alias: (*Alias)(e),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	e.Timestamp = time.Unix(aux.Timestamp, 0)

	return nil
}
