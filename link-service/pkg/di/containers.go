package di

import (
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"link-service/internal/links/repositories"
	"link-service/internal/links/services"
	"link-service/pkg/utils"
	"os"
)

type Dependencies struct {
	Logger           *zap.Logger
	KeyDB            *redis.Client
	Memcached        *memcache.Client
	DB               *gorm.DB
	Validator        *validator.Validate
	RabbitMQ         *utils.RabbitMQConnection
	LinkRepo         repositories.LinkRepository
	SlugService      *services.SlugService
	LinkService      *services.LinkService
	LinkCacheService *services.LinkCacheService
}

func InitDependencies() *Dependencies {
	// Infrastructure
	logger := utils.InitLogs()
	utils.LoadEnv()
	redisConn := utils.CreateRedisConn()
	memcachedConn := utils.CreateMemcachedConn(os.Getenv("MEMCACHED_HOST"), os.Getenv("MEMCACHED_PORT"))
	dbConn := utils.InitDBConnection(os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))
	utils.StartMigrations()
	validate := utils.InitValidator()
	rabbitMQ := utils.ConnectRabbitMQ(os.Getenv("RABBITMQ_HOST"))
	utils.InitializeQueues()

	// Repositories
	linkRepository := repositories.NewLinkDBRepository(dbConn)

	// Services
	slugService := services.NewSlugService()
	linkCacheService := services.NewLinkCacheService(memcachedConn, redisConn, logger)
	linkService := services.NewLinkService(linkRepository, linkCacheService, dbConn, logger)

	return &Dependencies{
		Logger:           logger,
		KeyDB:            redisConn,
		Memcached:        memcachedConn,
		DB:               dbConn,
		Validator:        validate,
		RabbitMQ:         rabbitMQ,
		SlugService:      slugService,
		LinkRepo:         linkRepository,
		LinkService:      linkService,
		LinkCacheService: linkCacheService,
	}
}
