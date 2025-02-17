package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/vmihailenco/msgpack/v5"
	"go.uber.org/zap"
	"link-service/internal/links/models"
	"time"
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

func (s *LinkCacheService) GetLinkFromGlobalCache(slug string) (string, error) {
	link, err := s.memcachedClient.Get(slug)
	if err != nil {
		return "", err
	}

	return string(link.Value), nil
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

func (s *LinkCacheService) RemoveLinkFromGlobalCache(slug string) error {
	err := s.memcachedClient.Delete(slug)
	if err != nil {
		s.logger.Error("Error removing link from global cache: ", zap.Error(err))
		return err
	}

	return nil
}

func (s *LinkCacheService) SetLinkInLocalCache(link *models.Link) error {
	body, err := msgpack.Marshal(link)
	if err != nil {
		s.logger.Error("Error marshalling link: ", zap.Error(err))
		return err
	}

	return s.redisClient.Set(context.Background(), s.createRedisKey(link.ID, link.UserID), body, 24*time.Hour).Err()
}

func (s *LinkCacheService) GetLinkFromLocalCache(id uuid.UUID, userID uuid.UUID) (*models.Link, error) {
	link, err := s.redisClient.Get(context.Background(), s.createRedisKey(id, userID)).Result()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			s.logger.Error("Error getting link from local cache: ", zap.Error(err))
		}

		return nil, err
	}

	var linkModel *models.Link
	err = msgpack.Unmarshal([]byte(link), &linkModel)
	if err != nil {
		s.logger.Error("Error unmarshalling link from local cache: ", zap.Error(err))
		return nil, err
	}

	return linkModel, nil
}

func (s *LinkCacheService) RemoveLinkFromLocalCache(id uuid.UUID, userID uuid.UUID) error {
	return s.redisClient.Del(context.Background(), s.createRedisKey(id, userID)).Err()
}

func (s *LinkCacheService) createRedisKey(id uuid.UUID, userID uuid.UUID) string {
	return fmt.Sprintf("%s:%s", id, userID)
}
