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
	mqClient     *queue.RabbitMQClient // P1.5: продьюсер crawl_product для POST .../wishes/crawl
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

func (s *Service) CreateWishItem(ctx context.Context, item *domain.WishItem) (*domain.WishItem, error) {
	if item.MarketCurrency == "" {
		item.MarketCurrency = "RUB"
	}
	if err := s.repo.CreateWishItem(item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *Service) UpdateWishItem(ctx context.Context, id int64, shareCode uuid.UUID, upd domain.WishItemUpdate) error {
	return s.repo.UpdateWishItem(id, shareCode, upd)
}

func (s *Service) DeleteWishItem(ctx context.Context, id int64) error {
	return s.repo.DeleteWishItem(id)
}
