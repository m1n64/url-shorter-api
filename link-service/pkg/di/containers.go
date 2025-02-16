package di

import (
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"link-service/internal/links/services"
	"link-service/pkg/utils"
	"os"
)

type Dependencies struct {
	Logger      *zap.Logger
	KeyDB       *redis.Client
	DB          *gorm.DB
	Validator   *validator.Validate
	RabbitMQ    *utils.RabbitMQConnection
	SlugService *services.SlugService
}

func InitDependencies() *Dependencies {
	// Infrastructure
	logger := utils.InitLogs()
	utils.LoadEnv()
	redisConn := utils.CreateRedisConn()
	dbConn := utils.InitDBConnection(os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))
	utils.StartMigrations()
	validate := utils.InitValidator()
	rabbitMQ := utils.ConnectRabbitMQ(os.Getenv("RABBITMQ_HOST"))
	utils.InitializeQueues()

	// Repositories

	// Services
	slugService := services.NewSlugService()

	return &Dependencies{
		Logger:      logger,
		KeyDB:       redisConn,
		DB:          dbConn,
		Validator:   validate,
		RabbitMQ:    rabbitMQ,
		SlugService: slugService,
	}
}
