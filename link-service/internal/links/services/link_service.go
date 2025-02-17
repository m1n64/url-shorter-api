package services

import (
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"link-service/internal/links/models"
	"link-service/internal/links/repositories"
)

type LinkService struct {
	linkRepository   repositories.LinkRepository
	linkCacheService *LinkCacheService
	db               *gorm.DB
	logger           *zap.Logger
}

type CreateLinkRequest struct {
	UserId string  `validate:"required,uuid"`
	Url    string  `validate:"required,url"`
	Slug   *string `validate:"omitempty"`
}

func NewLinkService(linkRepository repositories.LinkRepository, linkCacheService *LinkCacheService, db *gorm.DB, logger *zap.Logger) *LinkService {
	return &LinkService{
		linkRepository:   linkRepository,
		linkCacheService: linkCacheService,
		db:               db,
		logger:           logger,
	}
}

func (s *LinkService) GetAll() ([]*models.Link, error) {
	return s.linkRepository.GetAll()
}

func (s *LinkService) SaveAllInGlobalCache() error {
	links, err := s.GetAll()
	if err != nil {
		s.logger.Error("Error getting links: ", zap.Error(err))
		return err
	}

	for _, link := range links {
		err = s.linkCacheService.SaveLinkInGlobalCache(link)
		if err != nil {
			s.logger.Error("Error saving link in global cache: ", zap.Error(err))
		}
	}

	return nil
}

func (s *LinkService) GetAllByUserID(userID uuid.UUID) ([]*models.Link, error) {
	return s.linkRepository.GetAllByUserID(userID)
}

func (s *LinkService) GetByID(userID uuid.UUID, id uuid.UUID) (*models.Link, error) {
	link, err := s.linkRepository.GetByID(userID, id)
	if err != nil {
		s.logger.Error("Error getting link by ID: ", zap.Error(err))
		return nil, err
	}

	return link, nil
}

func (s *LinkService) GetBySlug(slug string) (*models.Link, error) {
	if cachedLink, err := s.linkCacheService.GetLinkFromLocalCache(slug); err == nil {
		return cachedLink, nil
	}

	link, err := s.linkRepository.GetBySlug(slug)

	err = s.linkCacheService.SetLinkInLocalCache(link)
	if err != nil {
		s.logger.Error("Error setting link in local cache: ", zap.Error(err))
	}

	return link, nil
}

func (s *LinkService) Create(link *CreateLinkRequest) (*models.Link, error) {
	linkModel := &models.Link{
		UserID: uuid.MustParse(link.UserId),
		Slug:   *link.Slug,
		URL:    link.Url,
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		err := s.linkRepository.Save(linkModel)
		if err != nil {
			return err
		}

		err = s.linkCacheService.SaveLinkInGlobalCache(linkModel)
		if err != nil {
			tx.Rollback()
			return err
		}

		err = s.linkCacheService.SetLinkInLocalCache(linkModel)
		if err != nil {
		}

		return nil
	})
	if err != nil {
		s.logger.Error("Error creating link: ", zap.Error(err))
		return nil, err
	}

	return linkModel, nil
}

func (s *LinkService) Delete(userID uuid.UUID, id uuid.UUID) error {
	link, err := s.GetByID(userID, id)
	if err != nil {
		s.logger.Error("Error getting link by ID: ", zap.Error(err))
		return err
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		err := s.linkRepository.DeleteByID(userID, id)
		if err != nil {
			s.logger.Error("Error deleting link: ", zap.Error(err))
			return err
		}

		err = s.linkCacheService.RemoveLinkFromGlobalCache(link.Slug)
		if err != nil {
			tx.Rollback()
			return err
		}

		err = s.linkCacheService.RemoveLinkFromLocalCache(link.Slug)
		if err != nil {
			tx.Rollback()
			return err
		}

		return nil
	})

	return err
}
