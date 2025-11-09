package wishitem

import (
	"context"
	"wishlist-go/internal/domain"
	"wishlist-go/internal/infrastructure/queue"

	"github.com/google/uuid"
)

type Service struct {
	repo         domain.WishItemRepository
	wishlistRepo domain.WishlistRepository
	mqClient     *queue.RabbitMQClient
}

func NewService(repo domain.WishItemRepository, wishlistRepo domain.WishlistRepository, mqClient *queue.RabbitMQClient) *Service {
	return &Service{
		repo:         repo,
		wishlistRepo: wishlistRepo,
		mqClient:     mqClient,
	}
}

func (s *Service) GetWishItemsByWishlist(ctx context.Context, shareCode uuid.UUID, limit int, offset int) ([]*domain.WishItem, error) {
	return s.repo.GetWishItemsByWishlistID(shareCode, limit, offset)
}

func (s *Service) GetWishItemByID(ctx context.Context, id int64, wlCode uuid.UUID) (*domain.WishItem, error) {
	return s.repo.GetWishItemByID(id, wlCode)
}

func (s *Service) CreateWishItem(ctx context.Context, wlCode uuid.UUID, marketURL string, name *string, marketPictureURL *string, marketPrice *float64) (*domain.WishItem, error) {

	if name == nil {
		//	TODO: отправляем в очередь на парсинг
		return nil, nil
	}

	wi := &domain.WishItem{
		WishListCode:     wlCode,
		MarketURL:        marketURL,
		Name:             name,
		MarketPictureURL: marketPictureURL,
		MarketPrice:      marketPrice,
	}
	if err := s.repo.CreateWishItem(wi); err != nil {
		return nil, err
	}
	return wi, nil
}

func (s *Service) UpdateWishItem(ctx context.Context, wishItem *domain.WishItem) error {
	if err := s.repo.UpdateWishItem(wishItem); err != nil {
		return err
	}
	return nil
}

func (s *Service) DeleteWishItem(ctx context.Context, id int64) error {
	return s.repo.DeleteWishItem(id)
}
