package services

import (
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"link-service/internal/links/models"
)

type LinkCacheService struct {
	memcachedClient *memcache.Client
	redisClient     *redis.Client
	logger          *zap.Logger
}

func NewLinkCacheService(memcachedClient *memcache.Client, redisClient *redis.Client, logger *zap.Logger) *LinkCacheService {
	return &LinkCacheService{
		memcachedClient: memcachedClient,
		redisClient:     redisClient,
		logger:          logger,
	}
}

func (s *LinkCacheService) SaveLinkInGlobalCache(link *models.Link) error {
	err := s.memcachedClient.Set(&memcache.Item{
		Key:        link.Slug,
		Value:      []byte(link.URL),
		Expiration: 0,
	})
	if err != nil {
		s.logger.Error("Error saving link in global cache: ", zap.Error(err))
		return err
	}

	return nil
}

func (s *LinkCacheService) RemoveLinkFromGlobalCache(link *models.Link) error {
	err := s.memcachedClient.Delete(link.Slug)
	if err != nil {
		s.logger.Error("Error removing link from global cache: ", zap.Error(err))
		return err
	}

	return nil
}
