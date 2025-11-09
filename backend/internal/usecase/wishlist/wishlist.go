package wishlist

import (
	"context"
	"log"
	"wishlist-go/internal/domain"

	"github.com/google/uuid"
)

type Service struct {
	repo domain.WishlistRepository
}

func NewService(repo domain.WishlistRepository) *Service {
	return &Service{repo: repo}
}

// GetWishlistsByOwner получаем список всех вишлистов пользователя
func (s *Service) GetWishlistsByOwner(ctx context.Context, ownerId int64, offset int, limit int) ([]*domain.Wishlist, error) {
	return s.repo.GetWishlistsByOwnerID(ctx, ownerId, offset, limit)
}

// CreateWishlist создаем новый вишлист
func (s *Service) CreateWishlist(ctx context.Context, ownerId int64, name string, description *string) (*domain.Wishlist, error) {
	wl := &domain.Wishlist{
		ShareCode:   uuid.New(),
		Name:        name,
		OwnerID:     ownerId,
		Description: description,
	}
	if err := s.repo.CreateWishlist(ctx, wl); err != nil {
		return nil, err
	}
	return wl, nil
}

// UpdateWishlist обновляем вишлист
func (s *Service) UpdateWishlist(ctx context.Context, shareCode uuid.UUID, name *string, description *string) error {
	if name == nil && description == nil {
		return nil
	}
	wl := &domain.Wishlist{
		ShareCode: shareCode,
	}
	if name != nil {
		wl.Name = *name
	}
	if description != nil {
		wl.Description = description
	}

	return s.repo.UpdateWishlist(ctx, wl)
}

// DeleteWishlist удаляем вишлист по shareCode
func (s *Service) DeleteWishlist(ctx context.Context, shareCode uuid.UUID) error {
	return s.repo.DeleteWishlist(ctx, shareCode)
}

// CheckAccess проверяем, есть ли у пользователя доступ к вишлисту
func (s *Service) CheckAccess(ctx context.Context, shareCode uuid.UUID, userId int64) bool {
	wl, err := s.repo.GetWishlistByCode(ctx, shareCode)
	if err != nil {
		log.Printf("error checking wishlist access: %v", err)
		return false
	}
	if wl == nil {
		log.Printf("requested wishlist not found: %s", shareCode.String())
		return false
	}
	return wl.OwnerID == userId
}
