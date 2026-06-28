package wishitem

import (
	"context"
	"errors"
	"net/url"
	"strings"
	"wishlist-go/internal/domain"

	"github.com/google/uuid"
)

// ErrInvalidMarketURL — market_url не из поддерживаемого маркетплейса.
var ErrInvalidMarketURL = errors.New("unsupported market url")

// Publisher публикует сообщения в очередь задач (продьюсер).
type Publisher interface {
	PublishMessage(msgType string, payload map[string]interface{}) error
}

type Service struct {
	repo         domain.WishItemRepository
	wishlistRepo domain.WishlistRepository
	publisher    Publisher
}

func NewService(repo domain.WishItemRepository, wishlistRepo domain.WishlistRepository, publisher Publisher) *Service {
	return &Service{
		repo:         repo,
		wishlistRepo: wishlistRepo,
		publisher:    publisher,
	}
}

// RequestCrawl публикует задачу на парсинг товара по URL; consumer дозаполнит WishItem.
func (s *Service) RequestCrawl(ctx context.Context, shareCode uuid.UUID, ownerID int64, marketURL string) error {
	if !isAllowedMarketURL(marketURL) {
		return ErrInvalidMarketURL
	}
	return s.publisher.PublishMessage("crawl_product", map[string]interface{}{
		"product_url":    marketURL,
		"wish_list_code": shareCode.String(),
		"owner_id":       ownerID,
	})
}

// isAllowedMarketURL пропускает только https-URL поддерживаемого маркетплейса
// (host сверяется строго, без подстрочного обхода) — защита от SSRF на произвольный хост.
func isAllowedMarketURL(raw string) bool {
	u, err := url.Parse(raw)
	if err != nil || u.Scheme != "https" {
		return false
	}
	if u.Host == "market.yandex.ru" {
		return true
	}
	return u.Host == "yandex.ru" && strings.HasPrefix(u.Path, "/products")
}

func (s *Service) GetWishItemsByWishlist(ctx context.Context, shareCode uuid.UUID, limit int, offset int) ([]*domain.WishItem, error) {
	return s.repo.GetWishItemsByWishlistID(shareCode, limit, offset, false)
}

// GetVisibleByWishlist — для гостевого просмотра: только невыполненные (is_done=false).
func (s *Service) GetVisibleByWishlist(ctx context.Context, shareCode uuid.UUID, limit int, offset int) ([]*domain.WishItem, error) {
	return s.repo.GetWishItemsByWishlistID(shareCode, limit, offset, true)
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
