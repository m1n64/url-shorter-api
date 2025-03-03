package entities

import "time"

type Click struct {
	ShortLink   string    `gorm:"column:short_link;type:String"`
	TotalClicks uint32    `gorm:"column:total_clicks;type:UInt32"`
	Timestamp   time.Time `gorm:"column:timestamp;type:DateTime"`
}
