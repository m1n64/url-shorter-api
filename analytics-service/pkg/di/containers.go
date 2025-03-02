package di

import (
	"analytics-service/internal/analytics/consumers"
	"analytics-service/internal/analytics/repositories"
	"analytics-service/internal/analytics/services"
	"analytics-service/pkg/chromium"
	services2 "analytics-service/pkg/infrastructure/services"
	"analytics-service/pkg/utils"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"os"
)

type Dependencies struct {
	Logger                *zap.Logger
	Redis                 *redis.Client
	DB                    *gorm.DB
	ClickHouseDB          *gorm.DB
	Validator             *validator.Validate
	RabbitMQ              *utils.RabbitMQConnection
	ChromeAllocator       *chromium.ChromeAllocator
	AnalyticsEventRepo    repositories.AnalyticsEventRepository
	CountryService        *services2.CountryService
	AnalyticsEventService *services.AnalyticsEventService
}

func InitDependencies() *Dependencies {
	// Infrastructure
	logger := utils.InitLogs()
	utils.LoadEnv()
	redisConn := utils.CreateRedisConn()
	dbConn := utils.InitDBConnection(os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))
	dbClickHouse := utils.InitClickHouseConnection(os.Getenv("CLICKHOUSE_HOST"), os.Getenv("CLICKHOUSE_PORT"), os.Getenv("CLICKHOUSE_USER"), os.Getenv("CLICKHOUSE_PASSWORD"), os.Getenv("CLICKHOUSE_DB"))
	utils.StartMigrations(dbConn, dbClickHouse)
	validate := utils.InitValidator()
	rabbitMQ := utils.ConnectRabbitMQ(os.Getenv("RABBITMQ_HOST"))
	utils.InitializeQueues()

	chromeAllocator := chromium.NewChromeAllocator()
	chromeAllocator.Init()

	// Repositories
	analyticsEventRepository := repositories.NewAnalyticsEventClickHouseRepository(dbClickHouse)

	// Services
	countryService := services2.NewCountryService(logger)
	analyticsEventService := services.NewAnalyticsEventService(analyticsEventRepository, countryService, logger)

	return &Dependencies{
		Logger:                logger,
		Redis:                 redisConn,
		DB:                    dbConn,
		ClickHouseDB:          dbClickHouse,
		Validator:             validate,
		RabbitMQ:              rabbitMQ,
		ChromeAllocator:       chromeAllocator,
		AnalyticsEventRepo:    analyticsEventRepository,
		CountryService:        countryService,
		AnalyticsEventService: analyticsEventService,
	}
}

func InitializeQueuesConsumer(dependencies *Dependencies) {
	analyticsConsumer := consumers.NewAnalyticsConsumer(dependencies.AnalyticsEventService, dependencies.Logger)

	err := utils.ListenToQueue(utils.AnalyticsQueue, analyticsConsumer)
	if err != nil {
		dependencies.Logger.Error("Error listening to queue", zap.Error(err))
	}
}
