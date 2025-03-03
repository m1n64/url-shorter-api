package utils

import (
	"gorm.io/gorm"
)

func StartMigrations(db *gorm.DB, dbClickHouse *gorm.DB) {
	db.AutoMigrate()

	/*dbClickHouse.AutoMigrate(&entities.AnalyticsEvent{})*/

	dbClickHouse.Exec(`
		CREATE TABLE IF NOT EXISTS analytics_events (
			short_link String,
			destination String,
			ip String,
			country String,
			referer String,
			user_agent String,
			device String,
			os String,
			os_version String,
			browser String,
			browser_version String,
			timestamp DateTime
		) ENGINE = MergeTree()
		ORDER BY (timestamp, short_link)
	`)

	dbClickHouse.Exec(`
		CREATE TABLE IF NOT EXISTS clicks_summary (
			short_link String,
			total_clicks UInt32,
			timestamp DateTime
		) ENGINE = SummingMergeTree()
		ORDER BY (short_link, timestamp)
	`)
}
