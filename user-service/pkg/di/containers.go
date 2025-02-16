package di

import (
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"os"
	"user-service/internal/users/repositories"
	"user-service/internal/users/services"
	"user-service/pkg/utils"
)

type Dependencies struct {
	Logger      *zap.Logger
	Redis       *redis.Client
	DB          *gorm.DB
	Validator   *validator.Validate
	RabbitMQ    *utils.RabbitMQConnection
	UserRepo    repositories.UserRepository
	TokenRepo   repositories.TokenRepository
	AuthService *services.AuthService
}

func InitDependencies() *Dependencies {
	// Infrastructure
	logger := utils.InitLogs()
	utils.LoadEnv()
	redisConn := utils.CreateRedisConn()
	dbConn := utils.InitDBConnection(os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"))
	utils.StartMigrations()
	validate := utils.InitValidator()
	rabbitMQ := utils.ConnectRabbitMQ(os.Getenv("RABBITMQ_HOST"))
	utils.InitializeQueues()

	// Repositories
	userRepo := repositories.NewUserGormRepository(dbConn)
	tokenRepo := repositories.NewTokenGormRepository(dbConn)

	// Services
	authService := services.NewAuthService(userRepo, tokenRepo, dbConn)

	return &Dependencies{
		Logger:      logger,
		Redis:       redisConn,
		DB:          dbConn,
		Validator:   validate,
		RabbitMQ:    rabbitMQ,
		UserRepo:    userRepo,
		TokenRepo:   tokenRepo,
		AuthService: authService,
	}
}
