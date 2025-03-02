package utils

import (
	"gorm.io/gorm"
)

func StartMigrations(db *gorm.DB, dbClickHouse *gorm.DB) {
	db.AutoMigrate()

	/*dbClickHouse.AutoMigrate(&entities.AnalyticsEvent{})*/

	dbClickHouse.Exec(`
		CREATE TABLE IF NOT EXISTS analytics_events (
			short_link String PRIMARY KEY,
			destination String,
			ip String,
			country String,
			referer String,
			user_agent String,
			timestamp DateTime
		) ENGINE = MergeTree()
		ORDER BY short_link
	`)
}
