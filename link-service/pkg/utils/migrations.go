package utils

func StartMigrations() {
	db := GetDBConnection()

	db.AutoMigrate()
}
