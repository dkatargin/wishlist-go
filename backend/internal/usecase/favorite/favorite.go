// Package favorite — сохранение чужих шаренных списков в «избранное».
package favorite

import (
	"context"
	"wishlist-go/internal/domain"

	"github.com/google/uuid"
)

type Service struct {
	favRepo      domain.FavoriteRepository
	wishlistRepo domain.WishlistRepository
}

func NewService(favRepo domain.FavoriteRepository, wishlistRepo domain.WishlistRepository) *Service {
	return &Service{favRepo: favRepo, wishlistRepo: wishlistRepo}
}

// Add сохраняет список в избранное. Список должен существовать; добавлять можно и
// собственный список (фронт не различает). Идемпотентность — на стороне репозитория.
func (s *Service) Add(ctx context.Context, accountID int64, shareCode uuid.UUID) error {
	if _, err := s.wishlistRepo.GetWishlistByCode(ctx, shareCode); err != nil {
		return err // ErrWishlistNotFound пробрасывается как есть
	}
	return s.favRepo.Add(ctx, &domain.Favorite{AccountID: accountID, WishlistCode: shareCode})
}

func (s *Service) Remove(ctx context.Context, accountID int64, shareCode uuid.UUID) error {
	return s.favRepo.Remove(ctx, accountID, shareCode)
}

// List возвращает сами карточки списков, добавленных пользователем в избранное.
func (s *Service) List(ctx context.Context, accountID int64) ([]*domain.Wishlist, error) {
	return s.favRepo.ListByAccount(ctx, accountID)
}
